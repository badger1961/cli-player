package main

import (
	"fmt"
	"os"

	"github.com/DavidGamba/go-getoptions"
	"gitlab.com/aag031/cli_player/internal/common"
	"gitlab.com/aag031/cli_player/internal/player"
)

type playMode int

const (
	fileMode = playMode(iota)
	folderMode
	playListMode
	unknownMode
)

const VERSION = "1.4.0"

func main() {
	fileName, mode, isRandomMode := parseCommandLine()
	if mode == fileMode {
		fileError := common.CheckInputFile(fileName)
		common.CheckErrorPanic(fileError)
		errorPlayer := player.PlayFile(fileName)
		common.CheckErrorPanic(errorPlayer)
	}
	if mode == folderMode {
		fileError := common.CheckInputFolder(fileName)
		common.CheckErrorPanic(fileError)
		errorPlayer := player.PlayFolder(fileName, isRandomMode)
		common.CheckErrorPanic(errorPlayer)
	}
	if mode == playListMode {
		fileError := common.CheckInputFile(fileName)
		common.CheckErrorPanic(fileError)
		errorPlayer := player.PlayPlayList(fileName, isRandomMode)
		common.CheckErrorPanic(errorPlayer)
	}
}

func parseCommandLine() (string, playMode, bool) {
	var fileName string
	var folderName string
	var playListName string
	opt := getoptions.New()
	opt.Bool("help", false, opt.Alias("h", "?"))
	opt.Bool("version", false, opt.Alias("v"))
	opt.StringVarOptional(&fileName, "file", "", opt.Description("Name of file with composition for playing"), opt.Alias("f"))
	opt.StringVarOptional(&folderName, "dir", "", opt.Description("Name of folder with compositions for playing"), opt.Alias("d"))
	opt.StringVarOptional(&playListName, "plist", "", opt.Description("Name of file with playlist for playing"), opt.Alias("l"))
	opt.Bool("random", false, opt.Description("Random mode of playing or not"), opt.Alias("r"))
	_, err := opt.Parse(os.Args[1:])
	if opt.Called("help") {
		fmt.Printf(opt.Help())
		os.Exit(1)
	}
	if opt.Called("version") {
		fmt.Printf("Version : " + VERSION)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("ERROR %s\n\n", err)
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if fileName == "" && folderName == "" && playListName == "" {
		fmt.Printf("ERROR %s\n\n", "fileName or folderName ot playlist should be passed")
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if opt.Called("file") && opt.Called("dir") && opt.Called("plist") {
		fmt.Printf("ERROR %s\n\n", "fileName and folderName and playlist in the same time should not be passed")
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if opt.Called("file") && opt.Called("dir") {
		fmt.Printf("ERROR %s\n\n", "fileName and folderName in the same time should not be passed")
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if opt.Called("file") && opt.Called("plist") {
		fmt.Printf("ERROR %s\n\n", "fileName andplaylist in the same time should not be passed")
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if opt.Called("dir") && opt.Called("plist") {
		fmt.Printf("ERROR %s\n\n", "folderName and playlist in the same time should not be passed")
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}
	if opt.Called("random") && opt.Called("file") {
		fmt.Printf("ERROR %s\n\n", "fileName and random option is not compatible")
		fmt.Printf(opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)

	}
	var isRandomMode bool
	if opt.Called("random") {
		isRandomMode = true
	}
	if opt.Called("file") {
		return fileName, fileMode, false
	}
	if opt.Called("dir") {
		return folderName, folderMode, isRandomMode
	}
	if opt.Called("plist") {
		return playListName, playListMode, isRandomMode
	}
	return "", unknownMode, false
}
