package main

import (
	"os"

	"github.com/dedgarsites/s3-browser/bucket"
	"github.com/dedgarsites/s3-browser/routers"
	"github.com/dedgarsites/s3-browser/tree"
)

func main() {
	s3Out := bucket.ListContents()

	go tree.CreateTree(s3Out)

	//go tree.OldTree()

	e := routers.Routers
	if localPort := os.Getenv("LOCAL_TESTING"); localPort != "" {
		e.Logger.Info(e.Start(":" + localPort))
	} else {
		e.Logger.Info(e.StartTLS(":8443", "/cert/lego/certificates/dashboard.crt", "/cert/lego/certificates/dashboard.key"))
	}
}
