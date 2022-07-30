package Scan

import (
	"log"
	"os"

	"gocv.io/x/gocv"
)

type Scan struct {
	BGR       gocv.Mat
	Gray      gocv.Mat
	Thres     gocv.Mat
	Output    gocv.Mat
	Hierarchy gocv.Mat
	Contours  gocv.PointsVector
}

func NewScan(fileName string) Scan {
	s := Scan{
		BGR:       gocv.NewMat(),
		Gray:      gocv.NewMat(),
		Thres:     gocv.NewMat(),
		Output:    gocv.NewMat(),
		Hierarchy: gocv.NewMat(),
		Contours:  gocv.NewPointsVector(),
	}

	_, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
	}
	img := gocv.IMRead(fileName, gocv.IMReadColor)
	s.BGR = img
	imgGray := gocv.NewMat()
	gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)
	s.Gray = imgGray
	_ = gocv.Threshold(s.Gray, &s.Thres, 127, 255, gocv.ThresholdBinary)

	return s
}

func (s *Scan) SaveImg() {
	// gocv.IMWrite("output.JPG", s.Output)
	gocv.IMWrite("output.JPG", s.Output)
}
