package mapper

import (
	"testing"

	"github.com/Chibas/TaskTrackerGo/src/storage"
)

type stubTaskService struct {
	addCalled            bool
	addDescription       string
	updateCalled         bool
	updateID             int
	updateDescription    string
	deleteCalled         bool
	deleteID             int
	markDoneCalled       bool
	markDoneID           int
	markInProgressCalled bool
	markInProgressID     int
	listCalls            []string
	listTasks            []storage.Task
	listErr              error
}

func (s *stubTaskService) AddTask(description string) error {
	s.addCalled = true
	s.addDescription = description
	return nil
}
func (s *stubTaskService) UpdateTask(id int, description string) error {
	s.updateCalled = true
	s.updateID = id
	s.updateDescription = description
	return nil
}
func (s *stubTaskService) DeleteTask(id int) error {
	s.deleteCalled = true
	s.deleteID = id
	return nil
}
func (s *stubTaskService) MarkDone(id int) error {
	s.markDoneCalled = true
	s.markDoneID = id
	return nil
}
func (s *stubTaskService) MarkInProgress(id int) error {
	s.markInProgressCalled = true
	s.markInProgressID = id
	return nil
}
func (s *stubTaskService) List(status string) ([]storage.Task, error) {
	s.listCalls = append(s.listCalls, status)
	return s.listTasks, s.listErr
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Command
		wantErr bool
	}{
		{name: "valid add", input: "add", want: Add},
		{name: "valid list", input: "list", want: List},
		{name: "invalid command", input: "unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := NewCommandMapper(&stubTaskService{})
			got, err := mapper.ParseCommand(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestExecuteCommand(t *testing.T) {
	t.Run("add command uses description", func(t *testing.T) {
		service := &stubTaskService{}
		mapper := NewCommandMapper(service)

		err := mapper.ExecuteCommand(Add, "buy milk", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !service.addCalled || service.addDescription != "buy milk" {
			t.Fatalf("expected add task to be called with description %q", "buy milk")
		}
	})

	t.Run("update command uses valid id and description", func(t *testing.T) {
		service := &stubTaskService{}
		mapper := NewCommandMapper(service)

		err := mapper.ExecuteCommand(Update, "7", "new description")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !service.updateCalled || service.updateID != 7 || service.updateDescription != "new description" {
			t.Fatalf("expected update task to be called with id 7 and description %q", "new description")
		}
	})

	t.Run("delete command uses id", func(t *testing.T) {
		service := &stubTaskService{}
		mapper := NewCommandMapper(service)

		err := mapper.ExecuteCommand(Delete, "3", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !service.deleteCalled || service.deleteID != 3 {
			t.Fatalf("expected delete task to be called with id 3")
		}
	})

	t.Run("mark in progress command uses id", func(t *testing.T) {
		service := &stubTaskService{}
		mapper := NewCommandMapper(service)

		err := mapper.ExecuteCommand(MarkInProgress, "4", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !service.markInProgressCalled || service.markInProgressID != 4 {
			t.Fatalf("expected mark in progress to be called with id 4")
		}
	})

	t.Run("mark done command uses id", func(t *testing.T) {
		service := &stubTaskService{}
		mapper := NewCommandMapper(service)

		err := mapper.ExecuteCommand(MarkDone, "5", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !service.markDoneCalled || service.markDoneID != 5 {
			t.Fatalf("expected mark done to be called with id 5")
		}
	})

	t.Run("list command passes filter to service", func(t *testing.T) {
		service := &stubTaskService{listTasks: []storage.Task{{ID: 1, Description: "task", Status: "done"}}}
		mapper := NewCommandMapper(service)

		err := mapper.ExecuteCommand(List, "done", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(service.listCalls) != 1 || service.listCalls[0] != "done" {
			t.Fatalf("expected list to be called with filter 'done', got %v", service.listCalls)
		}
	})

	t.Run("list uses empty filter when arg is missing", func(t *testing.T) {
		service := &stubTaskService{}
		mapper := NewCommandMapper(service)

		err := mapper.ExecuteCommand(List, 42, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(service.listCalls) != 1 || service.listCalls[0] != "" {
			t.Fatalf("expected empty filter, got %v", service.listCalls)
		}
	})

	t.Run("invalid arguments return error", func(t *testing.T) {
		mapper := NewCommandMapper(&stubTaskService{})
		err := mapper.ExecuteCommand(Add, 123, nil)
		if err == nil {
			t.Fatal("expected an error for invalid add arguments")
		}
	})

	t.Run("unknown command returns error", func(t *testing.T) {
		mapper := NewCommandMapper(&stubTaskService{})
		err := mapper.ExecuteCommand(Command("unknown"), nil, nil)
		if err == nil {
			t.Fatal("expected an error for unknown command")
		}
	})
}
