package main

import (
	"image"
	"image/color"

	"github.com/phrozen/render"
)

// The Gradient type stores just enough data to make it work, but this could
// be a whole scene definition for a Raytracer or source images for compositing.
type Gradient struct {
	width, height int
	// Any formated type from the 'image' package in the standard library
	// already satisfies the render.Image interface needed.
	image *image.RGBA
}

// We create a new Gradient with the given width and height
func NewGradient(width, height int) *Gradient {
	return &Gradient{width, height, image.NewRGBA(image.Rect(0, 0, width, height))}
}

// By implementing the Render(x, y) color.Color function, the Gradient type satisfies
// the render.Renderer interface needed for the render engine.
func (gr *Gradient) Render(x, y int) color.Color {
	r := uint8((x * 255) / gr.width)
	g := uint8((y * 255) / gr.height)
	b := uint8(255 - r)
	return color.RGBA{r, g, b, 255}
}

func main() {
	gr := NewGradient(1024, 1024)
	e := render.NewEngine(gr, gr.image)
	e.Run()
	// Common utility function included
	render.SavePNG("gradient.png", gr.image)
}
