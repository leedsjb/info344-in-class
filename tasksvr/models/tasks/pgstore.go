package tasks

import "database/sql"

type PGSTORE struct {
	DB *sql.DB
}

func (ps *PGSTORE) Insert(newtask *NewTask) (*Task, error) {

	// could use batch query and prepared statements

	t := newtask.ToTask() // convert task to a full task object
	tx, err := ps.DB.Begin()
	if err != nil {
		return nil, err
	}

	sql := `INSERT into tasks (title, createdAt, modifiedAt, complete) VALUES ($1, $2, $3, $4)
			RETURNING id` // parameter markers protect against SQL injection

	row := tx.QueryRow(sql, t.Title, t.CreatedAt, t.ModifiedAt, t.Complete)

	err = row.Scan(&t.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	sql = `INSERT INTO tags (taskID, tag) VALUES ($1, $2)`
	for _, tag := range t.Tags { // range == for each
		_, err := tx.Exec(sql, t.ID, tag)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return t, err

}

func (ps *PGSTORE) Get(ID interface{}) (*Task, error) {
	return nil, nil
}

func (ps *PGSTORE) GetAll() ([]*Task, error) {
	return nil, nil
}

func (ps *PGSTORE) Update(task *Task) error {
	return nil
}
