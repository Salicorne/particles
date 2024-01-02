package main

import (
	"errors"
	"log"
	"math/rand"
	"syscall/js"
)

type vec2i struct {
	x, y int
}

type vec2f struct {
	x, y float64
}

func (v vec2f) tovec2i() vec2i { return vec2i{int(v.x), int(v.y)} }

type entity struct {
	pos    vec2f
	newPos vec2f
	speed  vec2f
}

type App struct {
	ctx      js.Value
	window   js.Value
	size     vec2i
	entities map[string][]*entity
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
		ctx:      ctx,
		window:   jsWindow,
		size:     vec2i{x: canvas.Get("width").Int(), y: canvas.Get("height").Int()},
		entities: map[string][]*entity{},
	}, nil
}

func drawRect(ctx js.Value, color string, pos, size vec2i) {
	ctx.Set("fillStyle", color)
	ctx.Call("fillRect", pos.x, pos.y, size.x, size.y)
}

func (app *App) addEntity(color string) {
	if _, ok := app.entities[color]; !ok {
		app.entities[color] = []*entity{}
	}

	initialPos := vec2f{float64(rand.Intn(app.size.x)), float64(rand.Intn(app.size.y))}
	app.entities[color] = append(app.entities[color], &entity{
		pos:    initialPos,
		newPos: initialPos,
		speed:  vec2f{},
	})
}

func (app *App) initialize() {
	for i := 0; i < 50; i++ {
		app.addEntity("red")
	}
	for i := 0; i < 50; i++ {
		app.addEntity("blue")
	}
}

func (app *App) update() {
	//todo update physics
}

func (app *App) render() {
	drawRect(app.ctx, "black", vec2i{}, app.size)

	for color := range app.entities {
		for i := range app.entities[color] {
			drawRect(app.ctx, color, app.entities[color][i].pos.tovec2i(), vec2i{2, 2})
		}
	}
}

func (app *App) run() {
	app.update()
	app.render()

	app.window.Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) any {
		app.run()
		return nil
	}))
}

func main() {
	js.Global().Set("startApp", js.FuncOf(func(this js.Value, args []js.Value) any {
		log.Println("Hello world !")
		app, err := newApp()
		if err != nil {
			log.Println(err)
			return err
		}

		app.initialize()
		app.run()

		return nil
	}))

	<-make(chan bool) //wait
}
