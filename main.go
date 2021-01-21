package main

import (
	"Mmx/Router"
	"Mmx/Service"
)

func main() {
	Service.InitDatabase()
	Router.InitRouter()
}
