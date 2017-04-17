package tasks

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	Session        *mgo.Session
	DatabaseName   string
	CollectionName string
}

func (ms *MongoStore) Insert(newtask *NewTask) (*Task, error) { // function executes on an existing MongoStore struct
	t := newtask.ToTask()     // create new task
	t.ID = bson.NewObjectId() // unique to time and machine created on
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Insert(t)
	return t, err // do not need to handle error, client code does
}

func (ms *MongoStore) Get(ID interface{}) (*Task, error) {
	task := &Task{} // essentially, task = new Task()
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(ID).One(task)
	return task, err
}
