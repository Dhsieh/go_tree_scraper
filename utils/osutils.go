package utils

import (
	"fmt"
	"os"
)

func CreateDirectory(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		fmt.Printf("Directory %s has not been created yet. Creating it!\n", dirName)
		err = os.MkdirAll(dirName, os.ModePerm)

		if err != nil {
			fmt.Printf("Could not create directory %s!\n", dirName)
			fmt.Println(err)
		}
	} else {
		fmt.Printf("Directory %s has already been created\n", dirName)
	}
}
