package scrape

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/marcus-crane/khinsider/v2/pkg/util"

	"github.com/PuerkitoBio/goquery"
	"github.com/pterm/pterm"

	"github.com/marcus-crane/khinsider/v2/pkg/types"
)

const (
	LetterBase = "https://downloads.khinsider.com/game-soundtracks/browse/"
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

func GetResultsForLetter(letter string) (types.SearchResults, bool, error) {
	url := fmt.Sprintf("%s%s", LetterBase, letter)
	res, err := DownloadPage(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	if err != nil {
		return nil, false, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, false, err
	}
	results := make(types.SearchResults)
	doc.Find("table.albumList tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		s.Find("td a").Each(func(i int, t *goquery.Selection) {
			if i == 1 {
				title := strings.TrimSpace(t.Text())
				results[title] = "#"
				trackUrl, exists := t.Attr("href")
				if exists {
					results[title] = trackUrl
				}
			}
		})
	})
	more := false
	doc.Find(".pagination-next a").Each(func(i int, s *goquery.Selection) {
		_, more = s.Attr("href")
	})
	return results, more, nil
}

func RetrieveAlbum(slug string) (types.Album, error) {
	var album types.Album
	album.Slug = slug
	albumUrl := fmt.Sprintf("%s/game-soundtracks/album/%s", util.SiteBase, slug)

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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return album, err
	}
	metadata := doc.Find("#EchoTopic p[align='left'] b")
	if metadata.Length() == 5 {
		album.FlacAvailable = true
	}
	metadata.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			album.Name = s.Text()
		}
		if i == 1 {
			album.FileCount, err = strconv.Atoi(s.Text())
			if err != nil {
				album.FileCount = 0
			}
		}
		if i == 2 {
			album.MP3FileSize = s.Text()
		}
		if i == 3 && album.FlacAvailable {
			album.FlacFileSize = s.Text()
		}
	})
	flacLabel := ""
	if album.FlacAvailable {
		flacLabel = "[FLAC]"
	}
	pterm.Success.Printfln(
		"Found %s (%d tracks) %s %s",
		album.Name,
		album.FileCount,
		"[MP3]",
		flacLabel)
	pterm.Warning.Println("Searching for track locations up front. This may take seconds or minutes depending on album length!")
	_ = doc.Find("#EchoTopic table:not(#songList) tr div a").Each(func(i int, s *goquery.Selection) {
		coverUrl, exists := s.Attr("href")
		if exists {
			// TODO: Use proper escaping. Tried stdlib but it escaped everything
			coverUrl = strings.ReplaceAll(coverUrl, " ", "%20")
			album.Covers = append(album.Covers, coverUrl)
		}
	})
	songMeta := make(map[int]string)
	_ = doc.Find("#songlist_header th").Each(func(i int, s *goquery.Selection) {
		header := strings.TrimSpace(s.Text())
		if header == "CD" {
			songMeta[i] = "CD"
		}
		if header == "#" {
			songMeta[i] = "TrackNumber"
		}
		if header == "Song Name" {
			songMeta[i] = "SongName"
		}
		if header == "MP3" {
			songMeta[i] = "TrackLength"
			songMeta[i+1] = "MP3FileSize"
		}
		if header == "FLAC" {
			songMeta[i+1] = "FlacFileSize"
		}
	})
	doc.Find("#songlist tr:not(#songlist_header, #songlist_footer)").Each(func(i int, s *goquery.Selection) {
		var track types.Track
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			if songMeta[i] == "CD" {
				track.CDNumber = strings.TrimSpace(s.Text())
			}
			if songMeta[i] == "TrackNumber" {
				track.Number = strings.TrimSpace(s.Text())
			}
			if songMeta[i] == "SongName" {
				track.Name = strings.TrimSpace(s.Text())
			}
			if songMeta[i] == "TrackLength" {
				track.Duration = strings.TrimSpace(s.Text())
			}
			if songMeta[i] == "MP3FileSize" {
				track.MP3FileSize = strings.TrimSpace(s.Text())
				url, exists := s.Children().Attr("href")
				if exists {
					if !strings.Contains(url, "://") {
						url = fmt.Sprintf("%s%s", util.SiteBase, url)
					}
					res, err := DownloadPage(url)
					if err != nil {
						panic(err)
					}
					defer func(Body io.ReadCloser) {
						err := Body.Close()
						if err != nil {
							panic(err)
						}
					}(res.Body)
					if err != nil {
						panic(err)
					}
					page, err := goquery.NewDocumentFromReader(res.Body)
					if err != nil {
						panic(err)
					}
					src, exists := page.Find("audio").Attr("src")
					if exists {
						track.URL = src
					}
				}
			}
			if songMeta[i] == "FlacFileSize" {
				track.FlacFileSize = strings.TrimSpace(s.Text())
			}
		})
		album.Tracks = append(album.Tracks, track)
	})
	return album, nil
}
