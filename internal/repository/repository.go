package repository

import (
	"context"
	"fmt"

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

// only delete in mongo
func (r Repository) DeleteMany(dbName string, filter bson.D) (*mongo.DeleteResult, error) {
	result, err := r.DB.Collection(dbName).DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
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
	_, err := r.CreateTokenURIIndexModel()
	if err != nil {
		return err
	}

	_, err = r.CreateProjectIndexModel()
	if err != nil {
		return err
	}

	_, err = r.CreateMarketplaceListingsIndexModel()
	if err != nil {
		return err
	}

	_, err = r.CreateMarketplaceOffersIndexModel()
	if err != nil {
		return err
	}

	_, err = r.CreateProposalIndexModel()
	if err != nil {
		return err
	}

	_, err = r.CreateProposalVotesIndexModel()
	if err != nil {
		return err
	}

	_, err = r.CreateBTCWalletIndexModel()
	if err != nil {
		return err
	}

	_, err = r.CreateMintBTCCIndexModel()
	if err != nil {
		return err
	}

	_, err = r.CreateWalletTrackTxIndexModel()
	if err != nil {
		return err
	}
	
	_, err = r.CreateReferalIndexModel()
	if err != nil {
		return err
	}
	
	_, err = r.CreateVolumnIndexModel()
	if err != nil {
		return err
	}

	return nil
}
