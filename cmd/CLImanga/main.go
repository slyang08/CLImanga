package main

import (
	"bufio"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"github.com/scinac/CLImanga/internal/manga"
	"github.com/scinac/CLImanga/internal/ui"
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
		appMode, err := ui.SelectAppMode()
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

		mangaID, err := ui.SelectManga(&mangaName)
		if err != nil {
			fmt.Println(err)
		}

		switch appMode {
		case APPMODE_DOWNLOAD:
			downloadMangaMode(&mangaID, &mangaName)
		case APPMODE_READ:
			readMangaMode(&mangaID, &mangaName)
		}
	}
}

func readMangaMode(mangaID *string, mangaName *string) {
	selectedChapter, chapterList, errSelect := ui.SelectChapterFromManga(mangaID)
	if errSelect != nil {
		fmt.Println(errSelect)
	}

	appInstance, appWindow := ui.InitGUIApp(*mangaName, chapterList, 900, 900)

	readChapter(mangaName, selectedChapter, appInstance, appWindow)
}

func downloadMangaMode(mangaID *string, mangaName *string) {
	manga.DownloadEntireManga(mangaID, mangaName)
}

func readChapter(mangaName *string, selectedChapter *manga.ChapterSelect, appInstance fyne.App, appWindow fyne.Window) {
	manga.DownloadMangaChapter(&selectedChapter.ID, mangaName, &selectedChapter.ChapterNumber, DRECTORY_CACHE)

	wd, _ := os.Getwd()
	var folder string = wd + "/resources/cache/" + *mangaName + "/chapter-" + selectedChapter.ChapterNumber

	ui.DisplayChapter(appWindow, folder)
	appInstance.Run()
}

func checkForDependencies() bool {
	return true
}
