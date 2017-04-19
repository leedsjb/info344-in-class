package tasks

import (
	"time"

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

	if sID, ok := ID.(string); ok { // REVIEW HOW THIS SYNTAX WORKS ***
		ID = bson.ObjectIdHex(sID) // converts string object ID to bson ID via hex conversion
	}

	task := &Task{} // essentially, task = new Task()
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(ID).One(task)
	return task, err
}

func (ms *MongoStore) GetAll() ([]*Task, error) {
	tasks := []*Task{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(nil).All(&tasks) // find all documents in collection in MongoDB

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

//
func (ms *MongoStore) Update(task *Task) error {

	// may need to tesk incoming task.ID to see if it is a string using a go lang assertion
	if sID, ok := task.ID.(string); ok {
		task.ID = bson.ObjectIdHex(sID)
	}

	task.ModifiedAt = time.Now() // update the time the task was modified

	col := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName) // store reference to MongoDB collection

	updates := bson.M{"$set": bson.M{"complete": task.Complete, "modifiedat": task.ModifiedAt}}

	return col.UpdateId(task.ID, updates)
}

// $set: {
//        item: "ABC123",
//        "info.publisher": "2222",
//        tags: [ "software" ],
//        "ratings.1": { by: "xyz", rating: 3 }
//      }
