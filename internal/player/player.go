package player

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
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
	COMMENT      = "#"
	NOTIFICATION = "\rStart Play Composition : %v duration : %v:%v"
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

func loadPlayList(playListFileName string) (map[int]string, error) {
	idx := 0
	playList := make(map[int]string)
	file, err := os.Open(playListFileName)
	common.CheckErrorPanic(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileName := scanner.Text()
		if strings.HasPrefix(fileName, COMMENT) {
			continue
		}

		playList[idx] = fileName
		idx = idx + 1
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return playList, nil
}

func loadFolderToPlayList(folderName string) (map[int]string, error) {
	idx := 0
	playList := make(map[int]string)
	err := filepath.Walk(folderName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() {
			fmt.Println("Start folder " + info.Name())
			return nil
		}
		playList[idx] = path
		idx = idx + 1
		return nil
	})

	common.CheckErrorPanic(err)
	return playList, nil
}

func randomizePlayList(size int) []int {
	keyList := make([]int, 0, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		next := rand.Intn(size)
		keyList = append(keyList, next)
	}
	return keyList
}

func orderedPlayList(size int) []int {
	keyList := make([]int, 0, size)
	for i := 0; i < size; i++ {
		keyList = append(keyList, i)
	}
	return keyList
}

func PlayPlayList(playListName string, isRandomMode bool) error {
	fmt.Println("Start Play Compositions from playlist : " + playListName)
	playList, error := loadPlayList(playListName)
	common.CheckErrorPanic(error)
	playInternalPlayList(playList, isRandomMode)
	return nil
}

func PlayFolder(folderName string, isRandomMode bool) error {
	fmt.Println("Start Play Compositions from folder : " + folderName)
	playList, err := loadFolderToPlayList(folderName)
	common.CheckErrorPanic(err)
	playInternalPlayList(playList, isRandomMode)
	return nil
}

func playInternalPlayList(playList map[int]string, isRandomMode bool) error {
	var keyList []int
	if isRandomMode {
		keyList = randomizePlayList(len(playList))
	} else {
		keyList = orderedPlayList(len(playList))
	}
	for _, key := range keyList {
		PlayFile(playList[key])
	}
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
	fmt.Println()
	fmt.Printf(NOTIFICATION, fileName, size.Round(time.Second), 0)
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
	seconds := time.Tick(time.Second)
endPlay:
	for {
		select {
		case event := <-keysEvents:
			if event.Key == keyboard.KeyEsc {
				os.Exit(1)
			}
			if event.Key == keyboard.KeySpace {
				speaker.Lock()
				ctrl.Paused = !ctrl.Paused
				speaker.Unlock()
			}
			if event.Key == keyboard.KeyArrowRight {
				break endPlay
			}
		case _ = <-seconds:
			pos := format.SampleRate.D(streamHandler.Position())
			fmt.Printf(NOTIFICATION, fileName, size.Round(time.Second), pos.Round(time.Second))
		case isEnd := <-done:
			if isEnd {
				break endPlay
			}
			fmt.Println()
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
