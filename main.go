package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"sync"

	"github.com/gorilla/mux"

	_ "embed"
	"net/http"
)

//go:embed help.txt
var embedHelp []byte

var printerMu sync.Mutex
var printer *os.File

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

	r := mux.NewRouter()

	/*imageFile, err := os.OpenFile(img, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("Failed to open image at %s: %s", img, err)
	}*/

	r.HandleFunc("/", handleHelp).Methods("GET")
	r.HandleFunc("/print", handlePrint).Methods("PUT")
	r.HandleFunc("/cut", handleCut).Methods("GET")
	r.HandleFunc("/text", handleText).Methods("PUT")

	err = http.ListenAndServe("[::]:8080", r)
	log.Fatalf("Failed to listenandserve: %s", err)
}

func handleHelp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	w.Write(embedHelp)
}

func handleCut(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] is cutting", r.RemoteAddr)

	printerMu.Lock()
	defer printerMu.Unlock()

	CutPaper(printer)

	w.Write([]byte("Cut Sucessfull!\n"))
}

func handleText(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] is printing some text", r.RemoteAddr)

	printerMu.Lock()
	defer printerMu.Unlock()

	body, _ := io.ReadAll(r.Body)
	printer.Write(body)
	LineBreak(printer)
	LineBreak(printer)

	w.Write([]byte("Print Sucessfull!\n"))
}

func handlePrint(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] is printing an image", r.RemoteAddr)

	decimg, _, err := image.Decode(r.Body)
	if err != nil {
		log.Printf("Failed to decoe image: %s", err)
		w.Write([]byte("NONONONO!!! " + err.Error()))
		return
	}

	const MAX_WIDTH = 560
	// resize image to right dimensions:
	resized := resize.Resize(MAX_WIDTH, 0, decimg, resize.NearestNeighbor)

	dit := dither.NewDitherer(grays())
	dit.Mapper = dither.Bayer(8, 8, 1.0)

	dithered := dit.Dither(resized)

	printerMu.Lock()
	defer printerMu.Unlock()

	LineBreak(printer)
	Image(printer, dithered)
	LineBreak(printer)
	LineBreak(printer)
	LineBreak(printer)

	time.Sleep(time.Second * 5)

	//CutPaper(printer)

	w.Write([]byte("Printed Sucessfully!\n"))
}

func grays() []color.Color {
	s := make([]color.Color, 0, 254)

	for i := 0; i < 254; i++ {
		s = append(s, &color.Gray{uint8(i)})
	}

	return s
}
