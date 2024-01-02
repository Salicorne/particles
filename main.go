package main

import (
	"errors"
	"log"
	"syscall/js"
)

type App struct {
	ctx    js.Value
	window js.Value
	w      int
	h      int
}

func newApp() (*App, error) {
	jsWindow := js.Global().Get("window")
	if !jsWindow.Truthy() {
		return nil, errors.New("Failed to initialize DOM window")
	}

	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		return nil, errors.New("Failed to initialize DOM document")
	}

	canvas := jsDoc.Call("getElementById", "canvas")
	if !canvas.Truthy() {
		return nil, errors.New("Failed to initialize DOM canvas")
	}

	ctx := canvas.Call("getContext", "2d")
	if !ctx.Truthy() {
		return nil, errors.New("Failed to initialize canvas context")
	}

	return &App{
		ctx:    ctx,
		window: jsWindow,
		w:      canvas.Get("width").Int(),
		h:      canvas.Get("height").Int(),
	}, nil
}

func drawRect(ctx js.Value, color string, posx, posy, w, h int) {
	ctx.Set("fillStyle", color)
	ctx.Call("fillRect", posx, posy, w, h)
}

func (app *App) render() {
	drawRect(app.ctx, "black", 0, 0, app.w, app.h)
}

func main() {
	js.Global().Set("startApp", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Println("Hello world !")
		app, err := newApp()
		if err != nil {
			log.Println(err)
			return err
		}

		app.render()

		return nil
	}))

	<-make(chan bool) //wait
}
