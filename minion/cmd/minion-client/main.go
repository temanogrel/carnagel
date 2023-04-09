package main

import (
	"fmt"
	"os"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain/upstore"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	client := upstore.NewUpstoreClient("admin@camgirlcaps.com", "k6YiN9Ul4Tb73Bnmyu9Y6045b6x342vRtWeRR1", logger)

	fp, err := os.Open(".gitignore")
	if err != nil {
		panic(err)
	}

	if hash, err := client.Upload("Makefile", fp); err != nil {
		panic(err)
	} else {
		fmt.Println(hash)
	}

	fp.Close()

	fp, err = os.Open(".gitignore")
	if err != nil {
		panic(err)
	}

	if hash, err := client.Upload("Makefite", fp); err != nil {
		panic(err)
	} else {
		fmt.Println(hash)
	}

	fp.Close()
}
