package main

import (
	"errors"
	"log"
	"math"
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

type rule struct {
	color1      string
	color2      string
	force       float64
	effectRange float64
}

type App struct {
	ctx      js.Value
	window   js.Value
	size     vec2i
	entities map[string][]*entity
	rules    []rule

	globalLoss float64
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

	rect := canvas.Call("getBoundingClientRect")
	if !rect.Truthy() {
		return nil, errors.New("Failed to initialize canvas bounding rect")
	}

	log.Printf("Canvas has size %dx%d", rect.Get("width").Int(), rect.Get("height").Int())

	canvas.Set("width", rect.Get("width").Int())
	canvas.Set("height", rect.Get("height").Int())

	return &App{
		ctx:        ctx,
		window:     jsWindow,
		size:       vec2i{x: rect.Get("width").Int(), y: rect.Get("height").Int()},
		entities:   map[string][]*entity{},
		rules:      []rule{},
		globalLoss: 0.98,
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

func (app *App) addRule(color1, color2 string, force float64) {
	app.rules = append(app.rules, rule{color1: color1, color2: color2, force: force, effectRange: float64(app.size.x / 5.0)})
}

func (app *App) initialize() {
	for i := 0; i < 20; i++ {
		app.addEntity("red")
	}
	for i := 0; i < 50; i++ {
		app.addEntity("blue")
	}

	app.addRule("red", "blue", -0.2)
	app.addRule("blue", "blue", -0.1)
	app.addRule("red", "red", 0.8)
}

func (app *App) update() {
	//* Apply rules
	for u := range app.rules {
		for i := range app.entities[app.rules[u].color1] {
			force := vec2f{}
			for j := range app.entities[app.rules[u].color2] {
				if app.rules[u].color1 != app.rules[u].color2 || i != j {
					dx := app.entities[app.rules[u].color1][i].pos.x - app.entities[app.rules[u].color2][j].pos.x
					dy := app.entities[app.rules[u].color1][i].pos.y - app.entities[app.rules[u].color2][j].pos.y
					dist := math.Sqrt(dx*dx + dy*dy)

					if dist < app.rules[u].effectRange {
						force.x -= app.rules[u].force / dist * dx
						force.y -= app.rules[u].force / dist * dy
					}
				}
			}

			app.entities[app.rules[u].color1][i].speed.x *= app.globalLoss
			app.entities[app.rules[u].color1][i].speed.y *= app.globalLoss
			app.entities[app.rules[u].color1][i].speed.x += force.x
			app.entities[app.rules[u].color1][i].speed.y += force.y

			app.entities[app.rules[u].color1][i].newPos.x += app.entities[app.rules[u].color1][i].speed.x
			app.entities[app.rules[u].color1][i].newPos.y += app.entities[app.rules[u].color1][i].speed.y

			// clamp
			if app.entities[app.rules[u].color1][i].newPos.x < 0 {
				app.entities[app.rules[u].color1][i].newPos.x = 0
				app.entities[app.rules[u].color1][i].speed.x *= -1
			}
			if app.entities[app.rules[u].color1][i].newPos.x > float64(app.size.x) {
				app.entities[app.rules[u].color1][i].newPos.x = float64(app.size.x)
				app.entities[app.rules[u].color1][i].speed.x *= -1
			}
			if app.entities[app.rules[u].color1][i].newPos.y < 0 {
				app.entities[app.rules[u].color1][i].newPos.y = 0
				app.entities[app.rules[u].color1][i].speed.y *= -1
			}
			if app.entities[app.rules[u].color1][i].newPos.y > float64(app.size.y) {
				app.entities[app.rules[u].color1][i].newPos.y = float64(app.size.y)
				app.entities[app.rules[u].color1][i].speed.y *= -1
			}
		}
	}

	//* Update positions
	for i := range app.entities {
		for j := range app.entities[i] {
			app.entities[i][j].pos = app.entities[i][j].newPos
		}
	}
}

func (app *App) render() {
	drawRect(app.ctx, "black", vec2i{}, app.size)

	for color := range app.entities {
		for i := range app.entities[color] {
			drawRect(app.ctx, color, app.entities[color][i].pos.tovec2i(), vec2i{4, 4})
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
