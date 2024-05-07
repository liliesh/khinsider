package scrape

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/marcus-crane/khinsider/v3/cmd/khinsider/env"
	"github.com/marcus-crane/khinsider/v3/pkg/types"
	"github.com/marcus-crane/khinsider/v3/pkg/util"
)

const (
	IntIndexAlbumBase = "albums"
	ExtIndexAlbumBase = "https://khindex.utf9k.net/albums"
)

func DownloadPage(url string) (*http.Response, error) {
	res, err := util.MakeRequest(url, http.Header{})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received a non-200 status code: %d", res.StatusCode)
	}
	return res, err
}

func RetrieveAlbum(slug string) (types.Album, error) {
	var album types.Album
	var albumUrl string

	if env.GetAppFlags().LocalIndex {
		albumUrl = fmt.Sprintf("%s/%s.json", IntIndexAlbumBase, slug)
		filePath := fmt.Sprintf("%s/%s", env.GetCachePath(), albumUrl)

		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		if err := util.LoadJSON(file, &album); err != nil {
			return album, nil
		}

	} else {
		albumUrl = fmt.Sprintf("%s/%s.json", ExtIndexAlbumBase, slug)
		res, err := DownloadPage(albumUrl)
		if err != nil {
			return album, err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				panic(err)
			}
		}(res.Body)

		err = util.LoadJSON(res.Body, &album)
		if err != nil {
			return album, err
		}
	}

	return album, nil
}
