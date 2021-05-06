package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	register(r)
	log.Printf("run toad ocr api service...")
	if err := r.Run(":18889"); err != nil {
		log.Printf("run api service fail!")
	}
}
