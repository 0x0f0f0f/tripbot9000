package glitcher

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/repr"
	"gonum.org/v1/gonum/mat"
)

type MatrixCommand struct {
	fs          *flag.FlagSet
	infilename  string
	outfilename string
	invert      bool
	gauss       bool
}

func NewMatrixCommand() *MatrixCommand {
	c := &MatrixCommand{
		fs: flag.NewFlagSet("matrix", flag.ContinueOnError),
	}

	now := time.Now()
	defaultfile := fmt.Sprintf("output-%s.png", now.Format("Jan-2-2006-15-04-05"))
	c.fs.StringVar(&c.outfilename, "o", defaultfile, "output image file name or directory")
	c.fs.StringVar(&c.infilename, "i", "", "input image file name")
	c.fs.BoolVar(&c.invert, "invert", false, "invert matrix")
	c.fs.BoolVar(&c.gauss, "gauss", false, "perform gaussian reduction")

	stat, err := os.Stat(c.outfilename)
	if err == nil {
		if stat.IsDir() {
			c.outfilename = filepath.Join(c.outfilename, defaultfile)
		}
	}

	return c
}

func (c *MatrixCommand) Init(args []string) error {
	err := c.fs.Parse(args)

	repr.Println(*c)

	if strings.TrimSpace(c.infilename) == "" {
		return errors.New("please specify an input file")
	}
	return err
}

func (c *MatrixCommand) Name() string {
	return c.fs.Name()
}

func (c *MatrixCommand) Run() error {

	inputImageFile, err := os.Open(c.infilename)
	if err != nil {
		return err
	}
	defer inputImageFile.Close()

	img, err := png.Decode(inputImageFile)
	if err != nil {
		return err
	}
	fmt.Println(&img)

	img_w := img.Bounds().Max.X
	img_h := img.Bounds().Max.Y

	pixData := make([]float64, img_w*img_h)

	for y := img.Bounds().Min.Y; y < img_h; y++ {
		for x := img.Bounds().Min.X; x < img_w; x++ {
			col := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			pixData[(y*img_w)+x] = float64(col.Y)
		}
	}

	m := mat.NewDense(img_h, img_w, pixData)

	if c.invert {
		m.Inverse(m)
	}

	if c.gauss {
		var lu mat.LU
		lu.Factorize(m)
		dest := lu.UTo(nil)
		m.Copy(dest)
	}

	// Convert back into an image
	outImg := image.NewRGBA(image.Rect(0, 0, img_w, img_h))
	out, err := os.Create(c.outfilename)
	for y := outImg.Bounds().Min.Y; y < img_h; y++ {
		for x := outImg.Bounds().Min.X; x < img_w; x++ {
			pix := m.At(y, x)
			col := color.Gray{Y: uint8(pix)}
			outImg.Set(x, y, col)
		}
	}
	//fmt.Printf("u= %v\n", mat.Formatted(m, mat.Prefix("    ")))

	png.Encode(out, outImg)
	out.Close()
	return nil
}
