package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/scinac/CLImanga/internal/manga"
)

func main() {
	fmt.Println("Welcome to CLImanga!")
	fmt.Println("Checking Dependecies... ")
	if !checkForDependencies() {
		fmt.Println("Missing Dependencies")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	if len(os.Args) < 2 {
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

	var mangasFound map[string]string
	var err error

	mangasFound, err = manga.FetchMangaNames(mangaName)
	if err != nil {
		fmt.Print(err.Error())
	}

	if len(mangasFound) == 0 {
		fmt.Println("Sorry no mangas found with that name")
		return
	}

	fmt.Println("Mangas found: ")
	for _, manga := range mangasFound {
		fmt.Printf("%v \n", manga)
	}
}

func checkForDependencies() bool {
	return true
}
