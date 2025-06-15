package memory

import (
	"encoding/json"
	"fmt"
	"io"
	"nyxze/fayth/model"
	"os"
)

type Memory interface {
	Save(msg []model.Message)
	Load() ([]model.Message, error)
}
type FileMemory struct {
	fileName string
}

func NewFileMemory(name string) *FileMemory {
	fm := &FileMemory{
		fileName: name,
	}
	return fm
}

func (f *FileMemory) Save(messages []model.Message) {
	json, err := json.Marshal(&messages)
	if err != nil {
		fmt.Println("failed to save message", err)
		return
	}
	w := f.open()
	n, err := w.Write(json)
	if err != nil {
		fmt.Println("writing as failed", err)
		return
	}
	fmt.Printf("%d bytes written ", n)

}
func (f *FileMemory) open() io.ReadWriter {

	file, err := os.OpenFile(f.fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to create memory file", err)
		os.Exit(1)
	}

	return file

}
func (f *FileMemory) Load() ([]model.Message, error) {
	reader := f.open()
	var message []model.Message
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []model.Message{}, nil
	}
	err = json.Unmarshal(data, &message)
	return message, err
}
