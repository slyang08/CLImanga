package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/scinac/CLImanga/internal/manga"
)

func SelectAppMode() (int, error) {
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

func SelectChapterFromManga(mangaID *string) (*manga.ChapterSelect, *[]manga.ChapterSelect, error) {
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
		return nil, nil, fmt.Errorf("chapter selection cancelled or failed: %w", err)
	}

	selectedChapter := chapterList[i]
	return &selectedChapter, &chapterList, nil
}

func SelectManga(mangaName *string) (string, error) {
	fmt.Println("Searching mangas with Name:", *mangaName)

	mangasFound, err := manga.FetchMangasByNameSearch(mangaName)
	if err != nil {
		return "", err
	}

	if len(mangasFound) == 0 {
		return "", fmt.Errorf("sorry no mangas found with that name")
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
