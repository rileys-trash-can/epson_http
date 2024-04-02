package main

import (
	"fmt"
	"log"
	"os"

	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

var printer *os.File

const width = len("                                               .")

func Center(str string) string {
	return strings.Repeat(" ", width/2-len(str)/2) + str
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: print <printer>")
	}

	//img := os.Args[1]
	printerpath := os.Args[1]

	var err error
	printer, err = os.OpenFile(printerpath, os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("Failed to open printer at '%s': %s", printerpath, err)
	}

	InitPrinter(printer)

	printImage("header.png")
	fmt.Fprintln(printer, Center("Schaffenburg e.V."))
	fmt.Fprintln(printer, Center"Dorfstr. 1"))
	fmt.Fprintln(printer, Center("63741 Damm"))
	fmt.Fprintln(printer, Center("ZWEIGSTELLE 46523C"))
	fmt.Fprintln(printer, "ö")
	fmt.Fprintln(printer, "   Es bedienten sie:")
	fmt.Fprintln(printer, "                     normale menschen (endlich)")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "   Sie kauften:")
	fmt.Fprintln(printer, "        MAAATEEE . . . . . . . . . . 0.99 EUR")
	fmt.Fprintln(printer, "          1 x 0.99 EUR")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "        AVM FRITZ! Kola  . . . . . . 2.98 EUR")
	fmt.Fprintln(printer, "          2 x 1.49 EUR")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "        Belasto(TM) RIEGEL . . . . . 0.98 EUR")
	fmt.Fprintln(printer, "          2 x 0.49 EUR")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "        Radithor . . . . . . . . . . 4.90 EUR")
	fmt.Fprintln(printer, "          1 x 4.90 EUR")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "        Ferrero DUSPOL Riegel. . . . viel EUR")
	fmt.Fprintln(printer, "          1 x viele EUR")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "   Total: undefined EUR")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, " Zahl: BAR")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "   Gegeben: {object Object} EUR")
	fmt.Fprintln(printer, "   Rueckgeld: NaN EUR")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "Transaktion 6405GFNGK34t9uedflvß0WnagDI43H9S.MFLKFDJH	490gDFJaertjsfdiohgaT$HSDILVNWlVGRIHWNGkgNHIHIFAnhergji4$$")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "")
	fmt.Fprintln(printer, "")

	CutPaper(printer)
}

func printImage(path string) {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("Failed to open image: %s", err)
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
