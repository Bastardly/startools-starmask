package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"

	"github.com/sqweek/dialog"
	"golang.org/x/image/tiff"
)

const fancyMessage = "Oh no! We could not read that file. Are you sure it is a 16bit tif?"

func selectFile() (image.Image, string, error) {
	uri, _ := dialog.File().Filter("16bit tif file", "tif", "tiff").Load()

	fmt.Println(uri)
	data, err := ioutil.ReadFile(uri)

	if err != nil {
		fmt.Println("Unable to read that file")
		dialog.Message(fancyMessage)

	}

	img, err2 := tiff.Decode(bytes.NewReader(data))

	return img, uri, err2
}
