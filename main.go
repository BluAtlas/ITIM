package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nfnt/resize"
)

var PHONESIZEX int
var PHONESIZEY int
var APPSIZE int
var FROMTOP int
var FROMLEFT int
var FROMBOTTOM int
var BETWEENX int
var BETWEENY int
var FROMLEFTDOCK int
var MAXAPPS int
var DOCKCOUNT int
var APPCOUNT int
var FILETYPE string

func main() {
	// set log flags
	log.SetFlags(log.Lshortfile)
	// check for args
	if len(os.Args) < 2 {
		log.Fatalf("Wrong arguments, '1 [jpg or png]' for 12 pro max, '2 [jpg or png]' for 12 mini,\n\tor see the README for custom sizing information\n")
	}

	// set global variables based on args
	if len(os.Args) == 3 {
		switch os.Args[1] {
		case "1":
			PHONESIZEX = 1284
			PHONESIZEY = 2778
			APPSIZE = 192
			FROMTOP = 246
			FROMLEFT = 105
			BETWEENX = 102
			BETWEENY = 126
			FROMBOTTOM = 105
			FROMLEFTDOCK = 252
			MAXAPPS = 24
			DOCKCOUNT = 3
		case "2":
			PHONESIZEX = 1125
			PHONESIZEY = 2436
			APPSIZE = 180
			FROMTOP = 231
			FROMLEFT = 81
			BETWEENX = 81
			BETWEENY = 105
			FROMBOTTOM = 81
			FROMLEFTDOCK = 81
			MAXAPPS = 24
			DOCKCOUNT = 4
		}
		FILETYPE = os.Args[2]
	} else {
		if len(os.Args) != 13 {
			log.Fatalf("Wrong arguments, '1' for pro max, '2' for mini,\n\tor see the README for custom sizing information\n")
		}
		var args []int
		for i, arg := range os.Args {
			if i > 1 {
				val, err := strconv.Atoi(arg)
				check(err)
				args = append(args, val)
			}
		}
		FILETYPE = os.Args[1]
		PHONESIZEX = args[0]
		PHONESIZEY = args[1]
		APPSIZE = args[2]
		FROMTOP = args[3]
		FROMLEFT = args[4]
		BETWEENX = args[5]
		BETWEENY = args[6]
		FROMBOTTOM = args[7]
		FROMLEFTDOCK = args[8]
		MAXAPPS = args[9]
		DOCKCOUNT = args[10]
	}

	// ensure all dirs exist
	newpath := filepath.Join("input")
	err := os.MkdirAll(newpath, os.ModePerm)
	check(err)
	newpath = filepath.Join("output")
	err = os.MkdirAll(newpath, os.ModePerm)
	check(err)
	newpath = filepath.Join("input/icons")
	err = os.MkdirAll(newpath, os.ModePerm)
	check(err)
	newpath = filepath.Join("output/icons")
	err = os.MkdirAll(newpath, os.ModePerm)
	check(err)

	// create temp dir
	os.RemoveAll("temp")
	err = os.Mkdir("temp", 0755)
	check(err)
	defer os.RemoveAll("temp")

	// read background image
	reader, err := os.Open(fmt.Sprintf("input/background.%v", FILETYPE))
	check(err)
	m, _, err := image.Decode(reader)
	check(err)

	// resize the image, convert to a RGBA image
	m = ResizeImage(m)

	// create resized icons in temp dir
	ResizeIcons()

	// take pixels from background to create app squares, then overlay temp icons
	CreateIcons(m)

	// write the image to disk
	f, err := os.Create("output/background.png")
	check(err)
	err = png.Encode(f, m)
	check(err)
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

	//create new zero point
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
	// open directory for looping
	files, err := ioutil.ReadDir("input/icons")
	check(err)

	for _, file := range files { // for file in directory
		// open file into reader
		reader, err := os.Open(fmt.Sprintf("input/icons/%v", file.Name()))
		check(err)
		defer reader.Close()
		// decode reader as image
		m, _, err := image.Decode(reader)
		check(err)

		//resize image read
		m = resize.Resize(uint(APPSIZE), uint(APPSIZE), m, resize.Lanczos3)

		// write the resized image to disk
		f, err := os.Create(fmt.Sprintf("temp/%v", file.Name()))
		check(err)
		err = png.Encode(f, m)
		check(err)

		//increase app count
		APPCOUNT++
	}
}

func CreateIcons(m image.Image) {
	// open directory for looping
	files, err := ioutil.ReadDir("temp")
	check(err)

	// get starting point based off global variables
	cpoint := image.Point{X: m.Bounds().Min.X + FROMLEFT, Y: m.Bounds().Min.Y + FROMTOP}
	i := 0
	for _, file := range files { // for file in directory
		i++
		// open file into reader
		reader, err := os.Open(fmt.Sprintf("temp/%v", file.Name()))
		check(err)
		defer reader.Close()
		// decode reader as image
		n, _, err := image.Decode(reader)
		check(err)

		// create blank icon image, and draw onto it from background and icon overlay
		icon := image.NewRGBA(image.Rect(0, 0, APPSIZE, APPSIZE))
		draw.Draw(icon, icon.Bounds(), m, cpoint, draw.Over)
		draw.Draw(icon, icon.Bounds(), n, image.Point{X: icon.Bounds().Min.X, Y: icon.Bounds().Min.Y}, draw.Over)

		// write the icon to disk
		f, err := os.Create(fmt.Sprintf("output/icons/%v", file.Name()))
		check(err)
		err = png.Encode(f, icon)
		check(err)

		cpoint.X += BETWEENX + APPSIZE
		if i == APPCOUNT-DOCKCOUNT { // if the remaining icons are for the dock
			cpoint.Y = m.Bounds().Max.Y - FROMBOTTOM - APPSIZE
			cpoint.X = m.Bounds().Min.X + FROMLEFTDOCK
		} else if i > APPCOUNT-DOCKCOUNT {
		} else if i%4 == 0 { // move down to next row every four
			cpoint.Y += BETWEENY + APPSIZE
			cpoint.X = m.Bounds().Min.X + FROMLEFT
			if i%MAXAPPS == 0 { // if we're out of space on current page, move to next page
				cpoint.Y = m.Bounds().Min.Y + FROMTOP
			}
		}

	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
