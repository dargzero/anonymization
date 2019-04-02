package anondb

import "github.com/globalsign/mgo/bson"

const dbName = "anondb"
const dataPrefix = "anon_"

func FetchUnanonymizedData(dataset string) (data []bson.M, err error) {
	session := globalSession.Copy()
	defer session.Close()
	err = session.DB(dbName).
		C(dataPrefix + dataset).
		Find(bson.M{"__anonymized": false}).
		All(&data)
	return
}

func PersistAnonymizedData(dataset string, data []bson.M) (err error) {
	session := globalSession.Copy()
	defer session.Close()
	for _, doc := range data {
		err = session.DB(dbName).
			C(dataPrefix+dataset).
			UpdateId(doc["_id"], doc)
		if err != nil {
			return
		}
	}
	return

}
