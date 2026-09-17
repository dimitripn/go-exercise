package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"example.com/greetings"
)

func main() {
	name, age := "", 0

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("What is your name? ")
		scanner.Scan()

		name = scanner.Text()

		if name == "" {
			fmt.Println("Your name is empty")
			continue
		}

		break
	}

	for {
		fmt.Print("What is your age? ")
		scanner.Scan()

		input := scanner.Text()

		var err error
		age, err = strconv.Atoi(input)

		if err != nil {
			fmt.Println("Age must be an integer!")
			continue
		}

		break
	}

	message, err := greetings.Hello(name, age)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
}
