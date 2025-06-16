package service

import (
	"dummy-backend/lib/domain"
	"dummy-backend/lib/repository"
	"errors"
)

type TaskService interface {
	CreateTask(req *domain.CreateTaskRequest) (*domain.Task, error)
	GetAllTasks() ([]domain.Task, error)
	GetTaskByID(id uint) (*domain.Task, error)
	UpdateTask(id uint, req *domain.UpdateTaskRequest) (*domain.Task, error)
	DeleteTask(id uint) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) CreateTask(req *domain.CreateTaskRequest) (*domain.Task, error) {
	task := &domain.Task{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	err := s.taskRepo.Create(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetAllTasks() ([]domain.Task, error) {
	return s.taskRepo.GetAll()
}

func (s *taskService) GetTaskByID(id uint) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (s *taskService) UpdateTask(id uint, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	existingTask, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("task not found")
	}

	// Update only provided fields
	if req.Title != nil {
		existingTask.Title = *req.Title
	}
	if req.Description != nil {
		existingTask.Description = *req.Description
	}
	if req.Completed != nil {
		existingTask.Completed = *req.Completed
	}

	err = s.taskRepo.Update(id, existingTask)
	if err != nil {
		return nil, err
	}

	return existingTask, nil
}

func (s *taskService) DeleteTask(id uint) error {
	_, err := s.taskRepo.GetByID(id)
	if err != nil {
		return errors.New("task not found")
	}

	return s.taskRepo.Delete(id)
}
