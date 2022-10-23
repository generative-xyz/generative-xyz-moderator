package mongodb

import (
	"context"
	"reflect"
	"time"

	pg "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const MongoIdGenCollectionName = "id-gen"

type Repository interface {
	Database() *mongo.Database
	Collection() *mongo.Collection
	Filter(ctx context.Context, filters map[string]interface{}, sortFields []string, sortValues []int, page int64, limit int64, result interface{}) (int64, error)
	Find(ctx context.Context, filters map[string]interface{}, result interface{}, opts ...*options.FindOptions) error
	FindById(ctx context.Context, id primitive.ObjectID, value interface{}) error
	FindOne(ctx context.Context, filters map[string]interface{}, value interface{}, opts ...*options.FindOneOptions) error
	Update(ctx context.Context, model interface{}, id primitive.ObjectID, dateModified time.Time, opts ...*options.FindOneAndReplaceOptions) error
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error)
	Create(ctx context.Context, model interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error)
	CreateMany(ctx context.Context, models []interface{}, opts ...*options.InsertManyOptions) ([]primitive.ObjectID, error)
	CreateManyWithTransaction(ctx context.Context, models []interface{}, opts ...*options.TransactionOptions) ([]primitive.ObjectID, error)
	Delete(ctx context.Context, id []primitive.ObjectID, opts ...*options.DeleteOptions) error
	Aggregate(ctx context.Context, page int64, limit int64, result interface{}, agg ...interface{}) (int64, error)
	WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error)
	CreateCompoundIndex(ctx context.Context, compoundIndex []string, unique bool, opts ...*options.CreateIndexesOptions) error
	CreateIndices(ctx context.Context, indices []string, unique bool, opts ...*options.CreateIndexesOptions) error
	NextId(ctx context.Context, sequenceName *string) (uint, error)
}

type BaseRepository struct {
	CollectionName string
	DB             *mongo.Database
}

func (b *BaseRepository) Database() *mongo.Database {
	return b.DB
}
func (b *BaseRepository) Collection() *mongo.Collection {
	return b.DB.Collection(b.CollectionName)
}
func (b *BaseRepository) FindById(ctx context.Context, id primitive.ObjectID, value interface{}) error {
	return b.FindOne(ctx, bson.M{"_id": id}, value)
}
func (b *BaseRepository) FindOne(ctx context.Context, filters map[string]interface{}, value interface{}, opts ...*options.FindOneOptions) error {
	res := b.DB.Collection(b.CollectionName).FindOne(ctx, filters, opts...)
	if res.Err() != nil {
		return res.Err()
	}
	if err := res.Decode(value); err != nil {
		return err
	}
	return nil
}
func (b *BaseRepository) Find(ctx context.Context, filters map[string]interface{}, result interface{}, opts ...*options.FindOptions) error {
	cur, err := b.DB.Collection(b.CollectionName).Find(ctx, filters, opts...)
	if err != nil {
		return err
	}
	if err := cur.All(ctx, result); err != nil {
		return err
	}
	return nil
}
func (b *BaseRepository) Filter(ctx context.Context, filters map[string]interface{}, sortFields []string, sortValue []int, page int64, limit int64, result interface{}) (int64, error) {
	query := pg.New(b.DB.Collection(b.CollectionName)).Decode(result).Context(ctx).Page(page).Limit(limit)
	if len(sortFields) < 1 {
		query = query.Sort("date_modified", -1)
	} else {
		for i, sort := range sortFields {
			if i < len(sortValue) {
				query = query.Sort(sort, sortValue[i])
			}
		}
	}
	aggPaginatedData, err := query.Filter(filters).Find()
	if err != nil {
		return 0, err
	}
	return aggPaginatedData.Pagination.Total, err
}
func (b *BaseRepository) Create(ctx context.Context, model interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
	result, err := b.DB.Collection(b.CollectionName).InsertOne(ctx, model, opts...)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id, nil
	}
	return primitive.NilObjectID, nil
}
func (b *BaseRepository) CreateManyWithTransaction(ctx context.Context, models []interface{}, opts ...*options.TransactionOptions) ([]primitive.ObjectID, error) {
	session, err := b.DB.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		return b.DB.Collection(b.CollectionName).InsertMany(sessionContext, models)
	}
	results, err := session.WithTransaction(ctx, callback, opts...)
	if err != nil {
		return nil, err
	}
	var ids []primitive.ObjectID
	if insertResult, ok := results.(*mongo.InsertManyResult); ok {
		for _, result := range insertResult.InsertedIDs {
			if id, ok := result.(primitive.ObjectID); ok {
				ids = append(ids, id)
			}
		}
	}
	return ids, err
}
func (b *BaseRepository) Update(ctx context.Context, model interface{}, id primitive.ObjectID, dateModified time.Time, opts ...*options.FindOneAndReplaceOptions) error {
	res := b.DB.Collection(b.CollectionName).FindOneAndReplace(
		ctx,
		bson.D{
			{Key: "_id", Value: id},
			{Key: "date_modified", Value: dateModified},
		},
		model,
		opts...,
	)
	return res.Err()
}
func (b *BaseRepository) Delete(ctx context.Context, ids []primitive.ObjectID, opts ...*options.DeleteOptions) error {
	_, err := b.DB.Collection(b.CollectionName).DeleteMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
		opts...,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseRepository) CreateCompoundIndex(ctx context.Context, compoundIndex []string, unique bool, opts ...*options.CreateIndexesOptions) error {
	var indices bsonx.Doc
	for _, index := range compoundIndex {
		indices = append(indices, bsonx.Elem{
			Key:   index,
			Value: bsonx.Int32(1),
		})
	}
	index := mongo.IndexModel{
		Keys:    indices,
		Options: options.Index().SetUnique(unique),
	}
	_, err := b.DB.Collection(b.CollectionName).Indexes().CreateOne(ctx, index, opts...)
	return err
}

func (b *BaseRepository) Aggregate(ctx context.Context, page int64, limit int64, result interface{}, agg ...interface{}) (int64, error) {
	query := pg.New(b.DB.Collection(b.CollectionName)).Context(ctx).Page(page).Limit(limit)
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

func (b *BaseRepository) WithTransaction(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error) {
	session, err := b.DB.Client().StartSession()
	defer session.EndSession(ctx)
	if err != nil {
		return nil, err
	}
	return session.WithTransaction(ctx, callback, opts...)
}

func (b *BaseRepository) CreateIndices(ctx context.Context, indices []string, unique bool, opts ...*options.CreateIndexesOptions) error {
	var indexModels []mongo.IndexModel
	for _, index := range indices {
		indexModel := mongo.IndexModel{
			Keys: bsonx.Doc{{Key: index,
				Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(unique),
		}
		indexModels = append(indexModels, indexModel)
	}
	_, err := b.DB.Collection(b.CollectionName).Indexes().CreateMany(ctx, indexModels, opts...)
	return err
}

type Counter struct {
	ID  string `json:"id" bson:"_id"`
	Seq uint   `json:"seq" bson:"seq"`
}

func (b *BaseRepository) NextId(ctx context.Context, sequenceName *string) (uint, error) {
	findOptions := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	name := b.CollectionName
	if sequenceName != nil && len(*sequenceName) > 0 {
		name = *sequenceName
	}
	var counter Counter
	err := b.Database().Collection(MongoIdGenCollectionName).
		FindOneAndUpdate(ctx,
			bson.M{"_id": name},
			bson.M{"$inc": bson.M{"seq": 1}},
			findOptions,
		).Decode(&counter)

	if err != nil {
		return 0, err
	}
	return counter.Seq, nil
}
func (b *BaseRepository) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	result, err := b.DB.Collection(b.CollectionName).UpdateByID(ctx, id, update, opts...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
func (b *BaseRepository) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	result, err := b.DB.Collection(b.CollectionName).UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
func (b *BaseRepository) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	result, err := b.DB.Collection(b.CollectionName).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
func (b *BaseRepository) CreateMany(ctx context.Context, models []interface{}, opts ...*options.InsertManyOptions) ([]primitive.ObjectID, error) {
	results, err := b.DB.Collection(b.CollectionName).InsertMany(ctx, models, opts...)
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
