package internal

import (
	"fmt"
	"os"
)

func HandleError(msg string, e error) {
	errorMessage := fmt.Sprintf("Error: %s \n Cause: %s", msg, e.Error())
	fmt.Println(errorMessage)
	os.Exit(1)
}
