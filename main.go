package main

import (
	"context"
	"go-http-server-template/cmd"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := cmd.Run(os.Args, os.Stdout); err != nil {
		switch err {
		case context.Canceled:
			// not considered error
		case http.ErrServerClosed:
			// not considered error
		default:
			log.Fatalf("could not run application: %v", err)
		}
	}
}
