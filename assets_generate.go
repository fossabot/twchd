package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var err = vfsgen.Generate(
		http.Dir("./assets"),
		vfsgen.Options{
			Filename: "assets.go",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
}
