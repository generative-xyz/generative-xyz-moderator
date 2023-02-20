package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) GetInscriptionEventsByInscriptionID(id string, limit, offset int64) (*entity.InscriptionEvent, error) {
	resp := &entity.InscriptionEvent{}

	f := bson.D{
		{Key: "id", Value: id},
	}

	inscribeInfo, err := r.FilterOne(utils.INSCRIBE_INFO, f, &options.FindOneOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(inscribeInfo, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) CreateInscriptionEvent(event *entity.InscriptionEvent) error {
	err := r.InsertOne(event.TableName(), event)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) CreateInscriptionEvents(events []entity.InscriptionEvent) error {
	_events := make([]entity.IEntity, 0, len(events))
	for _, tokenHolder := range events {
		_event := tokenHolder
		_events = append(_events, &_event)
	}

	err := r.InsertMany(utils.INSCRIPTION_EVENT, _events)
	if err != nil {
		return err
	}

	return nil
}
