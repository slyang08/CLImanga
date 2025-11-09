package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
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

	mangasFound, err := manga.FetchMangasByNameSearch(mangaName)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	if len(mangasFound) == 0 {
		return "", fmt.Errorf("Sorry no mangas found with that name")
	}

	prompt := promptui.Select{
		Label: "Select a Manga",
		Items: mangasFound,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "▶ {{ .Name | cyan }}",
			Inactive: "  {{ .Name }}",
			Selected: "✔ {{ .Name | green }}",
		},
		Size: 10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("selection cancelled or failed: %w", err)
	}

	selectedManga := mangasFound[i]
	return selectedManga.ID, nil
}

func selectChapterFromManga(mangaID *string) (*manga.ChapterSelect, error) {
	fmt.Println("Select a Chapter from the list...")

	chapterList, err := manga.GetAllChapterListOfManga(mangaID)
	if err != nil {
		fmt.Println(err.Error())
	}

	prompt := promptui.Select{
		Label: "Select a Chapter",
		Items: chapterList,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "▶ Chapter {{ .ChapterNumber }}: {{ .Title | cyan }}",
			Inactive: "  Chapter {{ .ChapterNumber }}: {{ .Title }}",
			Selected: "✔ Chapter {{ .ChapterNumber }}: {{ .Title | green }}",
		},
		Size: 10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("chapter selection cancelled or failed: %w", err)
	}

	selectedChapter := chapterList[i]
	return &selectedChapter, nil
}

func checkForDependencies() bool {
	return true
}
