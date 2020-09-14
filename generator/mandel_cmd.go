package generator

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/0x0f0f0f/tripbot9000/util"
	"gonum.org/v1/plot/vg"
)

type MandelbrotCommand struct {
	fs        *flag.FlagSet
	filename  string
	img_w     vg.Length
	img_h     vg.Length
	img_w_str string
	img_h_str string
	dpi       int
}

func NewMandelbrotCommand() *MandelbrotCommand {
	c := &MandelbrotCommand{
		fs: flag.NewFlagSet("mandel", flag.ContinueOnError),
	}

	now := time.Now()
	defaultfile := fmt.Sprintf("output-%s.png", now.Format("Jan-2-2006-15-04-05"))
	c.fs.StringVar(&c.filename, "o", defaultfile, "output image file name or directory")
	c.fs.StringVar(&c.img_w_str, "w", "80cm", "image width")
	c.fs.StringVar(&c.img_h_str, "h", "80cm", "image height")
	c.fs.IntVar(&c.dpi, "dpi", 120, "resolution in dpi")

	return c
}

func (c *MandelbrotCommand) Init(args []string) error {
	err := c.fs.Parse(args)

	c.img_w, err = vg.ParseLength(c.img_w_str)
	if err != nil {
		log.Fatalf("could not parse width: %+v", err)
	}
	c.img_h, err = vg.ParseLength(c.img_h_str)
	if err != nil {
		log.Fatalf("could not parse height: %+v", err)
	}

	return err
}

func (c *MandelbrotCommand) Name() string {
	return c.fs.Name()
}

func (c *MandelbrotCommand) Run() error {
	img := RandomMandelbrot(c.filename, c.img_w, c.img_h, c.dpi)
	util.SavePNG(img, c.filename)
	return nil
}
