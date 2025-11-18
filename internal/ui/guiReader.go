package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/scinac/CLImanga/internal/manga"
	"github.com/scinac/CLImanga/internal/utils"
)

func InitGUIApp(title string, width float32, height float32) (fyne.App, fyne.Window) {
	appGui := app.New()
	window := appGui.NewWindow(title)
	window.Resize(fyne.NewSize(width, height))
	return appGui, window
}

func DisplayChapter(w fyne.Window, mode rune, mangaName string, chapterInfo *manga.ChapterSelect, chapterList *[]manga.ChapterSelect) {
	wd, _ := os.Getwd()
	var folder string = wd
	folder += "/resources/cache/" + mangaName + "/chapter-" + chapterInfo.ChapterNumber

	var saveFile string = wd
	saveFile += "/resources/history/history.json"

	imgContainer := container.NewVBox()
	scroll := container.NewVScroll(imgContainer)
	scroll.SetMinSize(fyne.NewSize(600, 800))

	nextBtn := widget.NewButton("Next Chapter", func() {
		deleteAllCurrentChaptersInCache(folder) // return value error maybe check with log system
		DisplayChapter(w, mode, mangaName, &(*chapterList)[chapterInfo.Index+1], chapterList)
	})

	prevBtn := widget.NewButton("Previous Chapter", func() {
		deleteAllCurrentChaptersInCache(folder)
		DisplayChapter(w, mode, mangaName, &(*chapterList)[chapterInfo.Index-
			2], chapterList)
	})

	content := container.NewBorder(prevBtn, nextBtn, nil, nil, scroll)
	w.SetContent(content)
	w.Show()

	if mode == 'r' { // if in read mode (download mode wont need to download again)
		ch := make(chan string)

		go manga.DownloadMangaChapter(&(chapterInfo.ID), &mangaName, &chapterInfo.ChapterNumber, "cache", ch) // TODO add multithreading go

		newentry := manga.HistorySave{
			ChapterID:    chapterInfo.ID,
			MangaName:    mangaName,
			ChapterNuber: chapterInfo.ChapterNumber,
		}

		err := utils.SaveEntryToFile(saveFile, newentry)
		if err != nil {
			fmt.Println("Error saving history:", err)
		}

		go func() {
			for file := range ch {
				fyne.DoAndWait(func() {
					img := canvas.NewImageFromFile(file)
					img.FillMode = canvas.ImageFillContain
					img.SetMinSize(fyne.NewSize(600, 800))
					imgContainer.Add(img)
					imgContainer.Refresh()
				})
			}
		}()
	} else if mode == 'd' { // download mode
		files, err := filepath.Glob(folder + "/*.jpg")
		if err != nil {
			log.Fatal(err)
		}
		if len(files) == 0 {
			log.Fatal("No images found in folder")
		}
	}
}

func deleteAllCurrentChaptersInCache(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
