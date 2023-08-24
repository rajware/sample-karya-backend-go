package gormrepo

import (
	"errors"
	"log"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks"
	"gorm.io/gorm"
)

// GORMTaskRepository is an implementation of the tasks.TaskRepository
// interface.
type GORMTaskRepository struct {
	db *gorm.DB
}

// GetAll returns all Tasks, or an error.
func (gtr *GORMTaskRepository) GetAll() ([]*tasks.Task, error) {
	var results []*tasks.Task
	err := gtr.db.Find(&results).Error
	if err != nil {
		return results, err
	}

	return results, nil
}

// GetByID returns the Task with the given ID, or an error.
func (gtr *GORMTaskRepository) GetByID(id uint) (*tasks.Task, error) {
	var task tasks.Task
	err := gtr.db.First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, tasks.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Add creates a new Task in the data store. The store must provide
// an auto-generated ID. Any supplied ID will be ignored.
func (gtr *GORMTaskRepository) Add(tsk *tasks.Task) (*tasks.Task, error) {
	err := gtr.db.Create(tsk).Error
	if err != nil {
		return nil, err
	}

	return tsk, nil
}

// Update updates the given task in the data store, and
// returns the updated task or an error. The task must
// exist in the data store.
func (gtr *GORMTaskRepository) Update(tsk *tasks.Task) (*tasks.Task, error) {
	log.Printf("The task as received for update:%+v", tsk)
	eng := gtr.db.Model(tsk).Updates(map[string]interface{}{
		"Description": tsk.Description,
		"Deadline":    tsk.Deadline,
		"Completed":   tsk.Completed,
	})
	err := eng.Error
	if err != nil {
		return nil, err
	}
	if eng.RowsAffected == 0 {
		return nil, tasks.ErrNotUpdated
	}

	return tsk, nil
}

// DeleteByID deletes the Task with the given ID from the
// data store, and returns the deleted task or an error.
// The Task must exist in the data store.
func (gtr *GORMTaskRepository) DeleteByID(id uint) (*tasks.Task, error) {
	deltask, err := gtr.GetByID(id)
	if err != nil {
		return nil, err
	}
	err = gtr.db.Delete(tasks.Task{ID: id}).Error
	if err != nil {
		return nil, err
	}

	return deltask, nil
}

// New returns an instance of GORMTaskRepository. It requires
// a GORM dialector and options. The dialector has to be
// initialized before calling New.
func New(d gorm.Dialector, o gorm.Option) *GORMTaskRepository {
	newRepo := &GORMTaskRepository{
		db: nil,
	}

	db, err := gorm.Open(d, o)
	if err != nil {
		panic("failed to connect database")
	}

	newRepo.db = db
	db.AutoMigrate(&tasks.Task{})

	return newRepo
}
