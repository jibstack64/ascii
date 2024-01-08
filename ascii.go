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
	closeClr   bool
    trueClr    bool

	// w/h values of image
	width  int
	height int

	pixels [][]RGBA

    ansi = map[string][2]interface{}{
        "Black":        { "\u001b[30m", RGBA{0, 0, 0, 0} },
        "Red":          { "\u001b[31m", RGBA{255, 0, 0, 0} },
        "Green":        { "\u001b[32m", RGBA{0, 255, 0, 0} },
        "Yellow":       { "\u001b[33m", RGBA{255, 255, 0, 0} },
        "Blue":         { "\u001b[34m", RGBA{0, 0, 255, 0} },
        "Magenta":      { "\u001b[35m", RGBA{255, 0, 255, 0} },
        "Cyan":         { "\u001b[36m", RGBA{0, 255, 255, 0} },
        "White":        { "\u001b[37m", RGBA{255, 255, 255, 0} },
        "BrightBlack":  { "\u001b[30;1m", RGBA{85, 85, 85, 0} },
        "BrightRed":    { "\u001b[31;1m", RGBA{255, 85, 85, 0} },
        "BrightGreen":  { "\u001b[32;1m", RGBA{85, 255, 85, 0} },
        "BrightYellow": { "\u001b[33;1m", RGBA{255, 255, 85, 0} },
        "BrightBlue":   { "\u001b[34;1m", RGBA{85, 85, 255, 0} },
        "BrightMagenta":{ "\u001b[35;1m", RGBA{255, 85, 255, 0} },
        "BrightCyan":   { "\u001b[36;1m", RGBA{85, 255, 255, 0} },
        "BrightWhite":  { "\u001b[37;1m", RGBA{255, 255, 255, 0} },
        "Reset":        { "\u001b[0m", RGBA{0, 0, 0, 0} },
    }

	successPrinter = colour.New(colour.FgHiGreen)
	errorPrinter   = colour.New(colour.FgHiRed)
)


func colorDistance(c1, c2 RGBA) float64 {
    return math.Sqrt(float64((c1.R-c2.R)*(c1.R-c2.R) + (c1.G-c2.G)*(c1.G-c2.G) + (c1.B-c2.B)*(c1.B-c2.B)))
}

func round(color RGBA) string {
    var minf = math.MaxFloat64
    var closest = ""
    for name, c := range ansi {
        distance := colorDistance(color, c[1].(RGBA))
        if distance < minf {
            minf = distance
            closest = name
        }
    }
    return ansi[closest][0].(string)
}

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
			if trueClr {
				final += rough(pixel) + characters[pos] + ansi["Reset"][0].(string)
            } else if closeClr {
                final += round(pixel) + characters[pos] + ansi["Reset"][0].(string)
            } else {
				final += characters[pos]
			}
		}
		final += "\n"
	}

	// write final to outfile
	if outPath != "" {
		err = os.WriteFile(outPath, []byte(final), 0644)
		if err != nil {
			errorPrinter.Println("error writing to outfile!")
		} else {
			successPrinter.Printf("success! written to '%s'.\n", outPath)
		}
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
	flag.StringVar(&outPath, "out", "", "Specifies the output .txt file.")
	flag.Float64Var(&scale, "scale", 0.5, "Specifies a scale factor.")
	flag.IntVar(&stretch, "stretch", 1, "Specifies a stretch factor.")
	flag.BoolVar(&prt, "print", false, "If passed, the result will be printed.")
	flag.BoolVar(&pretty, "pretty", false, "When '--print' is passed, output is printed layer-by-layer.")
	flag.BoolVar(&closeClr, "close-colour", false, "Colours the output by rounding RGB values to the closest available ANSI codes.")
    flag.BoolVar(&trueClr, "true-colour", false, "Colours the output using exact RGB-ANSI codes. Not supported on most consoles.") 

	// parse flags
	flag.Parse()

	// checks
	if scale <= 0 || stretch <= 0 {
		errorPrinter.Println("scale/stretch must be above 0.")
		os.Exit(1)
	}
}
