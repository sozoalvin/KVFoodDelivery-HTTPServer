package main

import "fmt"

func PrintNoOfTriesExceeded() {
	fmt.Println("\nNumber of tries Exceeded. Please approach our staff for further assistance")
}

func PrintWelcomeMessage() {
	fmt.Println("\n================================================")
	fmt.Println("\nWelcome to Kay Kafe's Backend Ordering System!")
	fmt.Println("\n================================================")

}

func PrintUserValidated(s string) {
	fmt.Printf("\nYour username: %s has been validated. Please start ordering!\n\nPlease let us know if you have any questions with the ordering system. - Kay Kafe\n\n", s)
}
func PrintUserNotValidated(s string) {
	fmt.Printf("\nUsername %s not found. Please enter your username again.\n\nIf you think this is a mistake, please let our staff know!\n", s)
}
