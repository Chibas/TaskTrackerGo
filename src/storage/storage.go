package storage

import (
	"encoding/json"
	"io"
	"os"
)

type Storage interface {
	ReadStorage() ([]Task, error)
	WriteStorage(data any) error
}

type storage struct {
	storageName string
}

func NewStorage(name string) Storage {
	return &storage{
		storageName: name,
	}
}

func (s *storage) ReadStorage() ([]Task, error) {
	err := checkFile(s.storageName)

	if err != nil {
		return nil, err
	}

	f, err := os.Open(s.storageName)

	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	defer f.Close()

	byteValue, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if len(byteValue) == 0 {
		return []Task{}, nil
	}

	var result []Task
	if err := json.Unmarshal(byteValue, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *storage) WriteStorage(data any) error {
	err := checkFile(s.storageName)
	if err != nil {
		return err
	}

	b, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(s.storageName, b, 0644)

	return nil
}

func checkFile(fileName string) error {
	_, err := os.Stat(fileName)

	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(fileName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
