package khinsider

import (
	"github.com/marcus-crane/khinsider/v3/pkg/download"
	"github.com/marcus-crane/khinsider/v3/pkg/indexer"
	"github.com/marcus-crane/khinsider/v3/pkg/scrape"
	"github.com/marcus-crane/khinsider/v3/pkg/search"
	"github.com/pterm/pterm"
)

func BeforeSearch() error {
	pterm.Warning.Println("Fetching the latest search index.")
	err := indexer.DownloadIndex()
	if err != nil {
		return err
	}
	return nil
}

func SearchAction(flacMode bool) error {
	//TODO: Add a possibilty to select the mode from the tui
	index, err := indexer.LoadLocalIndex()
	if err != nil {
		panic(err)
	}
	selectedSlugs, err := search.FilterAlbumList(index)
	if err != nil {
		panic(err)
	}
	err = DownloadAction(selectedSlugs, flacMode)
	if err != nil {
		pterm.Error.Println("Failed to download album")
	}
	return nil
}

func DownloadAction(slugs []string, flacMode bool) error {
	for _, slug := range slugs {
		album, err := scrape.RetrieveAlbum(slug)
		if err != nil {
			return err
		}
		download.GetAlbum(&album, flacMode)
	}
	return nil
}
