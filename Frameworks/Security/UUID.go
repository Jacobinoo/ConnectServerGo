package Security

import "github.com/google/uuid"

func GenerateUUID() (string, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
