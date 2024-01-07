package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"sync"
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

type Rule struct {
	Color1      string  `json:"c1,omitempty"`
	Color2      string  `json:"c2,omitempty"`
	Force       float64 `json:"f,omitempty"`
	EffectRange float64 `json:"r,omitempty"`
}

type App struct {
	ctx      js.Value
	window   js.Value
	size     vec2i
	entities map[string][]*entity
	rules    []Rule

	globalLoss float64
	mutex      sync.Mutex
}

type Settings struct {
	Colors map[string]uint `json:"c,omitempty"`
	Rules  []Rule          `json:"r,omitempty"`
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
		rules:      []Rule{},
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

func (app *App) addRule(color1, color2 string, force float64, effectRange float64) {
	app.rules = append(app.rules, Rule{Color1: color1, Color2: color2, Force: force, EffectRange: effectRange})
}

func (app *App) initialize(settings *Settings) {
	app.mutex.Lock()

	// Reset settings
	app.entities = map[string][]*entity{}
	app.rules = []Rule{}

	for c := range settings.Colors {
		for i := uint(0); i < settings.Colors[c]; i++ {
			app.addEntity(c)
		}
	}

	for r := range settings.Rules {
		app.addRule(settings.Rules[r].Color1, settings.Rules[r].Color2, settings.Rules[r].Force, settings.Rules[r].EffectRange)
	}

	app.mutex.Unlock()
}

func (app *App) getSettings() *Settings {
	s := &Settings{
		Colors: map[string]uint{},
		Rules:  []Rule{},
	}

	for i := range app.entities {
		s.Colors[i] = uint(len(app.entities[i]))
	}

	for i := range app.rules {
		s.Rules = append(s.Rules, Rule{
			Color1:      app.rules[i].Color1,
			Color2:      app.rules[i].Color2,
			Force:       app.rules[i].Force,
			EffectRange: app.rules[i].EffectRange,
		})
	}

	return s
}

func (app *App) update() {
	//* Apply rules
	for u := range app.rules {
		for i := range app.entities[app.rules[u].Color1] {
			force := vec2f{}
			for j := range app.entities[app.rules[u].Color2] {
				if app.rules[u].Color1 != app.rules[u].Color2 || i != j {
					dx := app.entities[app.rules[u].Color1][i].pos.x - app.entities[app.rules[u].Color2][j].pos.x
					dy := app.entities[app.rules[u].Color1][i].pos.y - app.entities[app.rules[u].Color2][j].pos.y
					dist := math.Sqrt(dx*dx + dy*dy)

					if dist < app.rules[u].EffectRange {
						force.x -= app.rules[u].Force / dist * dx
						force.y -= app.rules[u].Force / dist * dy
					}
				}
			}

			app.entities[app.rules[u].Color1][i].speed.x *= app.globalLoss
			app.entities[app.rules[u].Color1][i].speed.y *= app.globalLoss
			app.entities[app.rules[u].Color1][i].speed.x += force.x
			app.entities[app.rules[u].Color1][i].speed.y += force.y

			app.entities[app.rules[u].Color1][i].newPos.x += app.entities[app.rules[u].Color1][i].speed.x
			app.entities[app.rules[u].Color1][i].newPos.y += app.entities[app.rules[u].Color1][i].speed.y

			// clamp
			if app.entities[app.rules[u].Color1][i].newPos.x < 0 {
				app.entities[app.rules[u].Color1][i].newPos.x = 0
				app.entities[app.rules[u].Color1][i].speed.x *= -1
			}
			if app.entities[app.rules[u].Color1][i].newPos.x > float64(app.size.x) {
				app.entities[app.rules[u].Color1][i].newPos.x = float64(app.size.x)
				app.entities[app.rules[u].Color1][i].speed.x *= -1
			}
			if app.entities[app.rules[u].Color1][i].newPos.y < 0 {
				app.entities[app.rules[u].Color1][i].newPos.y = 0
				app.entities[app.rules[u].Color1][i].speed.y *= -1
			}
			if app.entities[app.rules[u].Color1][i].newPos.y > float64(app.size.y) {
				app.entities[app.rules[u].Color1][i].newPos.y = float64(app.size.y)
				app.entities[app.rules[u].Color1][i].speed.y *= -1
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
	drawRect(app.ctx, "#ddd", vec2i{}, app.size)

	for color := range app.entities {
		for i := range app.entities[color] {
			drawRect(app.ctx, color, app.entities[color][i].pos.tovec2i(), vec2i{6, 6})
		}
	}
}

func (app *App) run() {
	app.mutex.Lock()
	app.update()
	app.render()
	app.mutex.Unlock()

	app.window.Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) any {
		app.run()
		return nil
	}))
}

var app *App = nil

func main() {

	js.Global().Set("startApp", js.FuncOf(func(this js.Value, args []js.Value) any {
		var err error

		if app == nil {
			app, err = newApp()
			if err != nil {
				log.Println(err)
				return err
			}
		}

		settings := &Settings{
			Colors: map[string]uint{},
			Rules:  []Rule{},
		}
		defaultSettings := &Settings{
			Colors: map[string]uint{
				"red":  20,
				"blue": 20,
			},
			Rules: []Rule{
				//{"red", "blue", -0.2, float64(app.size.x / 5.0)},
				{"blue", "blue", -0.1, float64(app.size.x / 5.0)},
				{"red", "red", 0.8, float64(app.size.x / 8.0)},
				{"red", "blue", 0.5, float64(app.size.x / 5.0)},
				{"blue", "red", -0.8, float64(app.size.x / 8.0)},
			},
		}

		if len(args) > 0 {
			if !strings.EqualFold(args[0].Type().String(), "string") {
				log.Printf("startApp expected string, got %s", args[0].Type().String())
				settings = defaultSettings
			}
			if err := json.Unmarshal([]byte(args[0].String()), settings); err != nil {
				log.Printf("Error parsing settings %s: %s", []byte(args[0].String()), err)
				settings = defaultSettings
			}
		} else {
			settings = defaultSettings
		}

		app.initialize(settings)
		app.run()

		return nil
	}))

	js.Global().Set("setSettings", js.FuncOf(func(this js.Value, args []js.Value) any {
		if app != nil {
			settings := &Settings{
				Colors: map[string]uint{},
				Rules:  []Rule{},
			}

			if len(args) > 0 {
				if !strings.EqualFold(args[0].Type().String(), "string") {
					log.Printf("startApp expected string, got %s\n", args[0].Type().String())
					return nil
				}
				if err := json.Unmarshal([]byte(args[0].String()), settings); err != nil {
					log.Printf("Error parsing settings %s: %s\n", []byte(args[0].String()), err)
					return nil
				}
			} else {
				log.Println("setSettings expects an argument, got none.")
				return nil
			}
			app.initialize(settings)
		}
		return nil
	}))

	js.Global().Set("getSettings", js.FuncOf(func(this js.Value, args []js.Value) any {
		b, err := json.Marshal(app.getSettings())
		if err != nil {
			log.Printf("Error marshalling settings: %s", err)
			return nil
		}
		return fmt.Sprintf("%s", b)
	}))

	<-make(chan bool) //wait
}
