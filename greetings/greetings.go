package greetings

import (
	"errors"
	"fmt"

	"math/rand"
)

func Hello(name string, age int) (string, error) {

	if name == "" {
		return "", errors.New("No name")
	}

	message := fmt.Sprintf(randomFormat(), name, age)
	return message, nil
}

func randomFormat() string {
	formats := []string{
		"Hi, %v. Welcome! You're %d year old",
		"Great to see you, %v! You're %d year old",
		"Hail, %v! Well met! You're %d year old",
	}

	return formats[rand.Intn(len(formats))]
}
