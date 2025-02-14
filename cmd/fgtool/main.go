package main

import (
	"log"

	"github.com/cespare/subcmd"
)

func main() {
	log.SetFlags(0)

	cmds := []subcmd.Command{
		{
			Name:        "decode",
			Description: "turn game files into easy-to-edit format",
			Do:          decodeMain,
		},
	}

	subcmd.Run(cmds)
}

func decodeMain(args []string) {
	if err := doDecode(args); err != nil {
		log.Fatalf("decode error: %v", err)
	}
}
