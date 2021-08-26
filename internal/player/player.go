package player

import (
	"errors"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"gitlab.com/aag031/cli_player/internal/common"
	"os"
	"path/filepath"
	"time"
)
type TPlayFileFunc func(string) error
var controlTable  map[string]TPlayFileFunc

func init() {
	controlTable = make(map[string]TPlayFileFunc)
	controlTable[".mp3"] = playMp3File
	controlTable[".wav"] = playWavFile
}

func PlayFile (fileName string) error {
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
	return errors.New("Hmm ... WAV not implemented")
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
