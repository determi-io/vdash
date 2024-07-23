package main

import (
	"context"
	"fmt"
	"log"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	go runWebsocketServer()
	log.Print("before running term")
	go runTerm()
	log.Print("after running term")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {

	server := New(&Config{
		Host: "localhost",
		Port: "3333",
	})
	go server.Run()

	return fmt.Sprintf("Hello %s, It's show time!", name)
}
