package main

import (
	"flag"
	"fmt"
	"gitlab.com/aag031/cli_player/internal/common"
	"gitlab.com/aag031/cli_player/internal/player"
	"os"
	//"errors"
)

const VERSION = "0.1.0"

func main() {
	fileName, error := parseCommandLine()
	common.CheckErrorPanic(error)
	fileError := common.CheckInputFile(fileName)
	common.CheckErrorPanic(fileError)
	fmt.Println("Start Play Composition : " + fileName)
	errorPlayer := player.PlayFile(fileName)
	common.CheckErrorPanic(errorPlayer)
}

func parseCommandLine() (string, error) {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdout, "Usage: %s --file<name_of_composition>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	fileNamePtr := flag.String("file", "", "name of file with composition")
	flag.Parse()
	return *fileNamePtr, nil
}
