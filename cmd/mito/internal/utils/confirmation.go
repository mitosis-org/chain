package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfirmAction prompts the user to confirm an action
func ConfirmAction(skipConfirmation bool, message string) bool {
	if skipConfirmation {
		return true
	}

	fmt.Printf("%s\nType 'yes' to continue: ", message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return strings.ToLower(input) == "yes"
}
