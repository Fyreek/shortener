package security

import (
	"errors"

	"github.com/google/uuid"
)

// ErrEmptyUUID represents a uuid that is empty
var ErrEmptyUUID = errors.New("Empty uuid provided")

// GetUUIDString generates a new uuid and returns a string representation of it
func GetUUIDString() string {
	return uuid.New().String()
}

// ParseUUIDString parses the provided string as a uuid. Returns an empty string + error if string was not uuid, and the uuid string + no error if it was a uuid
func ParseUUIDString(uuidString string) (string, error) {
	if uuidString == "" {
		return "", ErrEmptyUUID
	}
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
