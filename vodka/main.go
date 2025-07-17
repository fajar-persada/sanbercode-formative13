package main

import (
	"vodka/routers"
)

func main() {
	var PORT = ":8080"

	routers.StartServer().Run(PORT)
}
