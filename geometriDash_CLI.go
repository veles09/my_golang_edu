package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

var lowerObsracle = "______🐍______"
var upperObsracle = "______🦇______"

func main() {
	score := 0
	greetings()
	for step := 1; step <= 10; step++ {
		obsracle := obsracle()
		input := inputStd()
		if obsracle == upperObsracle && input == "s" {
			score += 1
			fmt.Println("Ваш счёт", score)
		}
		if obsracle == lowerObsracle && input == "w" {
			score += 1
			fmt.Println("Ваш счёт", score)
		}
	}
	fmt.Println("Ваш счёт", score)
}

func greetings() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите имя: ")
	input.Scan()
	name := input.Text()
	fmt.Println("Привет", name)
	fmt.Println("Управление w - прыжок, s - приседание. У тебя десять попыток")
}

func obsracle() string {
	people := []string{lowerObsracle, upperObsracle}
	chosen := people[rand.Intn(len(people))]
	fmt.Print(chosen)
	return chosen
}

func inputStd() string {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("Делай выбор")
	input.Scan()
	command := input.Text()
	return command
}
