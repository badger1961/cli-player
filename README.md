# cli-player

This is simple audio player that can be started from command line. It bases on GoLang library "github.com/faiface/beep"
and support format that supported by this library.

Usage

SYNOPSIS:                                                                                  
    cli-player.exe [--dir|-d <string>] [--file|-f <string>] [--help|-h|-?]                 
                   [--plist|-l <string>] [--random|-r] [--version|-v] [<args>]             
                                                                                           
OPTIONS:                                                                                   
    --dir|-d <string>      Name of folder with compositions for playing (default: "")      
                                                                                           
    --file|-f <string>     Name of file with composition for playing (default: "")         
                                                                                           
    --help|-h|-?           (default: false)                                                
                                                                                           
    --plist|-l <string>    Name of file with playlist for playing (default: "")            
                                                                                           
    --random|-r            Random mode of playing or not (default: false)                  
                                                                                           
    --version|-v           (default: false)                                                
                                                                                           
    playlist is text file each line of this file contains path to audio file. For example


1. ./testdata/1.mp3
1. ./testdata/1.flac


Comment is marked by #

