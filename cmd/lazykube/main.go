package main

import "github.com/TNK-Studio/lazykube/pkg/app"

func main() {
	lazykube := app.NewApp()
	defer lazykube.Stop()
	lazykube.Run()
}
