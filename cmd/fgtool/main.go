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

		{
			Name:        "encode",
			Description: "turn decoded (potentially changed) files into game native format",
			Do:          encodeMain,
		},
	}

	subcmd.Run(cmds)
}

func decodeMain(args []string) {
	if err := doDecode(args); err != nil {
		log.Fatalf("decode error: %v", err)
	}
}

func encodeMain(args []string) {
	if err := doEncode(args); err != nil {
		log.Fatalf("encode error: %v", err)
	}
}
