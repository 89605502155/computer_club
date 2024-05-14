package main

import (
	"computerClub/pkg/scaner"
)

func main() {
	reader := scaner.NewScaner("test_file.txt")
	reader.Scaner.Read()
	reader.Scaner.Close()
}
