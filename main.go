package main

import "github.com/aldoramirezmartinez/fiber-api/config"

func main() {
	config.LoadEnv()
	app := NewApp()
	app.Run()
}
