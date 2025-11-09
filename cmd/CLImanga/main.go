package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome to CLImanga!")
	fmt.Println("Checking Dependecies... ")
	if !checkForDependencies() {
		fmt.Println("Missing Dependencies")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	if len(os.Args) == 0 {
		fmt.Print("Search Manga: ")
		mangaName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		searchMangas(&mangaName)
	}
}

func searchMangas(mangaName *string) {
	fmt.Println("Searching mangas with Name:", *mangaName)
}

func checkForDependencies() bool {
	return true
}
