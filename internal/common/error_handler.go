package common

import "log"

func CheckErrorPanic(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func CheckErrorNoPanic(err error) {
	if err != nil {
		log.Println(err)
	}
}
