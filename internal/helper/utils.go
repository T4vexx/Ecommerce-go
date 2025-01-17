package helper

import (
	"crypto/rand"
)

func RandomNumbers(length int) (string, error) {
	const numbers = "0123456789"

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	numLength := len(buffer)

	for i := 0; i < length; i++ {
		buffer[i] = numbers[int(buffer[i])%numLength]
	}

	return string(buffer), nil
}
