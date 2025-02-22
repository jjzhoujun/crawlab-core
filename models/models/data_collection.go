package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataCollection struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

func (dc *DataCollection) GetId() (id primitive.ObjectID) {
	return dc.Id
}

func (dc *DataCollection) SetId(id primitive.ObjectID) {
	dc.Id = id
}

type DataCollectionList []DataCollection

func (l *DataCollectionList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
