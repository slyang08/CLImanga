package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"github.com/scinac/CLImanga/internal/log"
	"github.com/scinac/CLImanga/internal/manga"
	"github.com/scinac/CLImanga/internal/ui"
)

const (
	APPMODE_DOWNLOAD   = 1
	APPMODE_READ       = 0
	DIRECTORY_DOWNLOAD = "downloads"
	DRECTORY_CACHE     = "cache"
)

func main() {
	log.Init()
	log.Info.Println("Application Started.")
	catchProgramExit()
	fmt.Println("Welcome to CLImanga!")
	fmt.Println("Checking Dependecies... ")
	if !checkForDependencies() {
		log.Error.Println("Missing dependencies")
		fmt.Println("Missing Dependencies")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	if len(os.Args) < 2 {
		appMode, err := ui.SelectAppMode()
		if err != nil {
			log.Error.Printf("App mode selection error: %v", err)
			fmt.Print(err.Error())
			return
		}

		fmt.Print("Search Manga: ")
		mangaName, err := reader.ReadString('\n')
		if err != nil {
			log.Error.Printf("Reading manga name error: %v", err)
			fmt.Print(err.Error())
			return
		}

		mangaID, err := ui.SelectManga(&mangaName)
		if err != nil {
			log.Error.Printf("Manga selection error: %v", err)
			fmt.Println(err)
		}

		switch appMode {
		case APPMODE_DOWNLOAD:
			log.Info.Printf("Starting Download mode for manga: %s", mangaName)
			log.WrapFunction(func() {
				downloadMangaMode(&mangaID, &mangaName)
			})()
		case APPMODE_READ:
			log.Info.Printf("Starting Read mode for manga: %s", mangaName)
			mangaName = strings.TrimSpace(mangaName)
			log.WrapFunction(func() {
				readMangaMode(&mangaID, &mangaName)
			})()
		default:
			log.Error.Printf("Unknown app mode selected: %d", appMode)
		}
	}
}

func readMangaMode(mangaID *string, mangaName *string) {
	log.LogFunctionName()
	selectedChapter, chapterList, errSelect := ui.SelectChapterFromManga(mangaID)
	if errSelect != nil {
		log.Error.Printf("Chapter selection error: %v", errSelect)
		fmt.Println(errSelect)
	} else {
		log.Info.Printf("Read mode: Manga %s, Chapter %s", *mangaName, selectedChapter.ChapterNumber)
	}

	appInstance, appWindow := ui.InitGUIApp(*mangaName, 900, 900)

	readChapter(mangaName, selectedChapter, chapterList, appInstance, appWindow)
}

func downloadMangaMode(mangaID *string, mangaName *string) {
	log.LogFunctionName()
	log.Info.Printf("Download mode: Manga %s", *mangaName)
	err := manga.DownloadEntireManga(mangaID, mangaName)
	if err != nil {
		log.Error.Printf("Download error for manga %s: %v", *mangaName, err)
	}
}

func readChapter(mangaName *string, selectedChapter *manga.ChapterSelect, chapterList *[]manga.ChapterSelect, appInstance fyne.App, appWindow fyne.Window) {
	log.LogFunctionName()
	ui.DisplayChapter(appWindow, 'r', *mangaName, selectedChapter, chapterList)

	appInstance.Run()
}

func checkForDependencies() bool {
	return true
}

func catchProgramExit() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM) // Catch Ctrl+C and SIGTERM

	go func() {
		<-sigs
		log.Info.Println("Gracefully shutting down...")
		fmt.Println("\nGracefully shutting down...")
		cleanupProgram()
		os.Exit(0)
	}()
}

func cleanupProgram() {
	log.Info.Println("Cleaning up cache files on exit.")
	err := deleteAllCacheFiles()
	if err != nil {
		log.Error.Printf("Error cleaning cache files: %v", err)
	}
}

func deleteAllCacheFiles() error {
	wd, _ := os.Getwd()
	var path string = wd + "/resources/cache"

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.RemoveAll(path + "/" + file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}
