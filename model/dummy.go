package model

import (
    "fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DummyDAO struct {
	conn *Conn
}

func NewDummyDAO(_conn *Conn) *DummyDAO {
	if nil == _conn {
		return &DummyDAO{
			conn: defaultConn,
		}
	} else {
		return &DummyDAO{
			conn: _conn,
		}
	}
}

func (this *DummyDAO) Filter(_offset int64, _count int64, _tag []string) (total_ int64, dummy_ []string, err_ error) {
    total_ = 0
    dummy_ = make([]string, 0)
    err_ = nil

    total_ = 0
	ctx, cancel := NewContext()
	defer cancel()

    filter := []bson.M{
        bson.M{"$match": bson.M{"name": bson.M{"$in":_tag}}},
    }
    opts := options.Distinct()
	result, err := this.conn.DB.Collection(CollectionName).Distinct(ctx, "dummy", filter, opts)
    if nil != err {
        err_ = err
        return
    }
    fmt.Println(result)

    return
}

func (this *DummyDAO) List(_offset int64, _count int64, _owner string) (total_ int64, tag_ []*Tag, err_ error) {
    total_ = 0
    tag_ = make([]*Tag, 0)
    err_ = nil

    total_ = 0
	ctx, cancel := NewContext()
	defer cancel()

    pipeline := []bson.M{
        bson.M{"$match": bson.M{"dummy": bson.M{"$in":[]string{_owner}}}},
    }
    opts := options.Aggregate()
	cur, err := this.conn.DB.Collection(CollectionName).Aggregate(ctx, pipeline, opts)
    defer cur.Close(ctx)

    if nil != err {
        err_ = err
        return
    }

    for cur.Next(ctx) {
		var tag Tag
		err = cur.Decode(&tag)
        tag_ = append(tag_, &tag)
    }
    return
}
