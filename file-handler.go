package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"starkiller/createfile"
	"starkiller/imagerunner"

	"github.com/sqweek/dialog"
	"golang.org/x/image/tiff"
)

const fancyMessage = "We regret to inform you that the lake is full of cats, and that we were unable to read your file."

func selectFile() {
	uri, _ := dialog.File().Filter("16bit tif file", "tif", "tiff").Load()

	fmt.Println(uri)

	// f, err := os.OpenFile(uri, os.O_RDONLY, 0644)

	data, err := ioutil.ReadFile(uri)

	if err != nil {
		fmt.Println("Unable to read that file")
		dialog.Message(fancyMessage)

	} else {

		img, err := tiff.Decode(bytes.NewReader(data))

		if err != nil {
			fmt.Println(err)
			dialog.Message(fancyMessage)

		} else {
			pixels, width, height := imagerunner.Start(img)
			createfile.CreateAlpha(pixels, uri, width, height)
			createfile.CreatePreview(pixels, uri, width, height)
		}

	}
}
