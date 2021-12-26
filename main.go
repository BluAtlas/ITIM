package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

var PHONE string
var PHONESIZEX int
var PHONESIZEY int
var APPSIZE int
var FROMTOP int
var FROMLEFT int
var FROMBOTTOM int
var BETWEEN int

func main() {
	// check for args
	if len(os.Args) != 2 {
		log.Fatalf("Need an argument, '1' for pro max, '2' for mini\n")
	}

	// set global variables based on args
	switch os.Args[1] {
	case "1":
		PHONE = "12px"
		PHONESIZEX = 1080
		PHONESIZEY = 2340
		APPSIZE = 192
		FROMTOP = 246
		FROMLEFT = 105
		BETWEEN = 102
		FROMBOTTOM = 105
	case "2":
		PHONE = "12m"
		PHONESIZEX = 1284
		PHONESIZEY = 2778
		APPSIZE = 180
		FROMTOP = 387
		FROMLEFT = 81
		BETWEEN = 81
		FROMBOTTOM = 81
	}

	// read background image
	reader, err := os.Open("images/background.jpg")
	if err != nil {
		log.Fatal(err)
	}
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	// resize the image, convert to a RGBA image
	m = ResizeImage(m)

	ResizeIcons()

	// write the image to disk
	f, err := os.Create("images/outimage.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
	}
}

// resizes the background image to fit iphones aspect ratio, and changes format to RGBA
func ResizeImage(m image.Image) image.Image {
	//get image width, height, and ratio
	initwidth := m.Bounds().Max.X - m.Bounds().Min.X
	initheight := m.Bounds().Max.Y - m.Bounds().Min.Y
	width := initwidth
	height := initheight
	ratio := float64(height) / float64(width)
	DesiredRatio := float64(PHONESIZEY) / float64(PHONESIZEX)

	// if background image is not the right aspect ratio, change it to be
	if ratio < DesiredRatio-.0025 || DesiredRatio+.0025 < ratio {
		// if ratio is bigger than desired, cut height
		for ratio > DesiredRatio+.0025 {
			height--
			ratio = float64(height) / float64(width)
		}
		// if less than desired, cut width
		for ratio < DesiredRatio-.0025 {
			width--
			ratio = float64(height) / float64(width)
		}
		fmt.Printf("Original image cropped from %vx%v to %vx%v\n", m.Bounds().Max.X-m.Bounds().Min.X, m.Bounds().Max.Y-m.Bounds().Min.Y, width, height)
	} else {
		fmt.Printf("Image correct size.\n")
	}
	point := image.Point{}
	if height < m.Bounds().Max.Y-m.Bounds().Min.Y {
		// if height changed, adjust new zero point
		difference := (initheight - height) / 2
		point.X = m.Bounds().Min.X
		point.Y = difference - m.Bounds().Min.Y
	} else if width < m.Bounds().Max.X-m.Bounds().Min.X {
		// if width changed, adjust new zero point
		difference := (initwidth - width) / 2
		point.X = difference - m.Bounds().Min.X
		point.Y = m.Bounds().Min.Y
	}

	// now that we have new dimensions, create new image to crop original photo
	n := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(n, n.Bounds(), m, point, draw.Src)

	// resize/upscale to exact phone size
	m = resize.Resize(uint(PHONESIZEX), uint(PHONESIZEY), n, resize.Lanczos3)
	return m
}

func ResizeIcons() {

}
