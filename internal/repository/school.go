package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
)

func (r Repository) InsertAISChoolJob(data *entity.AISchoolJob) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetAISchoolJobByUUID(uuid string) (*entity.AISchoolJob, error) {
	job := &entity.AISchoolJob{}
	err := r.Find(context.Background(), job.TableName(), bson.M{"uuid": uuid}, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r Repository) GetAISchoolJobByCreator(address string, limit, offset int64) ([]entity.AISchoolJob, error) {
	jobs := []entity.AISchoolJob{}
	err := r.Find(context.Background(), entity.AISchoolJob{}.TableName(), bson.M{"created_by": address}, &jobs, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r Repository) UpdateAISchoolJob(job *entity.AISchoolJob) error {
	filter := bson.D{{"uuid", job.UUID}}
	_, err := r.UpdateOne(job.TableName(), filter, job)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetAISchoolJobByStatus(status []string) ([]entity.AISchoolJob, error) {
	jobs := []entity.AISchoolJob{}
	err := r.Find(context.Background(), entity.AISchoolJob{}.TableName(), bson.M{"status": bson.M{"$in": status}}, &jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r Repository) GetFileByUUID(id string) (*entity.Files, error) {
	file := []entity.Files{}
	err := r.Find(context.Background(), entity.Files{}.TableName(), bson.M{"uuid": id}, &file)
	if err != nil {
		return nil, err
	}
	if len(file) == 0 {
		return nil, errors.New("file not found")
	}
	return &file[0], nil
}

func (r Repository) GetAISchoolUnClearedJob(before int64) ([]entity.AISchoolJob, error) {
	jobs := []entity.AISchoolJob{}
	err := r.Find(context.Background(), entity.AISchoolJob{}.TableName(), bson.M{"cleared_at": 0, "created_at": bson.M{"$lte": before}}, &jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r Repository) GetPresetDatasetByID(id string) (*entity.AISchoolPresetDataset, error) {
	file := []entity.AISchoolPresetDataset{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.Find(context.Background(), entity.AISchoolPresetDataset{}.TableName(), bson.M{"_id": objID}, &file)
	if err != nil {
		return nil, err
	}
	if len(file) == 0 {
		return nil, errors.New("file not found")
	}
	return &file[0], nil
}

func (r Repository) FindPresetDatasetByName(name, creator string) ([]entity.AISchoolPresetDataset, error) {
	files := []entity.AISchoolPresetDataset{}
	filter := bson.D{
		{"name", bson.D{{"$regex", primitive.Regex{Pattern: "a", Options: "i"}}}},
		{"$or",
			bson.A{
				bson.D{
					{"$and",
						bson.A{
							bson.D{{"creator", "0x3fc4053980c04ea4c517d82afbbb1cedbbbaa15b"}},
							bson.D{{"is_private", true}},
						},
					},
				},
				bson.D{
					{"$or",
						bson.A{
							bson.D{{"is_private", false}},
							bson.D{{"is_private", primitive.Null{}}},
						},
					},
				},
			},
		},
	}
	result, err := r.DB.Collection(entity.AISchoolPresetDataset{}.TableName()).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	err = result.All(context.Background(), &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r Repository) CreateDataset(data *entity.AISchoolPresetDataset) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}
