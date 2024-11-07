package tools

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ChoseOneOf(message string, choices []string, defaultChoice int) int {
	var ask func() int
	ask = func() int {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Select: [default - ", defaultChoice, "]: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		if text == "" {
			return defaultChoice
		}
		chosen, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Invalid choice, please enter value from 0 to ", len(choices)-1)
			return ask()
		}
		if chosen < 0 || chosen >= len(choices) {
			fmt.Println("Invalid choice, please enter value from 0 to ", len(choices)-1)
			return ask()
		}
		return chosen
	}
	fmt.Println(message, ": ")
	for i, c := range choices {
		fmt.Println(i, ":", c)
	}
	return ask()
}

func Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message, "(y/n): ")
	text, _ := reader.ReadString('\n')
	if text[0] != 'y' && text[0] != 'Y' && text[0] != 'n' && text[0] != 'N' {
		return Confirm(message)
	}
	return text[0] == 'y' || text[0] == 'Y'
}

func AskNum(message string, from, to, defaultNum int) int {
	var ask func() int
	ask = func() int {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(message, ". Select one from", from, "to", to, "): [default - ", defaultNum, "]: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		if text == "" {
			return defaultNum
		}
		num, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Invalid number, please enter valid number")
			return ask()
		}
		if num < from || num > to {
			fmt.Println("Invalid number, please enter number from", from, "to", to)
			return ask()
		}
		return num
	}
	return ask()
}
