package main

import (
	"github.com/mappedbyte/server-agent/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd Execute err: %v", err)
	}

}
