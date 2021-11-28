package main

import (
	"bot/envars"
	"bot/app"
)

func init() {
	envars.LoadEnvs()
}

func main() {
	app.Run()
}