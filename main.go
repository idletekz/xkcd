// grabxkcd get xkcd comic
package main

import (
	"os"

	"github.com/idletekz/xkcd/grab"
)

func main() {
	os.Exit(grab.CLI(os.Args[1:]))
}
