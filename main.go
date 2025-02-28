package main

import (
	"fmt"
	"log"
	"math/rand"
)

var guess = rand.Intn(10) + 1
var input int

func doGuess() {
	fmt.Printf("Guess a number between 1 and 10: ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("Invalid input")
		return
	}

	if input == guess {
		fmt.Println("You guessed correctly!")
	} else {
		fmt.Printf("You guessed wrong! ")

		if input > guess {
			fmt.Println("Guessed greater")
		} else {
			fmt.Println("Guessed smaller")
		}

		doGuess()
	}
}

func main() {
	doGuess()
	log.Println("Thanks for playing!")
}
