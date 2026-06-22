package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Storage interface {
	ReadStorage() ([]Task, error)
	WriteStorage(data any)
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
	f, err := os.Open(s.storageName)
	fmt.Printf("Storage name %s \n", s.storageName)
	if err != nil {
		fmt.Printf("Storage read error %s \n", err)
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

	fmt.Println("Storage successfully read")

	fmt.Printf("Storage bytes %d", len(byteValue))

	if len(byteValue) == 0 {
		return []Task{}, nil
	}

	var result []Task
	if err := json.Unmarshal(byteValue, &result); err != nil {
		fmt.Printf("Error Unmarshalling %s", err)
		return nil, err
	}

	fmt.Println("Storage successfully parsed")

	return result, nil
}

func (s *storage) WriteStorage(data any) {
	b, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(s.storageName, b, 0644)
}
