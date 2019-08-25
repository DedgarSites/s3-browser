package main

import (
	"os"

	"github.com/dedgarsites/s3-browser/bucket"
	"github.com/dedgarsites/s3-browser/routers"
	"github.com/dedgarsites/s3-browser/tree"

	"path/filepath"
)

func main() {
	e := routers.Routers

	if localPort := os.Getenv("LOCAL_TESTING"); localPort != "" {
		go tree.ExampleTree()

		e.Logger.Info(e.Start(":" + localPort))
	} else {
		s3Out := bucket.ListContents()
		go tree.CreateTree(s3Out)

		fullPath := filepath.Join("/", os.Getenv("CERT_PATH"), "certificates", os.Getenv("CERT_FILENAME"))

		e.Logger.Info(e.StartTLS(":"+os.Getenv("TLS_PORT"), fullPath+".crt", fullPath+".key"))
	}
}
