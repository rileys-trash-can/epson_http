package main

import (
	"log"
	"os"

	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"sync"
)

var printerMu sync.Mutex
var printer *os.File

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: print <printer> <imgage>")
	}

	//img := os.Args[1]
	printerpath := os.Args[1]

	var err error
	printer, err = os.OpenFile(printerpath, os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("Failed to open printer at '%s': %s", printerpath, err)
	}

	//InitPrinter(printer)

	printimg(os.Args[2])
}

func printimg(img string) {
	f, err := os.OpenFile(img, os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("Failed to open image: %s", err)
		return
	}

	decimg, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decoe image: %s", err)
		return
	}

	const MAX_WIDTH = 560

	bds := decimg.Bounds()
	width, height := bds.Max.Y-bds.Min.Y, bds.Max.X-bds.Min.X
	log.Printf("Width: %d ; height %d", width, height)

	var nw, nh uint = MAX_WIDTH, uint(float64(width) * MAX_WIDTH / float64(height))
	log.Printf("NW: %d NH: %d", nw, nh)

	// resize image to right dimensions:
	resized := resize.Resize(nw, nh, decimg, resize.NearestNeighbor)

	dit := dither.NewDitherer(grays())
	dit.Mapper = dither.Bayer(8, 8, 1.0)

	dithered := dit.Dither(resized)

	Image(printer, dithered)
}

func grays() []color.Color {
	s := make([]color.Color, 0, 254)

	for i := 0; i < 254; i++ {
		s = append(s, &color.Gray{uint8(i)})
	}

	return s
}
