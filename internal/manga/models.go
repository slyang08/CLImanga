package manga

type MangaSelect struct {
	ID   string
	Name string
}

type ChapterSelect struct {
	Index         int
	ID            string
	ChapterNumber string
	Title         string
	Pages         float64
}
type HistorySave struct {
	ChapterID    string
	MangaName    string
	ChapterNuber string
}
