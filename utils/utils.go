package utils

import (
	"log"
	"os"

	"gocv.io/x/gocv"
)

func LoadImg(fileName string) gocv.Mat {
	_, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
	}
	img := gocv.IMRead(fileName, gocv.IMReadColor)
	return img
}

func SaveImg(fileName string, img *gocv.Mat) {
	gocv.IMWrite(fileName, *img)
}
