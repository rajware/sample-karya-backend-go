package gormrepo

import (
	"errors"
	"log"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks"
	"gorm.io/gorm"
)

type GORMTaskRepository struct {
	db *gorm.DB
}

func (gtr *GORMTaskRepository) GetAll() ([]*tasks.Task, error) {
	var results []*tasks.Task
	err := gtr.db.Find(&results).Error
	if err != nil {
		return results, err
	}

	return results, nil
}

func (gtr *GORMTaskRepository) GetById(id uint) (*tasks.Task, error) {
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

func (gtr *GORMTaskRepository) Add(tsk *tasks.Task) (*tasks.Task, error) {
	err := gtr.db.Create(tsk).Error
	if err != nil {
		return nil, err
	}

	return tsk, nil
}

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

func (gtr *GORMTaskRepository) DeleteById(id uint) (*tasks.Task, error) {
	deltask, err := gtr.GetById(id)
	if err != nil {
		return nil, err
	}
	err = gtr.db.Delete(tasks.Task{ID: id}).Error
	if err != nil {
		return nil, err
	}

	return deltask, nil
}

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
