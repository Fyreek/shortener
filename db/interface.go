package db

import (
	"errors"
)

var ErrNoDocument = errors.New("No documents found")

type Database interface {
	Connect(ip string, port, timeout int) error
	IsConnected() bool
	SetDatabase(name string)
	GetSingleEntry(collection, column, value string, iStruct interface{}) error
	GetMultipleEntries(collection, column, value string, sort map[string]interface{}, limit int) ([][]byte, error)
	InsertSingleEntry(collection string, value interface{}) error
	UpdateSingleEntry(collection, filterColumn, filterValue string, values interface{}) error
	DeleteSingleEntry(collection, filterColumn, filterValue string) error
}
