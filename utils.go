package render

import (
	"image"
	"image/png"
	"os"
)

// SavePNG utility function, uses standard library encoder to easily
// save any image.Image interface in .png format.
func SavePNG(filename string, img image.Image) error {
	output, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer output.Close()
	if err = png.Encode(output, img); err != nil {
		return err
	}
	return nil
}
