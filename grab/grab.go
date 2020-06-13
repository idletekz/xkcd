package grab

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type appEnv struct {
	hc          *http.Client
	comicNumber int
	saveImage   bool
	outputJSON  bool
}

// CLI get commic and return error code status
func CLI(args []string) int {
	var app appEnv
	if err := app.parse(args); err != nil {
		return 2
	}
	if err := app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v", err)
		return 1
	}
	return 0
}

func (app *appEnv) parse(args []string) error {
	// Shallow copy of http.Client
	app.hc = &*http.DefaultClient
	fl := flag.NewFlagSet("xkcd-grab", flag.ContinueOnError)
	fl.IntVar(&app.comicNumber, "n", LatestComic, "Comic number to fetch")
	fl.DurationVar(&app.hc.Timeout, "t", 30*time.Second, "Client timeout")
	fl.BoolVar(&app.saveImage, "s", false, "Save image to current directory")
	outputType := fl.String("o", "text", "Print output in format: text/json")
	if err := fl.Parse(args); err != nil {
		return err
	}
	if *outputType != "text" && *outputType != "json" {
		fmt.Fprintf(os.Stderr, "got bad output type: %q\n", *outputType)
		fl.Usage()
		return flag.ErrHelp
	}
	app.outputJSON = *outputType == "json"
	return nil
}

func (app *appEnv) run() error {
	u := BuildURL(app.comicNumber)

	var result XKCD
	if err := app.fetchJSON(u, &result); err != nil {
		return err
	}
	if result.Date() == nil {
		return fmt.Errorf("could not parse date of comic: %q/%q/%q",
			result.Year, result.Month, result.Day)
	}
	if app.saveImage {
		if err := app.fetchAndSave(result.Img, result.Filename()); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Saved: %q\n", result.Filename())
	}
	if app.outputJSON {
		return printJSON(&result)
	}
	return pp(&result)
}

func (app *appEnv) fetch(url string) (*http.Response, error) {
	resp, err := app.hc.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch failed: %v", resp.Status)
	}
	return resp, nil
}

func (app *appEnv) fetchJSON(url string, data interface{}) error {
	resp, err := app.fetch(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func (app *appEnv) fetchAndSave(url, destPath string) error {
	resp, err := app.fetch(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, resp.Body)
	return err

}
