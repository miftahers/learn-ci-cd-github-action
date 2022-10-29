package main

import (
	"praktikum/config"
	"praktikum/routes"
)

func main() {

	db := config.InitDB()

	e := routes.Init(db)

	err := e.Start(config.APIPort)
	if err != nil {
		panic(err)
	}
}
