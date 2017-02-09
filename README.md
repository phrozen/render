# render
Go Multi-Threaded Rendering Engine for Computer Graphic Applications

## About
*(early version, work in progress, subject to change*)

In computer graphic applications, from raytracers or image compositing to fractal generators and anything in between, you find yourself with a lot of similar code, applications like those are commonly reduced to a double iteration (width and height) of an image, where you make calculations to obtain a certain **color** in a given position **(x, y)** and after calculating every pixel you save the result as a *jpg* or *png*.

It becomes more obvious and repetitive the more you develop this kind of projects that most of this algorithms are meant to be parallelized, as the state of the application does not change between iterations and refactoring your application each time to include code for multi-threading purposes becomes tedious and repetitive.

**render** tries to solve that in a simple way by providing an already parallelized rendering engine.


## Install
```
go get -u github.com/phrozen/render
```

## Usage
**render** works by making assumptions about the most common scenario of graphic applications, and exposes 2 simple interfaces.

```
type Renderer interface {
  Render(x, y int) color.Color
}
```

Following Go's idiomatic writting, a **Renderer** is any type that implements the **Render** function, which, (more often than not) usually is the code between the dual cycles that iterate over the width and height of the final image, calculating one pixel at a time.
 ```
 for y := 0; y < height; y++ { 
   for x := 0; x < width; x++ {

     // CALCULATE THE PIXEL (COLOR) VALUE FOR X,Y
     // STORE THE RESULTING COLOR IN IMAGE[X,Y]

   } 
}
 ```
That code can easily be extracted to a function ```Render(x, y int) color.Color``` and all the relevant data needed for the calculation can be stored inside a ```struct``` type so that it satisfies the **Renderer** interface. It is important to note that ```color.Color``` was used in order to make it easy and universal by utilizing teh standard library as much as possible.

```
type Image interface {
	Set(x, y int, c color.Color)
	Bounds() image.Rectangle
}
```

The second is an **Image** interface that can be satisfied by any of the image types inside the ```image``` pkg in Go standard library, there is nothing else to do but to initialize an image and save it afterwards.

```func NewEngine(r Renderer, img Image) *Engine```

To initialize a new rendering engine just call the function ```Engine``` and provide any **Renderer** and any **Image** and then simply call the function ```Run()``` on your engine. The ammout of workers (goroutines) will be ```runtime.NumCPU()``` by default, and it is usually the number of logical processors. The engine will run the double cycles of the image by rendering each line inside a different worker with a queue and saving the result back on the image. *(Note: Image can be a field of the Renderer provided both satisfy their respective interfaces)*

## Example

```
type Gradient struct {
	width, height int
	image         *image.RGBA
}

func NewGradient(width, height int) *Gradient {
	return &Gradient{width, height, image.NewRGBA(image.Rect(0, 0, width, height))}
}

func (gr *Gradient) Render(x, y int) color.Color {
	r := uint8((x * 255) / gr.width)
	g := uint8((y * 255) / gr.height)
	b := uint8(255 - r)
	return color.RGBA{r, g, b, 255}
}
```

Let's say you have a **Gradient** type like the one above, that implements the **Render** function. The **image.RGBA** type from the standard library, already satisfies the **Image** interface from **render**. The next step would be to create a new **Engine** by providing a **Renderer** (which our **Gradient** type satisfies) and an **Image** (provided by the **image** field in our **Gradient** struct type). By calling the function ```Run()``` the library automatically does the following:

1. Create N workers, a communication channel and synchronization mechanism. *(Workers are reused, only ```NumCPU()``` goroutines are created.)*
2. Iterate over all the image pixels calling ```img.Set(x, y, Render(x, y))``` asynchronously.
3. Close channels, and wait for all workers to finish rendering.

```
func main() {
  gr := NewGradient(1024, 1024)
  e := render.NewEngine(gr, gr.image)
  e.Run()
  // Utility function
  render.SavePNG("gradient.png", gr.image)
}
```

This will yield a nice image (like the one below) but most importantly you will see something like this in your console. You will also see the progress (by line) while it's being rendered, but more importantly it saves a lot of repetitive coding and execution time.

```
Running with 4 workers...
Rendering complete! (1024 lines)
Saved gradient.png
```