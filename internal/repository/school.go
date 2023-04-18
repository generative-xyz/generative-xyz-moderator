package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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
	file := &entity.Files{}
	err := r.Find(context.Background(), file.TableName(), bson.M{"uuid": id}, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}
