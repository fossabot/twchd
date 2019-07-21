// +build ignore

package main

import (
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	if err := vfsgen.Generate(
		http.Dir("./assets"),
		vfsgen.Options{}); err != nil {
		panic(err)
	}
}
