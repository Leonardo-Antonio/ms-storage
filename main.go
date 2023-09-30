package main

import (
	_ "github.com/Leonardo-Antonio/ms-storage/config"
	"github.com/Leonardo-Antonio/ms-storage/internal/router"
)

func main() {
	router.New().Run()
}
