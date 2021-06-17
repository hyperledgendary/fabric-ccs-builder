package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	logger.Println("::Detect 002")

	if err := run(); err != nil {
		logger.Printf("::Error: %v\n", err)
		os.Exit(1)
	} else {
		logger.Printf("::Type detected as external")
	}
}

type ChaincodeMetadata struct {
	Type string `json:"type"`
}

func run() error {
	if len(os.Args) > 2 {
		chaincodeMetaData := os.Args[2]
		mdbytes, err := ioutil.ReadFile(filepath.Join(chaincodeMetaData, "metadata.json"))
		if err != nil {
			return err
		}

		var metadata ChaincodeMetadata
		err = json.Unmarshal(mdbytes, &metadata)
		if err != nil {
			return err
		}

		switch strings.ToLower(metadata.Type) {
		case "external":
			return nil
		default:
			return fmt.Errorf("chaincode type not supported: %s", metadata.Type)
		}

	} else {
		return fmt.Errorf("Too few arguments")
	}

}
