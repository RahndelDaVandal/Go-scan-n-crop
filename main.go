package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

func main() {
  fileName := "./imgs-test/middle.JPG"
  scan := gocv.IMRead(fileName, gocv.IMReadColor)

  fmt.Printf("%+v\n", scan)

  scanGray := gocv.NewMat()
  gocv.CvtColor(scan, &scanGray, gocv.ColorBGRToGray)


  gocv.IMWrite("./imgs-test/test.JPG", scanGray)
}
