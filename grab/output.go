package grab

import (
	"encoding/json"
	"fmt"
)

// Output is the JSON output of this app
type Output struct {
	Title       string `json:"title"`
	Number      int    `json:"number"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func printJSON(x *XKCD) error {
	o := Output{
		Title:       x.Title,
		Number:      x.Num,
		Date:        x.Date().Format("2006-01-02"),
		Description: x.SafeTitle,
		Image:       x.Img,
	}
	b, err := json.MarshalIndent(&o, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func pp(x *XKCD) error {
	_, err := fmt.Printf(
		"Title: %s\nComic No: %d\nDate: %s\nDescription: %s\nImage: %s\n",
		x.Title,
		x.Num,
		x.Date().Format("02-01-2006"),
		x.SafeTitle,
		x.Img,
	)
	return err
}
