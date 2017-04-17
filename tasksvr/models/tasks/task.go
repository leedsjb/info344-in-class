package tasks

import (
	"fmt"
	"time"
)

//NewTask represents a new task posted to the server
type NewTask struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

//Task represents a task stored in the database
type Task struct {
	ID         interface{} `json:"id" bson:"_id"` // interface{} is go's any type (implements an interface w/ no methods), //bson: MongoDB compatibility
	Title      string      `json:"title"`
	Tags       []string    `json:"tags"`
	CreatedAt  time.Time   `json:"createdAt"`
	ModifiedAt time.Time   `json:"modifiedAt"`
	Complete   bool        `json:"complete"`
}

//Validate will validate the NewTask
func (nt *NewTask) Validate() error { // (nt *NewTask) = received, similar to Java this pointer but named. Struct upon which validate function is called
	//Title field must be non-zero in length
	if len(nt.Title) == 0 {
		return fmt.Errorf("title must be something")
	}
	return nil
}

//ToTask converts a NewTask to a Task
func (nt *NewTask) ToTask() *Task {
	t := &Task{ // could also use new keyword to create new struct
		// static initializer
		Title:      nt.Title,
		Tags:       nt.Tags,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	return t
}
