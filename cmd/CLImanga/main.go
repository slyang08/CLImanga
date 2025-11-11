package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/scinac/CLImanga/internal/manga"
)

const (
	APPMODE_DOWNLOAD   = 0
	APPMODE_READ       = 1
	DIRECTORY_DOWNLOAD = "downloads"
	DRECTORY_CACHE     = "cache"
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
		fmt.Println("Do you want to download a Manga or only read one?")
		appMode, err := selectAppMode()
		if err != nil {
			fmt.Print(err.Error())
			return
		}

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

		switch appMode {
		case APPMODE_DOWNLOAD:
		case APPMODE_READ:
			readMangaMode(&mangaID, &mangaName)
		}
	}
}

func readMangaMode(mangaID *string, mangaName *string) {
	selectedChapter, errSelect := selectChapterFromManga(mangaID)
	if errSelect != nil {
		fmt.Println(errSelect)
	}

	errDownload := manga.DownloadMangaChapter(&selectedChapter.ID, mangaName, &selectedChapter.ChapterNumber, DRECTORY_CACHE)
	if errDownload != nil {
		fmt.Print(errDownload.Error())
	}
}

func selectAppMode() (int, error) {
	prompt := promptui.Select{
		Label: "Choose an action",
		Items: []string{"Download Manga", "Read Manga"},
		Templates: &promptui.SelectTemplates{
			Active:   "▶ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✔ {{ . | green }}",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("selection cancelled or failed: %w", err)
	}

	return index, nil
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
