package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"
	"time"

	colour "github.com/fatih/color"
	"github.com/nfnt/resize"
)

// used for storing pixel values
type RGBA struct {
	R, G, B, A uint32
}

func (rgb *RGBA) Subtract(rgb_ RGBA) RGBA {
	s := func(a uint32, b uint32) uint32 {
		return (uint32)((int)(a) - (int)(b))
	}
	return RGBA{
		s(rgb.R, rgb_.R), s(rgb.G, rgb_.G), s(rgb.B, rgb_.B), s(rgb.A, rgb_.A),
	}
}

var (
	characters []string
	inPath     string
	outPath    string
	scale      float64
	stretch    int
	prt        bool
	pretty     bool
	clr        bool

	// w/h values of image
	width  int
	height int

	pixels [][]RGBA

	successPrinter = colour.New(colour.FgHiGreen)
	errorPrinter   = colour.New(colour.FgHiRed)
)

// calculates a pixel's luminance from its RGB values
func luminance(pixel RGBA) uint8 {
	// consideration for transparent to nearly transparent pixels
	if pixel.A == 0 {
		return 0
	}
	return uint8(math.Floor((0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)) / 256))
}

// detects the name of a pixel's colour and returns the corresponding ANSI colour string
func rough(pixel RGBA) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", pixel.R, pixel.G, pixel.B)
	// correct pixel
	/* if pixel.R == pixel.G && pixel.G == pixel.B {
		if pixel.R < 8 {
			return fmt.Sprintf(s, 16)
		}
		if pixel.R > 248 {
			return fmt.Sprintf(s, 231)
		}
		return fmt.Sprintf(s, (int)(math.Round((float64)(((int)(pixel.R)-8)/247)*24)+232))
	}
	return fmt.Sprintf(s, 16+(int)(36*math.Round((float64)(pixel.R/255*5)))+(int)(6*math.Round((float64)(pixel.G/255*5)))+(int)(math.Round((float64)(pixel.B/255*5)))) */
}

func main() {
	// ensure in path exists
	if inPath == "" {
		errorPrinter.Println("no input path provided.")
		return
	} else {
		if _, err := os.Stat(inPath); err != nil {
			errorPrinter.Printf("input path '%s' does not exist.\n", inPath)
			return
		}
	}

	// read from and parse image data
	imPipe, err := os.Open(inPath)
	if err != nil {
		errorPrinter.Println("failed to read from input image.")
		return
	}
	defer imPipe.Close()
	if imConf, _, err := image.DecodeConfig(imPipe); err != nil {
		errorPrinter.Println("failed to decode image data - is the image a valid format?")
		return
	} else {
		width, height = imConf.Width, imConf.Height
	}
	imPipe.Seek(0, 0) // seek at start
	im, _, err := image.Decode(imPipe)
	if err != nil {
		errorPrinter.Println("failed whilst reading image data.")
		return
	}

	// resize image based on scale and stretch values
	height, width = int(math.Floor(float64(height)*scale)), int(math.Floor(float64(width)*scale*2*float64(stretch)))
	im = resize.Resize(uint(width), uint(height), im, resize.Lanczos3)

	// generate pixel array
	pixels = make([][]RGBA, height) // gen based on height
	for y := 0; y < height; y++ {
		pixels[y] = make([]RGBA, width)
		for x := 0; x < width; x++ {
			r, g, b, a := im.At(x, y).RGBA()
			pixels[y][x] = RGBA{
				R: r, G: g, B: b, A: a,
			}
		}
	}
	//successPrinter.Println("succesfully generated pixel array...") // hooray!

	// ooooh heavens...
	final := ""
	for _, layer := range pixels {
		for _, pixel := range layer {
			lum := luminance(pixel)
			// basically just get bright as a fraction of the max brightness
			bright := uint8(math.Floor((float64(lum) / 255) * 70))
			pos := len(characters) - int(bright)
			if pos == 70 {
				pos -= 1
			}
			if clr {
				final += rough(pixel) + characters[pos]
			} else {
				final += characters[pos]
			}
		}
		final += "\n"
	}

	// write final to outfile
	err = os.WriteFile(outPath, []byte(final), 0644)
	if err != nil {
		errorPrinter.Println("error writing to outfile!")
	} else {
		successPrinter.Printf("success! written to '%s'.\n", outPath)
	}
	if prt {
		if pretty {
			for _, s := range strings.Split(final, "\n") {
				fmt.Println(s)
				time.Sleep(time.Millisecond * 75)
			}
		} else {
			fmt.Print(final)
		}
	}
}

func init() {
	// create list of chars from standard
	characters = strings.Split("$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. ", "")

	// parse arguments
	flag.StringVar(&inPath, "in", "", "Specifies the input .png/jpg/jpeg file.")
	flag.StringVar(&outPath, "out", "out.txt", "Specifies the output .txt file.")
	flag.Float64Var(&scale, "scale", 0.5, "Specifies a scale factor.")
	flag.IntVar(&stretch, "stretch", 1, "Specifies a stretch factor.")
	flag.BoolVar(&prt, "print", false, "If passed, the result will be printed.")
	flag.BoolVar(&pretty, "pretty", false, "When '--print' is passed, output is printed layer-by-layer.")
	flag.BoolVar(&clr, "colour", false, "Colours the output.")

	// parse flags
	flag.Parse()

	// checks
	if scale <= 0 || stretch <= 0 {
		errorPrinter.Println("scale/stretch must be above 0.")
		os.Exit(1)
	}
}
