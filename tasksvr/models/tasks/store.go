package tasks

//Store defines an abstract interface for a Task object store
type Store interface {
	//Insert inserts a NewTask and
	//returns the fully-populated Task or an error
	Insert(newtask *NewTask) (*Task, error)

	//GetID retrieves a new task and
	//returns the fully-populated Task or an error
	Get(ID interface{}) (*Task, error)

	//Get all tasks
	GetAll() ([]*Task, error)

	//
	Update(task *Task) error
}
