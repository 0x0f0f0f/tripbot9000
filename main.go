package main

import (

	// "gonum.org/v1/plot/vg/draw"
	// "gonum.org/v1/plot/vg/vgsvg"
	// "gonum.org/v1/plot/vg/vgeps"

	"log"
	"os"

	com "github.com/0x0f0f0f/tripbot9000/commands"
	gen "github.com/0x0f0f0f/tripbot9000/generator"
	// "strings"
)

func main() {

	cmds := []com.Command{
		gen.NewGemstoneCommand(),
		gen.NewMandelbrotCommand(),
		gen.NewMandelgemCommand(),
	}

	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("You must pass a subcommand") //TODO print usage
	}

	subcmd := args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcmd {
			cmd.Init(args[1:])
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
		}
	}

}
