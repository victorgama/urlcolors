package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"

	"io/ioutil"

	"github.com/aybabtme/rgbterm"
	"github.com/victorgama/colorarty"
)

func printUsageAndExit() {
	fmt.Println("Usage: urlcolors http(s)://url")
	fmt.Println()
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		printUsageAndExit()
		return
	}

	_, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		printUsageAndExit()
		return
	}

	tmpfile, err := ioutil.TempFile("", "urlcolors")
	if err != nil {
		println(err)
		return
	}

	defer os.Remove(tmpfile.Name())

	resp, err := http.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	img, _, err := image.Decode(tmpfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	result := colorarty.Analyse(img)
	if result != nil {
		fmt.Println()

		fmt.Printf("Background: %s\n", colorToRGB(result.BackgroundColor))
		fmt.Printf("   Primary: %s\n", colorToRGB(result.PrimaryColor))
		fmt.Printf(" Secondary: %s\n", colorToRGB(result.SecondaryColor))
		fmt.Printf("    Detail: %s\n", colorToRGB(result.DetailColor))
		fmt.Println()
	} else {
		fmt.Println("Color analysis failed.")
	}
}

func colorToRGB(c *color.Color) string {
	cr, cg, cb, _ := (*c).RGBA()
	r := float64(cr)
	g := float64(cg)
	b := float64(cb)
	r /= 0x101
	g /= 0x101
	b /= 0x101
	hex := fmt.Sprintf("#%02x%02x%02x", uint8(r), uint8(g), uint8(b))
	return fmt.Sprintf("%s (Hex: %s RGB: %.0f, %.0f, %.0f)", rgbterm.FgString("█", uint8(r), uint8(g), uint8(b)), hex, r, g, b)
}
