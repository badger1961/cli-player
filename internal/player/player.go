package player

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"azul3d.org/engine/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"gitlab.com/aag031/cli_player/internal/common"
)

const (
	COMMENT = "#"
)

type TDecodeFileFunc func(*os.File) (beep.StreamSeekCloser, beep.Format, error)

var controlTable map[string]TDecodeFileFunc

func init() {
	controlTable = make(map[string]TDecodeFileFunc)
	controlTable[".mp3"] = decodeMp3Composition
	controlTable[".wav"] = decodeWavComposition
	controlTable[".flac"] = decodeFlacComposition
	controlTable[".ogg"] = decodeOggComposition
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

func PlayFile(fileName string) error {
	var extension = filepath.Ext(fileName)
	if len(extension) == 0 {
		return errors.New("Hmm ... No extension for file")
	}

	fileHandle, error := os.Open(fileName)
	common.CheckErrorPanic(error)

	if decodeFileFuncPtr, ok := controlTable[extension]; ok {
		streamHandler, format, error := decodeFileFuncPtr(fileHandle)
		common.CheckErrorPanic(error)
		defer streamHandler.Close()
		size := format.SampleRate.D(streamHandler.Len())
		fmt.Printf("Start Play Composition : %v duration : %v", fileName, size.Round(time.Second))
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
		done := make(chan bool)
		//	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamHandler)}
		watcher := keyboard.NewWatcher()
		speaker.Play(beep.Seq(streamHandler, beep.Callback(func() {
			done <- true
		})))

		go func() {
			for {
				status := watcher.States()
				event := status[keyboard.P]
				if event == keyboard.Down {
					log.Println("P pressed")
				}
			}
		}()

		<-done
		return nil
	} else {
		return errors.New("Hmm ... " + extension + " is not supported ")
	}

	return nil
}

func decodeMp3Composition(fileHandle *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	streamHandler, format, error := mp3.Decode(fileHandle)
	common.CheckErrorPanic(error)
	return streamHandler, format, error
}

func decodeWavComposition(fileHandle *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	streamHandler, format, error := wav.Decode(fileHandle)
	common.CheckErrorPanic(error)
	return streamHandler, format, error
}

func decodeFlacComposition(fileHandle *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	streamHandler, format, error := flac.Decode(fileHandle)
	common.CheckErrorPanic(error)
	return streamHandler, format, error
}

func decodeOggComposition(fileHandle *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	streamHandler, format, error := vorbis.Decode(fileHandle)
	common.CheckErrorPanic(error)
	return streamHandler, format, error
}
