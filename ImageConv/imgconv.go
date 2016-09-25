/* Convert a GIF/PNG/JPEG file to a GIF/PNG/JPEG file. */
package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	outFileUsage  = "output file format"
	destFileUsage = "file directory path to store output"
	nameFileUsage = "output file name"
)

var (
	outFileType  = flag.String("type", "", outFileUsage)
	destFilePath = flag.String("path", "", destFileUsage)
	outFileName  = flag.String("name", "", nameFileUsage)
)

func main() {
	/* get source image files from command line
	and flag values */
	sourceImages, err := parseflags()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	/* concurrently convert each source image file */
	var wg sync.WaitGroup
	for _, sImg := range sourceImages {
		wg.Add(1)
		go imgconv(sImg, &wg)
	}
	wg.Wait()
}

/*
parseflags extracts data from the command line and
sends user error messages if there are any problems. */
func parseflags() ([]string, error) {
	flag.Parse()
	sourceImages := flag.Args()
	if len(sourceImages) == 0 {
		fmt.Fprintf(os.Stderr, "Specify the source image files as arguments.\n")
	}
	*outFileType = strings.ToLower(*outFileType)
	if *outFileType == "" {
		return nil, fmt.Errorf("Specify the output file format (ie. JPEG, PNG, or GIF) with the -type flag.")
	} else if *outFileType != "jpeg" && *outFileType != "jpg" && *outFileType != "png" && *outFileType != "gif" {
		return nil, fmt.Errorf("%s is not one of the supported file formats: JPEG, PNG, or GIF\n", *outFileType)
	}
	return sourceImages, nil
}

/*
imgconv takes a source image file name (sImg)
and the address of a sync.WaitGroup (wg) since
this function is meant to be run as a goroutine.

imgconv opens and validates the source image file,
creates the destination image file specified by
the user, and converts the source image file to the
desired format. */
func imgconv(sImg string, wg *sync.WaitGroup) {
	defer wg.Done()

	/* open the source image file */
	in, err := os.Open(sImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", sImg, err)
		return
	}

	/* ensure source image file is supported */
	nameAndForm := strings.Split(filepath.Base(sImg), ".")
	name, format := nameAndForm[0], nameAndForm[1]
	if format != "jpeg" && format != "jpg" && format != "png" && format != "gif" {
		fmt.Fprintf(os.Stderr, "%s of format %s is not one of the supported file formats: JPEG, PNG, or GIF\n",
			sImg, format)
		in.Close()
		return
	}

	/* ensure file directory path ends with a / */
	if *destFilePath != "" {
		*destFilePath += "/"
		*destFilePath = filepath.Dir(*destFilePath) + "/" /* Dir() removes trailing slashes. */
	}

	/* create output file */
	var outFileName string
	switch *outFileType {
	case "jpeg", "jpg":
		outFileName = ".jpg"
	case "png":
		outFileName = ".png"
	case "gif":
		outFileName = ".gif"
	}
	outFileName = *destFilePath + name + outFileName
	out, err := os.Create(outFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		in.Close()
		return
	}

	/* convert the files */
	if err := toImg(in, out); err != nil {
		fmt.Fprintf(os.Stderr, "toImg: %v\n", err)
	}

	in.Close()
	out.Close()
}

/*
toImg decodes the image file (in) and encodes it into
the other image file (out), completing the format conversion. */
func toImg(in io.Reader, out io.Writer) error {
	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}
	switch *outFileType {
	case "jpeg", "jpg":
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		err = png.Encode(out, img)
	case "gif":
		err = gif.Encode(out, img, &gif.Options{NumColors: 256})
	}
	return err
}
