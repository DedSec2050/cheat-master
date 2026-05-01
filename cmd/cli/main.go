package main

import (
	"cheat-master/internal/orchestrator"
	"fmt"
)

func main() {
	var email, password string

	fmt.Print("Email: ")
	fmt.Scanln(&email)

	fmt.Print("Password: ")
	fmt.Scanln(&password)

	orchestrator.Run(email, password)
}

