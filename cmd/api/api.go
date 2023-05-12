package main

import (
	"github.com/Parsa-Sedigh/rebottle/internal/app"
)

func main() {
	api := app.NewApp()
	api.Start()
}
