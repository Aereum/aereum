package logger

import "fmt"

func MustOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func LogError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
