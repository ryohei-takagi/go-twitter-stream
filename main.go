package main

import (
	"github.com/flatnyat/go-twitter-stream/app"
)

func main() {
	twitter := app.NewTwitterClient()
	app.NewRouter(8081, &twitter).Run()
}