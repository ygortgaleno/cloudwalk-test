package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/ygortgaleno/cloudwalk-test/internals/services"
)

func main() {
	var outputFile string
	var duration, help bool
	flag.BoolVar(&help, "help", false, "Help flag to see usage commands")
	flag.StringVar(&outputFile, "out", "", "Output filepath to put the parsed result")
	flag.BoolVar(&duration, "duration", false, "Mesure duration of parser execution")
	flag.Parse()

	if help {
		helpInfo()
		return
	}

	if len(flag.Args()) == 0 {
		log.Fatal("The input file is required as first argument")
	}

	quakelogFilepath := flag.Args()[0]

	start := time.Now()

	var w io.Writer
	w = os.Stdout
	if len(outputFile) > 0 {
		var err error
		w, err = os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	jsonEncoder := json.NewEncoder(w)
	svc := services.QuakeLogParserService{}
	result := svc.Exec(quakelogFilepath)
	jsonEncoder.Encode(result)

	if duration {
		fmt.Println(time.Since(start))
	}
}

func helpInfo() {
	fmt.Println(`This program is a parser from quake 3 game that get kill events from log and agregate into a json for each mach.
		
Usage:

  parser_quake_log [-out=out_filepath -duration] [quake_log_filepath]

Flags:

  out   	string 		Filepath to get the result of parser(stdout is default)
  duration	boolean		Mesure of parser duration
  help		boolean		Helps to se the usage of program`)
}
