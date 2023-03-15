package repository

import (
	"context"
	"fmt"
	"reflect"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"

	. "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoIdGenCollectionName = "id-gen"

type Repository struct {
	Connection *mongo.Client
	Logger     logger.Ilogger
	Cache      redis.IRedisCache
	DB         *mongo.Database
}

type PaginationKey struct {
	Colllection string
	Page        int64
	Limit       int64
	Filter      interface{}
	OrderBy     string
	Order       entity.SortType
}

type Sort struct {
	SortBy string
	Sort   entity.SortType
}

func NewRepository(g *global.Global) (*Repository, error) {

	clientOption := &options.ClientOptions{}
	opt := &options.DatabaseOptions{
		ReadConcern:    clientOption.ReadConcern,
		WriteConcern:   clientOption.WriteConcern,
		ReadPreference: clientOption.ReadPreference,
		Registry:       clientOption.Registry,
	}

	r := new(Repository)
	connection := g.DBConnection.GetType()
	r.Connection = connection.(*mongo.Client)
	r.Logger = g.Logger
	r.Cache = g.Cache
	r.DB = r.Connection.Database(g.Conf.Databases.Mongo.Name, opt)
	return r, nil
}

func (r Repository) InsertMany(dbName string, objs []entity.IEntity) error {
	bDatas := make([]interface{}, 0)

	for i := 0; i < len(objs); i++ {
		objs[i].SetID()
		objs[i].SetCreatedAt()
		bData, err := objs[i].ToBson()
		if err != nil {
			return err
		}
		bDatas = append(bDatas, *bData)
	}

	opts := options.InsertMany().SetOrdered(false)

	_, err := r.DB.Collection(dbName).InsertMany(context.TODO(), bDatas, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) InsertOne(dbName string, obj entity.IEntity) error {
	obj.SetID()
	obj.SetCreatedAt()

	bData, err := obj.ToBson()
	if err != nil {
		return err
	}

	inserted, err := r.DB.Collection(dbName).InsertOne(context.TODO(), &bData)
	if err != nil {
		return err
	}

	objIDObject := inserted.InsertedID.(primitive.ObjectID)
	objIDStr := objIDObject.Hex()

	r.CreateCache(dbName, objIDStr, obj)
	err = obj.Decode(bData)
	if err != nil {
		return err
	}

	return nil
}

type queriedChan struct {
	Err  error
	Data *primitive.M
}

func (r Repository) FindOne(dbName string, id string) (*bson.M, error) {
	go r.FindOneWithoutCache(dbName, id) // reload cache here

	data, err := r.GetCache(dbName, id)
	if err == nil {
		return data, nil
	}

	data, err = r.FindOneWithoutCache(dbName, id)
	if err == nil {
		return data, nil
	}

	return data, err
}

func (r Repository) FindOneWithoutCache(dbName string, id string) (*bson.M, error) {

	var err error
	data := &primitive.M{}
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{utils.KEY_UUID, id}},
				bson.D{{utils.KEY_DELETED_AT, nil}},
			},
		},
	}

	data, err = r.FilterOne(dbName, filter)
	if err != nil {
		return nil, err
	}
	r.CreateCache(dbName, id, data)
	return data, nil

}

func (r Repository) FilterOne(dbName string, filter bson.D, opts ...*options.FindOneOptions) (*bson.M, error) {
	data := &bson.M{}

	err := r.DB.Collection(dbName).FindOne(context.Background(), filter, opts...).Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r Repository) UpdateOne(dbName string, filter bson.D, obj entity.IEntity) (*mongo.UpdateResult, error) {
	obj.SetUpdatedAt()
	bData, err := obj.ToBson()
	if err != nil {
		return nil, err
	}

	update := bson.D{{"$set", bData}}
	result, err := r.DB.Collection(dbName).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	//Update cache
	id := obj.GetID()
	r.CreateCache(dbName, id, obj)

	return result, nil
}

func (r Repository) DeleteOne(dbName string, filter bson.D) (*mongo.DeleteResult, error) {
	result, err := r.DB.Collection(dbName).DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	filterArr := filter.Map()
	id, ok := filterArr[utils.KEY_UUID]
	if ok {
		r.DeleteCache(dbName, id.(string))
	}

	return result, nil
}

func (b Repository) DeleteMany(ctx context.Context, collectionName string, filters interface{}, opts ...*options.DeleteOptions) (int64, error) {
	rs, err := b.DB.Collection(collectionName).DeleteMany(
		ctx,
		filters,
		opts...,
	)
	if err != nil {
		return 0, err
	}
	return rs.DeletedCount, nil
}

func (r Repository) SoftDelete(obj entity.IEntity) (*mongo.UpdateResult, error) {
	id := obj.GetID()
	dbName := obj.TableName()

	obj.SetDeletedAt()
	filter := bson.D{{utils.KEY_UUID, id}}

	bData, err := obj.ToBson()
	if err != nil {
		return nil, err
	}

	update := bson.D{{"$set", bData}}
	result, err := r.DB.Collection(dbName).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	r.DeleteCache(dbName, id)
	return result, nil
}

func (r Repository) Paginate(dbName string, page int64, limit int64, filter interface{}, selectFields interface{}, sorts []Sort, returnData interface{}) (*PaginatedData, error) {
	paginatedData := New(r.DB.Collection(dbName)).
		Context(context.TODO()).
		Limit(int64(limit)).
		Page(int64(page))

	if len(sorts) > 0 {
		for _, sort := range sorts {
			if sort.Sort == entity.SORT_ASC || sort.Sort == entity.SORT_DESC {
				//sortValue := bson.D{{"created_at", -1}}
				paginatedData.Sort(sort.SortBy, sort.Sort)
			}
		}
	} else {
		paginatedData.Sort("created_at", entity.SORT_DESC)
		paginatedData.Sort(utils.KEY_UUID, entity.SORT_ASC)
	}

	data, err := paginatedData.
		Select(selectFields).
		Filter(filter).
		Decode(returnData).
		Find()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r Repository) Aggregate(dbName string, page int64, limit int64, filter interface{}, selectFields interface{}, sorts []Sort, returnData interface{}) (*PaginatedData, error) {
	paginatedData := New(r.DB.Collection(dbName)).
		Context(context.TODO()).
		Limit(int64(limit)).
		Page(int64(page))

	if len(sorts) > 0 {
		for _, sort := range sorts {
			if sort.Sort == entity.SORT_ASC || sort.Sort == entity.SORT_DESC {
				//sortValue := bson.D{{"created_at", -1}}
				paginatedData.Sort(sort.SortBy, sort.Sort)
			}
		}
	} else {
		paginatedData.Sort("created_at", entity.SORT_DESC)
		paginatedData.Sort(utils.KEY_UUID, entity.SORT_ASC)
	}

	data, err := paginatedData.
		Select(selectFields).
		Decode(returnData).
		Aggregate(filter)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r Repository) CreateCache(dbName string, objID string, obj interface{}) error {
	bytes, err := bson.Marshal(obj)
	if err != nil {
		return err
	}

	stringData := string(bytes)
	key := fmt.Sprintf(utils.DB_CACHE_KEY, dbName, objID)
	err = r.Cache.SetStringDataWithExpTime(key, stringData, utils.DB_CACHE_EXPIRED_TIME)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetCache(dbName string, objID string) (*bson.M, error) {
	key := fmt.Sprintf(utils.DB_CACHE_KEY, dbName, objID)
	data, err := r.Cache.GetData(key)
	if err != nil {
		return nil, err
	}

	resp := &bson.M{}
	err = helpers.ParseCache(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) DeleteCache(dbName string, objID string) error {
	key := fmt.Sprintf(utils.DB_CACHE_KEY, dbName, objID)
	err := r.Cache.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) CreateCollectionIndexes() error {
	// _, err := r.CreateTokenURIIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateProjectIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateMarketplaceListingsIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateMarketplaceOffersIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateProposalIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateProposalVotesIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateBTCWalletIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateMintBTCCIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateWalletTrackTxIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateReferalIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateVolumnIndexModel()
	// if err != nil {
	// 	return err
	// }

	// _, err = r.CreateCategoryIndexModel()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (b Repository) Find(ctx context.Context, collectionName string, filters map[string]interface{}, result interface{}, opts ...*options.FindOptions) error {
	cur, err := b.DB.Collection(collectionName).Find(ctx, filters, opts...)
	if err != nil {
		return err
	}
	if err := cur.All(ctx, result); err != nil {
		return err
	}
	return nil
}
func (r Repository) FindOneBy(ctx context.Context, collectionName string, filters map[string]interface{}, value interface{}, opts ...*options.FindOneOptions) error {
	res := r.DB.Collection(collectionName).FindOne(ctx, filters, opts...)
	if res.Err() != nil {
		return res.Err()
	}
	if err := res.Decode(value); err != nil {
		return err
	}
	return nil
}
func (b Repository) Create(ctx context.Context, collectionName string, model interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
	result, err := b.DB.Collection(collectionName).InsertOne(ctx, model, opts...)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id, nil
	}
	return primitive.NilObjectID, nil
}
func (b Repository) CreateMany(ctx context.Context, collectionName string, models []interface{}, opts ...*options.InsertManyOptions) ([]primitive.ObjectID, error) {
	results, err := b.DB.Collection(collectionName).InsertMany(ctx, models, opts...)
	if err != nil {
		return nil, err
	}
	var ids []primitive.ObjectID
	for _, idResult := range results.InsertedIDs {
		if id, ok := idResult.(primitive.ObjectID); ok {
			ids = append(ids, id)
		}
	}
	return ids, err
}
func (b Repository) UpdateByID(ctx context.Context, collectionName string, id interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	result, err := b.DB.Collection(collectionName).UpdateByID(ctx, id, update, opts...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
func (b Repository) UpdateMany(ctx context.Context, collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	result, err := b.DB.Collection(collectionName).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

type Counter struct {
	ID  string `json:"id" bson:"_id"`
	Seq uint   `json:"seq" bson:"seq"`
}

func (b Repository) NextId(ctx context.Context, sequenceName string) (uint, error) {
	findOptions := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var counter Counter
	err := b.DB.Collection(MongoIdGenCollectionName).
		FindOneAndUpdate(ctx,
			bson.M{"_id": sequenceName},
			bson.M{"$inc": bson.M{"seq": 1}},
			findOptions,
		).Decode(&counter)

	if err != nil {
		return 0, err
	}
	return counter.Seq, nil
}
func (b Repository) Aggregation(ctx context.Context, collectionName string, page int64, limit int64, result interface{}, agg ...interface{}) (int64, error) {
	query := New(b.DB.Collection(collectionName)).Context(ctx).Page(page).Limit(limit)
	aggPaginatedData, err := query.Aggregate(agg...)
	if err != nil {
		return 0, err
	}
	to := indirect(reflect.ValueOf(result))
	toType, _ := indirectType(to.Type())
	if to.IsNil() {
		slice := reflect.MakeSlice(reflect.SliceOf(to.Type().Elem()), 0, int(limit))
		to.Set(slice)
	}
	for i := 0; i < len(aggPaginatedData.Data); i++ {
		ele := reflect.New(toType).Elem().Addr()
		if marshallErr := bson.Unmarshal(aggPaginatedData.Data[i], ele.Interface()); marshallErr == nil {
			to.Set(reflect.Append(to, ele))
		} else {
			return 0, marshallErr
		}
	}
	return aggPaginatedData.Pagination.Total, nil
}
func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}
func indirectType(reflectType reflect.Type) (_ reflect.Type, isPtr bool) {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
		isPtr = true
	}
	return reflectType, isPtr
}
