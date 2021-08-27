package player

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"gitlab.com/aag031/cli_player/internal/common"
	"os"
	"path/filepath"
	"time"
	"strings"
)

const (
	COMMENT = "#"
)

type TPlayFileFunc func(string) error
var controlTable  map[string]TPlayFileFunc

func init() {
	controlTable = make(map[string]TPlayFileFunc)
	controlTable[".mp3"] = playMp3File
	controlTable[".wav"] = playWavFile
}

func PlayPlayList(playListName string) error {
	fmt.Println("Start Play Compositions from playlist : " + playListName)
	file, err := os.Open(playListName)
	common.CheckErrorPanic(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileName := scanner.Text()
		if strings.HasPrefix(fileName, COMMENT) {
			continue
		}
		PlayFile(fileName)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func PlayFolder(folderName string) error {
	fmt.Println("Start Play Compositions from folder : " + folderName)
	err := filepath.Walk(folderName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() {
			fmt.Println("Start folder " + info.Name())
			return nil
		}
		fmt.Printf("Play name: %s\n", path)
 		errFile := PlayFile(path)
 		common.CheckErrorNoPanic(errFile)
		return nil
	})

	common.CheckErrorPanic(err)
	return nil
}

func PlayFile (fileName string) error {
	fmt.Println("Start Play Composition : " + fileName)
	var extension = filepath.Ext(fileName)
	if len (extension) == 0 {
		return errors.New("Hmm ... No extension for file")
	}
	if playFileFuncPtr, ok := controlTable[extension]; ok {
		playFileFuncPtr(fileName)
	} else {
		return errors.New("Hmm ... " + extension + " is not supported ")
	}

	return nil
}

func playWavFile(fileName string) error {
	fileHandle, error := os.Open(fileName)
	common.CheckErrorPanic(error)

	streamHandler, format, error := wav.Decode(fileHandle)
	common.CheckErrorPanic(error)
	defer streamHandler.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamHandler, beep.Callback(func() {
		done <-true
	})))
	<-done
	return nil
}

func playMp3File(fileName string) error {
	fileHandle, error := os.Open(fileName)
	common.CheckErrorPanic(error)

	streamHandler, format, error := mp3.Decode(fileHandle)
	common.CheckErrorPanic(error)
	defer streamHandler.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamHandler, beep.Callback(func() {
		done <-true
	})))
	<-done
	return nil
}

func decodeComposition(fileName string)  (beep.StreamSeekCloser, beep.Format, error)  {

}
