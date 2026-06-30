package task

import (
	"testing"

	"github.com/Chibas/TaskTrackerGo/src/storage"
)

type stubStorage struct {
	readCalled  bool
	readData    []storage.Task
	readErr     error
	writeCalled bool
	writeData   any
}

func (s *stubStorage) ReadStorage() ([]storage.Task, error) {
	s.readCalled = true
	return s.readData, s.readErr
}

func (s *stubStorage) WriteStorage(data any) error {
	s.writeCalled = true
	s.writeData = data
	return nil
}

func TestTaskService(t *testing.T) {
	t.Run("list calls ReadStorage", func(t *testing.T) {
		stub := &stubStorage{}
		svc := NewTaskService(stub)

		_, err := svc.List("")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !stub.readCalled {
			t.Fatalf("Expected storage read to get called, called none")
		}
	})

	t.Run("list filters by status", func(t *testing.T) {
		stub := &stubStorage{
			readData: []storage.Task{
				{ID: 1, Description: "task one", Status: string(storage.Todo)},
				{ID: 2, Description: "task two", Status: string(storage.Done)},
			},
		}
		svc := NewTaskService(stub)

		result, err := svc.List("done")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(result) != 1 || result[0].ID != 2 || result[0].Status != string(storage.Done) {
			t.Fatalf("Expected only done task, got %#v", result)
		}
	})

	t.Run("AddTask returns error for duplicate description", func(t *testing.T) {
		stub := &stubStorage{
			readData: []storage.Task{
				{ID: 1, Description: "Write some code", Status: string(storage.Todo)},
			},
		}
		svc := NewTaskService(stub)

		err := svc.AddTask("Write some code")
		if err == nil {
			t.Fatal("Expected duplicate add to return error")
		}
		if stub.writeCalled {
			t.Fatal("Expected storage write not to be called on duplicate")
		}
	})

	t.Run("AddTask writes new task with next ID", func(t *testing.T) {
		stub := &stubStorage{
			readData: []storage.Task{{ID: 1, Description: "existing task", Status: string(storage.Todo)}},
		}
		svc := NewTaskService(stub)

		err := svc.AddTask("Write some code")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !stub.writeCalled {
			t.Fatal("Expected storage write to get called")
		}
		written, ok := stub.writeData.([]storage.Task)
		if !ok || len(written) != 2 {
			t.Fatalf("Expected two tasks written, got %#v", stub.writeData)
		}
		if written[1].ID != 2 || written[1].Description != "Write some code" {
			t.Fatalf("Expected appended task with ID 2, got %#v", written[1])
		}
	})

	t.Run("UpdateTask updates description and writes storage", func(t *testing.T) {
		stub := &stubStorage{
			readData: []storage.Task{{ID: 1, Description: "old", Status: string(storage.Todo)}},
		}
		svc := NewTaskService(stub)

		err := svc.UpdateTask(1, "new")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !stub.writeCalled {
			t.Fatal("Expected storage write to be called")
		}
		written, ok := stub.writeData.([]storage.Task)
		if !ok || len(written) != 1 {
			t.Fatalf("Expected one task written, got %#v", stub.writeData)
		}
		if written[0].Description != "new" {
			t.Fatalf("Expected description to be updated, got %#v", written[0].Description)
		}
	})

	t.Run("DeleteTask removes task and writes storage", func(t *testing.T) {
		stub := &stubStorage{
			readData: []storage.Task{{ID: 1, Description: "task one", Status: string(storage.Todo)}},
		}
		svc := NewTaskService(stub)

		err := svc.DeleteTask(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !stub.writeCalled {
			t.Fatal("Expected storage write to be called")
		}
		written, ok := stub.writeData.([]storage.Task)
		if !ok || len(written) != 0 {
			t.Fatalf("Expected zero tasks written, got %#v", stub.writeData)
		}
	})

	t.Run("MarkDone updates status and writes storage", func(t *testing.T) {
		stub := &stubStorage{
			readData: []storage.Task{{ID: 1, Description: "task one", Status: string(storage.Todo)}},
		}
		svc := NewTaskService(stub)

		err := svc.MarkDone(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !stub.writeCalled {
			t.Fatal("Expected storage write to be called")
		}
		written, ok := stub.writeData.([]storage.Task)
		if !ok || len(written) != 1 {
			t.Fatalf("Expected one task written, got %#v", stub.writeData)
		}
		if written[0].Status != string(storage.Done) {
			t.Fatalf("Expected status to be done, got %#v", written[0].Status)
		}
	})

	t.Run("MarkInProgress updates status and writes storage", func(t *testing.T) {
		stub := &stubStorage{
			readData: []storage.Task{{ID: 1, Description: "task one", Status: string(storage.Todo)}},
		}
		svc := NewTaskService(stub)

		err := svc.MarkInProgress(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !stub.writeCalled {
			t.Fatal("Expected storage write to be called")
		}
		written, ok := stub.writeData.([]storage.Task)
		if !ok || len(written) != 1 {
			t.Fatalf("Expected one task written, got %#v", stub.writeData)
		}
		if written[0].Status != string(storage.InProgress) {
			t.Fatalf("Expected status to be in-progress, got %#v", written[0].Status)
		}
	})
}
