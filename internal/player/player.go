package player

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
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

	decodeFileFuncPtr, ok := controlTable[extension]
	if !ok {
		return errors.New("Hmm ... " + extension + " is not supported ")
	}
	streamHandler, format, error := decodeFileFuncPtr(fileHandle)
	common.CheckErrorPanic(error)
	defer streamHandler.Close()
	size := format.SampleRate.D(streamHandler.Len())
	fmt.Printf("Start Play Composition : %v duration : %v", fileName, size.Round(time.Second))
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamHandler)}
	ctrl.Paused = false
	done := make(chan bool, 1)

	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		done <- true
	})))

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
endPlay:
	for {
		select {
		case event := <-keysEvents:
			if event.Key == keyboard.KeyEsc {
				os.Exit(1)
			}
			if event.Key == keyboard.KeySpace {
				speaker.Lock()
				fmt.Println(ctrl.Paused)
				ctrl.Paused = !ctrl.Paused
				fmt.Println(ctrl.Paused)
				speaker.Unlock()
			}
		case isEnd := <-done:
			fmt.Println(isEnd)
			if isEnd {
				break endPlay
			}
		default:
			continue
		}
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
