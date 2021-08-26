package main

import (
    "fmt"
    "flag"
    "os"
    "log"
    "gitlab.com/aag031/cli_player/internal/common"
    //"errors"
)

func main() {
    nameOfFile, error := parseCommandLine()
    common.CheckErrorPanic(error)
    fileInfo, err := os.Stat(nameOfFile)

    if os.IsNotExist(err) {
        log.Fatal("Hmm ... File " + nameOfFile + " not found")
        os.Exit(1)
    }

    if fileInfo.IsDir() {
        log.Fatal("Hmm ... File " + nameOfFile + " should be a file not folder")
        os.Exit(1)
    }
    fmt.Println("Start Play Composition : " + nameOfFile)
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