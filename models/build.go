package models

import (
	"fmt"
	"time"

	"github.com/gernest/wuxia/db"
	uuid "github.com/satori/go.uuid"
)

//BuildArtifact is an interface which defines a buildable command.
type BuildArtifact interface {
	User() string
	Project() string
	Source() string
}

//BuildTask is a task for building a project.
type BuildTask struct {
	ID        int64     `store:"id"`
	UUID      string    `store:"uuid"`
	Done      bool      `store:"done"`
	User      string    `store:"user"`
	Project   string    `store:"project"`
	Source    string    `store:"source"`
	CreatedAt time.Time `store:"created_on"`
	UpdateAt  time.Time `store:"updated_on"`
}

//CreateBuildTask creates a new task record and stores it into the database.
func CreateBuildTask(store *db.DB, t *BuildTask) (*BuildTask, error) {
	var query = `
	BEGIN TRANSACTION;
	  INSERT INTO %s VALUES ($1,$2,$3,$4,$5, $6);
	COMMIT;
	`
	query = fmt.Sprintf(query, TaskTable)
	tx, err := store.Begin()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	uid := uuid.NewV4()
	_, err = tx.Exec(query,
		uid.String(),
		t.Done,
		t.User,
		t.Project,
		t.Source,
		now, now)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	t.UUID = uid.String()
	t.CreatedAt = now
	t.UpdateAt = now
	return t, nil
}