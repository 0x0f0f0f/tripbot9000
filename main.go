package main

import (

	// "gonum.org/v1/plot/vg/draw"
	// "gonum.org/v1/plot/vg/vgsvg"
	// "gonum.org/v1/plot/vg/vgeps"

	"log"
	"os"

	com "github.com/0x0f0f0f/tripbot9000/commands"
	gen "github.com/0x0f0f0f/tripbot9000/generator"
	gli "github.com/0x0f0f0f/tripbot9000/glitcher"
	// "strings"
)

func main() {

	cmds := []com.Command{
		gen.NewGemstoneCommand(),
		gen.NewMandelbrotCommand(),
		gen.NewMandelgemCommand(),
		gli.NewMatrixCommand(),
	}

	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("You must pass a subcommand") //TODO print usage
	}

	subcmd := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == subcmd {
			err := cmd.Init(args[1:])
			if err != nil {
				panic(err)
			}
			err = cmd.Run()
			if err != nil {
				panic(err)
			}

			return
		}
	}

	// TODO print usage
	log.Fatal("Subcommand not found")
}
