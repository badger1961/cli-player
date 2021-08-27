package main

import (
	"github.com/DavidGamba/go-getoptions"
	"fmt"
	"gitlab.com/aag031/cli_player/internal/common"
	"gitlab.com/aag031/cli_player/internal/player"
	"os"
)

type playMode int

const (
	fileMode= playMode(iota)
	folderMode
	playListMode
	unknownMode
)

const VERSION = "0.2.0"

func main() {
	fileName, mode := parseCommandLine()
	if mode == fileMode {
		fileError := common.CheckInputFile(fileName)
		common.CheckErrorPanic(fileError)
		fmt.Println("Start Play Composition : " + fileName)
		errorPlayer := player.PlayFile(fileName)
		common.CheckErrorPanic(errorPlayer)
	}
	if mode == folderMode {
		fileError := common.CheckInputFolder(fileName)
		common.CheckErrorPanic(fileError)
		fmt.Println("Start Play Compositions from folder: " + fileName)
		errorPlayer := player.PlayFolder(fileName)
		common.CheckErrorPanic(errorPlayer)
	}

}

func parseCommandLine() (string, playMode) {
	var fileName string
	var folderName string
	opt := getoptions.New()
	opt.Bool("help", false, opt.Alias("h", "?"))
	opt.StringVarOptional(&fileName, "file", "", opt.Description("Name of file with composition for playing"), opt.Alias("f"))
	opt.StringVarOptional(&folderName, "dir", "", opt.Description("Name of folder with compositions for playing"), opt.Alias("d"))
	_, err := opt.Parse(os.Args[1:])
	if opt.Called("help") {
		fmt.Fprintf(os.Stderr, opt.Help())
		os.Exit(1)
	}
	if err  != nil {
		fmt.Printf("ERROR %s\n\n", err)
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if fileName == "" && folderName == "" {
		fmt.Printf("ERROR %s\n\n", "fileName or folderName should be passed")
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if opt.Called("file") && opt.Called("dir") {
		fmt.Printf("ERROR %s\n\n", "fileName  and folderName in the same time should not be passed")
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if len(fileName) > 0 {
		return fileName, fileMode
	}
	if len(folderName) > 0 {
		return folderName, folderMode
	}
	return "", unknownMode
}
