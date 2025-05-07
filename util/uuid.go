package util

import "github.com/google/uuid"

//GenerateUUIDV4 generate a randomized uuid
func GenerateUUIDV4() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
