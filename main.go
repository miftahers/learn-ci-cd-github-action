package main

import (
	"praktikum/config"
	"praktikum/routes"
)

func main() {

	db := config.InitDB()

	e := routes.Init(db)

	e.Start(config.APIPort)
}
