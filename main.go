package main

import (
	"log"
	"syscall/js"
)

func main() {
	js.Global().Set("startApp", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Println("Hello world !")
		return nil
	}))

	<-make(chan bool) //wait
}
