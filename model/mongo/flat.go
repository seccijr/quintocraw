package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/seccijr/quintocrawl/model"
	"math"
)

type MFlatRepo struct {
	fc mgo.Collection
}

func (mFlat MFlatRepo) FindByAddress(address string, skip int, limit int, cb model.FlatReact) error {
	iter := mFlat.fc.Find(bson.M{"address": address}).Skip(skip).Limit(limit).Iter()
	flat := &model.Flat{}

	for iter.Next(flat) {
		result := cb(flat)
		if result != nil {
			return result
		}
	}

	return nil
}

func (mFlat MFlatRepo) FindAllByAddress(address string, cb model.FlatReact) error {
	return mFlat.FindByAddress(address, 0, math.MaxInt64, cb)
}

func (mFlat MFlatRepo) Save(flat model.Flat) error {
	_, err := mFlat.fc.Upsert(bson.M{"ref": flat.Ref}, flat)

	return err
}