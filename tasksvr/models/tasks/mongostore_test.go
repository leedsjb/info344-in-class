package tasks

import "testing"
import "gopkg.in/mgo.v2"
import "fmt"

func TestCRUD(t *testing.T) { // can create functions that start w/ Test in a go file ending in _test
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		t.Fatalf("error dialing Mongo: %v", err) // Fatalf ends test execution
	}
	defer sess.Close() // executes when function ends

	store := &MongoStore{
		Session:        sess,
		DatabaseName:   "test",
		CollectionName: "tasks",
	}

	newtask := &NewTask{
		Title: "Learn MongoDB",
		Tags:  []string{"mongo", "info344"}, // *** STATIC INITIALIZER
	}

	task, err := store.Insert(newtask)
	if err != nil {
		t.Errorf("error inserting new task: %v", err)
	}
	fmt.Println(task.ID)

	task2, err := store.Get(task.ID)
	if err != nil {
		t.Errorf("Error fetching task: %v", err)
	}
	if task2.Title != task.Title {
		t.Errorf("Task title did not match. Expected %s, but received %s", task.Title, task2.Title)
	}

	// cleanup now that test complete
	sess.DB(store.DatabaseName).C(store.CollectionName).RemoveAll(nil)
}
