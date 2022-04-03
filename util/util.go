package util

import "log"

func checkError(err error) {
	if err != nil {
		log.Println("CheckError:", err)
	}
}
