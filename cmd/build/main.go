package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

var logger = log.New(os.Stderr, "", 0)

type ChaincodeMetadata struct {
	Type string `json:"type"`
}

type Connection struct {
	Address     string `json:"address"`
	DialTimeout string `json:"dial_timeout"`
	TLS         bool   `json:"tls_required"`
	ClientAuth  bool   `json:"client_auth_required"`
	RootCert    string `json:"root_cert"`
	ClientKey   string `json:"client_key"`
	ClientCert  string `json:"client_cert"`
}

func main() {
	logger.Println("::Build")

	if err := run(); err != nil {
		logger.Printf("::Error: %v\n", err)
		os.Exit(1)
	} else {
		logger.Printf("::Type detected as external")
	}
}

func run() error {
	sourceDir, metadataDir, outputDir := os.Args[1], os.Args[2], os.Args[3]
	connectionSrcFile := filepath.Join(sourceDir, "/connection.json")
	metadataFile := filepath.Clean(filepath.Join(metadataDir, "metadata.json"))
	connectionDestFile := filepath.Join(outputDir, "/connection.json")
	metainfoSrcDir := filepath.Join(sourceDir, "META-INF")
	metainfoDestDir := filepath.Join(outputDir, "META-INF")
	// Process and check the metadata file, then copy to the output location
	metadataFileContents, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		return err
	}

	var metadata ChaincodeMetadata
	if err := json.Unmarshal(metadataFileContents, &metadata); err != nil {
		return err
	}

	if strings.ToLower(metadata.Type) != "external" {
		return fmt.Errorf("chaincode type should be external, it is %s", metadata.Type)
	}

	if err := copy.Copy(metadataDir, outputDir); err != nil {
		return fmt.Errorf("failed to copy build metadata folder: %s", err)
	}

	if _, err := os.Stat(metainfoSrcDir); !os.IsNotExist(err) {
		if err := copy.Copy(metainfoSrcDir, metainfoDestDir); err != nil {
			return fmt.Errorf("failed to copy build META-INF folder: %s", err)
		}
	}

	// Process and update the connections file
	fileInfo, err := os.Stat(connectionSrcFile)
	if err != nil {
		return fmt.Errorf("connection.json not found in source folder: %s", err)
	}

	connectionFileContents, err := ioutil.ReadFile(connectionSrcFile)
	if err != nil {
		return err
	}

	var connectionData Connection
	if err := json.Unmarshal(connectionFileContents, &connectionData); err != nil {
		return err
	}

	if err := updateConnectionData(&connectionData); err != nil {
		return err
	}

	updatedConnectionBytes, err := json.Marshal(connectionData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated connection.json file: %s", err)
	}

	err = ioutil.WriteFile(connectionDestFile, updatedConnectionBytes, fileInfo.Mode())
	if err != nil {
		return err
	}

	return nil

}

func updateConnectionData(metadata *Connection) error {
	// do nothing for the moment
	return nil
}
