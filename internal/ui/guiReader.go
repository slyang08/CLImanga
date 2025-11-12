package ui

import (
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func InitGUIApp(title string, width, height float32) (fyne.App, fyne.Window) {
	appGui := app.New()
	window := appGui.NewWindow(title)
	window.Resize(fyne.NewSize(width, height))
	return appGui, window
}

func DisplayChapter(w fyne.Window, folder string) {
	log.Println(folder)
	files, err := filepath.Glob(folder + "/*.jpg")
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		log.Fatal("No images found in folder")
	}

	var images []fyne.CanvasObject
	for _, f := range files {
		img := canvas.NewImageFromFile(f)
		img.FillMode = canvas.ImageFillStretch
		img.SetMinSize(fyne.NewSize(600, 800))
		images = append(images, img)
	}

	scroll := container.NewVScroll(container.NewVBox(images...))
	scroll.SetMinSize(fyne.NewSize(600, 800))

	w.SetContent(scroll)
	w.Show()
}
