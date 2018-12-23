package main

import (
	"00-newapp-template/internal"
	"00-newapp-template/internal/pkg"
)

func main() {
	config := pkg.NewConfig()
	a := internal.NewApp(config)

	a.InvokeCLI()
	return
}
