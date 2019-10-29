package microbusiness

import (
	"log"
	"os"
)

func GetHomeDirectoryPath() string {
	homePath, err := os.UserHomeDir()
	HandleError(err)
	log.Print(homePath)
	if homePath != "" {
		log.Print("linux mode")
		return homePath
	}
	log.Print("windows mode")
	return os.Getenv("USERPROFILE")
}

func HandleError(err error) {
	if err != nil {
		log.Print(err)
		panic(err.Error())
	}
}

func Int32Ptr(i int32) *int32 {
	return &i
}
