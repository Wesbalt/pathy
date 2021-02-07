package main

import (
	"os"
	"errors"
	"image"
	"math"
	"golang.org/x/image/draw"
	"image/color"
	"image/jpeg"
	"github.com/llgcode/draw2d/draw2dimg"
)

/*
 * Creates an image based on a map.
 * White cells are open, black cells are blocked.
 */
func MakeMapImage(scale int) *image.RGBA {
	h := len(grid)
	w := len(grid[0])
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := color.White
			if grid[y][x] {
				c = color.Black
			}
			img.Set(x, y, c)
		}
	}
	newWidth  := img.Bounds().Max.X * scale
	newHeight := img.Bounds().Max.Y * scale
	newSize   := image.Rect(0, 0, newWidth, newHeight)
	scaledImg := image.NewRGBA(newSize)
	draw.NearestNeighbor.Scale(scaledImg, newSize, img, img.Bounds(), draw.Over, nil)
	return scaledImg
}

/*
 * The supplied filename should include a "jpg" suffix as the image
 * is encoded to a JPEG file.
 */
func SaveImage(img *image.RGBA, fname string) error {
	out, err := os.Create(fname)
	defer out.Close()
	if err != nil {
		return errors.New("Could not create the file")
	}
	options := jpeg.Options{}
	options.Quality = 100 // Highest
    jpeg.Encode(out, img, &options)
	return nil
}

func DrawPath(img *image.RGBA, path []Node, scale int) *image.RGBA {
	if len(path) == 0 {
		return img
	}
	scaleF    := float64(scale)
	lineWidth := math.Ceil(scaleF/4)

	gc := draw2dimg.NewGraphicContext(img)
	defer gc.Close()

	gc.SetLineWidth(lineWidth)
	var prevX, prevY float64
	for i, n := range(path) {
		x := float64(n.X)
		y := float64(n.Y)
		if i > 0 {
			// Line between path nodes
			gc.SetStrokeColor(color.RGBA{255,0,0,255})
			gc.BeginPath()
			gc.MoveTo(scaleF * prevX, scaleF * prevY)
			gc.LineTo(scaleF * x,     scaleF * y)
			gc.Stroke()
		}
		// Diamond at each path node
		size := 0.2
		gc.SetStrokeColor(color.RGBA{0,0,255,255})
		gc.BeginPath()
		gc.MoveTo(scaleF * x,          scaleF * (y - size))
		gc.LineTo(scaleF * (x + size), scaleF * y)
		gc.LineTo(scaleF * x,          scaleF * (y + size))
		gc.LineTo(scaleF * (x - size), scaleF * y)
		gc.LineTo(scaleF * x,          scaleF * (y - size))
		gc.Stroke()
		prevX = x
		prevY = y
	}
	return img
}
