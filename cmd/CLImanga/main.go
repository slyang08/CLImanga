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

		mangaID, err := selectManga(&mangaName)
		if err != nil {
			fmt.Println(err)
		}

		selectChapterFromManga(&mangaID)
	}
}

func selectManga(mangaName *string) (string, error) {
	fmt.Println("Searching mangas with Name:", *mangaName)

	var mangasFound map[string]string
	var err error

	mangasFound, err = manga.FetchMangasByNameSearch(mangaName)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	if len(mangasFound) == 0 {
		return "", fmt.Errorf("Sorry no mangas found with that name")
	}

	fmt.Println("Mangas found: ")

	var counter uint8 = 1

	for id, mangaName := range mangasFound {
		fmt.Printf("%v. %v \n", counter, mangaName)
		return id, nil
		counter++
	}
	return "", nil
}

func selectChapterFromManga(mangaID *string) {
	fmt.Println("Select a Chapter from the list..." + *mangaID)

	var chapterList map[string]map[string]any
	var err error

	chapterList, err = manga.GetAllChapterListOfManga(mangaID)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(chapterList)
}

func checkForDependencies() bool {
	return true
}
