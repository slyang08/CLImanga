package manga

type MangaSelect struct {
	ID   string
	Name string
}

type ChapterSelect struct {
	ID            string
	ChapterNumber string
	Title         string
	Pages         float64
}
