package main

import (
	"afire/internal/app/manager/service"
	"log"
)

func main()  {
	e := service.Start()
	if e != nil {
		log.Fatalln(e.Error())
	}
}
