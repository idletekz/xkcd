package grab

import (
	"fmt"
	"path"
	"time"
)

const (
	LatestComic = iota
)

// BuildURL returns json uri
func BuildURL(comicNumber int) string {
	if comicNumber == LatestComic {
		return "https://xkcd.com/info.0.json"
	}
	return fmt.Sprintf("https://xkcd.com/%d/info.0.json", comicNumber)
}

// XKCD struct
type XKCD struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

// Date returns a *time.Time based on the API strings (or nil if the response is malformed)
func (x *XKCD) Date() *time.Time {
	t, err := time.Parse(
		"2006-1-2",
		fmt.Sprintf("%s-%s-%s", x.Year, x.Month, x.Day),
	)
	if err != nil {
		return nil
	}
	return &t
}

func (x *XKCD) Filename() string {
	return path.Base(x.Img)
}
