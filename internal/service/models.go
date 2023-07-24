package service

import "fmt"

type Repository interface {
	GetAllNote() ([]JsonModel, error)
	AddNote(id string, jmodel []byte) error
}

type IServices interface {
}

type Event struct {
}

type Model struct {
}

type MyError struct {
	Message string
	Code    int
}

func (err *MyError) Error() string {
	return fmt.Sprintf("Code: %d\nMessage: %s\n", err.Code, err.Message)
}
