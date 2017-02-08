package render

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"sync"
)

// Renderer is an interface to be used with Engine
type Renderer interface {
	Render(x, y int) color.Color
}

// Image is compatible with all standard library formats
type Image interface {
	Set(x, y int, c color.Color)
	Bounds() image.Rectangle
}

// Engine is the rendering engine that has to be configured
type Engine struct {
	width, height int
	workers       int
	image         Image
	renderer      Renderer
}

// NewEngine return a new Engine value pointer with width and height
func NewEngine(r Renderer, img Image) *Engine {
	e := new(Engine)
	e.workers = runtime.NumCPU()
	e.width = img.Bounds().Max.X - img.Bounds().Min.X
	e.height = img.Bounds().Max.Y - img.Bounds().Min.Y
	e.image = img
	e.renderer = r
	return e
}

// SetWorkers changes the number of goroutines that will be run.
// default: runtime.NumCPU() - This usually returns the number of logical processors
func (e *Engine) SetWorkers(num int) {
	e.workers = num
}

// worker is just a wrapper around the Renderer that gets fed lines from the channel
// and renders each line until it's done (channel closed)
func (e *Engine) worker(line chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for y := range line {
		fmt.Printf("Rendering line %.3d of %d\r", y+1, e.height)
		for x := 0; x < e.width; x++ {
			e.image.Set(x, y, e.renderer.Render(x, y))
		}
	}
}

// Run executes the Render function of the Renderer in a multi threaded environment
// and sends it's result to the Set function of the image to save the color.
func (e *Engine) Run() {
	// Create a line channel to feed the workers
	fmt.Printf("Running with %d workers...\n", e.workers)
	line := make(chan int)
	// and a WaitGroup to sync them
	var wg sync.WaitGroup
	// We run the workers sharing the line channel and the WaitGroup
	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go e.worker(line, &wg)
	}
	// Iterate over the height to feed the channel
	for y := 0; y < e.height; y++ {
		line <- y
	}
	// Close the channel (no more work)
	close(line)
	// Wait for all workers to finish
	wg.Wait()
	fmt.Printf("Rendering complete! (%d lines)\n", e.height)
}
