package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"strings"
)

func (r *Repository) InsertAction(obj *entity.Actions) error {
	return r.InsertOne(obj.TableName(), obj)
}

func (r *Repository) FilterAction(filter *entity.FilterActions) bson.D {
	f := bson.D{}
	if filter.ObjectID != nil && *filter.ObjectID != "" {
		f = append(f, bson.E{"object_id", strings.ToLower(*filter.ObjectID)})
	}

	if filter.Parent != nil && *filter.Parent != "" {
		f = append(f, bson.E{"parent", strings.ToLower(*filter.Parent)})
	}

	if filter.ObjectType != nil {
		f = append(f, bson.E{"object_type", *filter.ObjectType})
	}

	if filter.Action != nil {
		f = append(f, bson.E{"action", filter.Action})
	}

	return f
}
