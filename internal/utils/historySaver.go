package utils

import (
	"encoding/json"
	"os"

	"github.com/scinac/CLImanga/internal/manga"
)

func SaveEntryToFile(filePath string, entry manga.HistorySave) error {
	var entries []manga.HistorySave

	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		json.Unmarshal(data, &entries)
	}

	entries = append(entries, entry)

	jsonData, err := json.MarshalIndent(entries, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0o644)
}
