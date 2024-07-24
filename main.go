package main

import (
	"THR/Database"
	route "THR/Route"
	"fmt"
	"log"
)

func main() {
	
	err := Database.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer Database.DBConnect.Close()
	
	fmt.Println("MySQL connected")
	route.RunServer()
	
}