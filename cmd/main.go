package main

import (
	"os"

	"github.com/ygortgaleno/cloudwalk-test/internals/services"
)

func main() {
	svc := services.QuakeLogParserService{}
	r := svc.Exec("../qgames.log")
	os.WriteFile("../result.json", r, 0644)
}
