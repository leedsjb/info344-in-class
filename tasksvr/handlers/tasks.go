package handlers

import (
	"encoding/json"
	"net/http"
	"path"

	"fmt"

	"github.com/leedsjb/info344-in-class/tasksvr/models/tasks"
)

//HandleTasks will handle requests for the /v1/tasks resource
func (ctx *Context) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body) // read JSON from POST body
		newtask := &tasks.NewTask{}
		if err := decoder.Decode(newtask); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err := newtask.Validate(); err != nil {
			http.Error(w, "error validationg task: "+err.Error(), http.StatusBadRequest)
			return
		}

		task, err := ctx.TasksStore.Insert(newtask)
		if err != nil {
			http.Error(w, "error inserting task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(task)

	// get all tasks in the database when http://localhost:4000/v1/tasks endpoint is called
	case "GET":
		tasks, err := ctx.TasksStore.GetAll()
		if err != nil {
			http.Error(w, "error getting tasks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(tasks)
	}
}

//HandleSpecificTask will handle requests for the /v1/tasks/some-task-id resource
// Handles get requests in the form of http://localhost:4000/v1/tasks/[task_id]
func (ctx *Context) HandleSpecificTask(w http.ResponseWriter, r *http.Request) {

	_, id := path.Split(r.URL.Path) // returns directory and file. unique id of task is the file portion, preceding URL is the directory

	switch r.Method {
	case "GET":

		task, err := ctx.TasksStore.Get(id) // retrieve task by ID

		if err != nil {
			http.Error(w, "error finding task: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(task)

	case "PATCH":

		decoder := json.NewDecoder(r.Body) // read JSON from PATCH body

		task := &tasks.Task{}

		if err := decoder.Decode(task); err != nil {
			http.Error(w, "error decoding JSON"+err.Error(), http.StatusBadRequest)
			return
		}

		task.ID = id

		fmt.Println(task.ID)

		if err := ctx.TasksStore.Update(task); err != nil {
			http.Error(w, "error updating: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("update successful"))

	}

	// mgo.ErrNotFound

}
