package main

import (
	"fmt"
	"log"
	"math/rand"
)

var guess = rand.Intn(10) + 1
var input int
var playAgain string

func doGuess() {
	fmt.Printf("Guess a number between 1 and 10: ")
	fmt.Scanln(&input)

	if input == guess {
		fmt.Println("You guessed correctly!")
	} else {
		fmt.Printf("You guessed wrong! Input: %d\n", input)

		if input > guess {
			fmt.Println("Guessed greater")
		} else {
			fmt.Println("Guessed smaller")
		}

		fmt.Print("Do you want to play again? (Y/n): ")
		fmt.Scanln(&playAgain)

		if playAgain == "y" || playAgain == "Y" || playAgain == "" {
			doGuess()
		} else if playAgain == "n" || playAgain == "N" {
			fmt.Printf("The correct number is %d", guess)
		} else {
			fmt.Println("Invalid input")
		}
	}
}

func main() {
	doGuess()
	log.Println("Game over")
}
