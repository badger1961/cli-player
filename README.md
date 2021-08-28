# cli-player

This is simple audio player that can be started from command line. It bases on GoLang library "github.com/faiface/beep"
and support format that supported by this library.

Usage

SYNOPSIS:
    cli-player [--dir|-d <string>] [--file|-f <string>] [--help|-h|-?]
               [--plist|-l <string>] [--version|-v] [<args>]

OPTIONS:
    --dir|-d <string>      Name of folder with compositions for playing (default: "")

    --file|-f <string>     Name of file with composition for playing (default: "")

    --help|-h|-?           print this text and exit

    --plist|-l <string>    Name of file with playlist for playing (default: "")

    --version|-v           print version and exit

    playlist is text file each line of this file contains path to audio file. For example


./testdata/1.mp3
./testdata/1.flac

Comment is marked by #

