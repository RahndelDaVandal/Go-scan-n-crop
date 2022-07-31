package cropper

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gocv.io/x/gocv"
)

type Cropper struct {
	FileName string
}

func (ptr *Cropper) Crop(fileName string, dstPath string) {
	// TODO - Add Arg for cropped size offset
	// TODO - Add Arg for dst dir
	ptr.FileName = fileName

	// Load Scan Image
  img := loadImg(fileName)

	// Generate Grayscale Image from Loaded Image
	imgGray := gocv.NewMat()
	gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)

	// Generate Threshold Image from Grayscale Image
	imgThresh := gocv.NewMat()
	_ = gocv.Threshold(imgGray, &imgThresh, 245, 255, gocv.ThresholdToZeroInv)
	imgGray.Close()

	// Get Contours from Threshold Image
	contours := getContours(&imgThresh)
	imgThresh.Close()

	// Find Photo Contours
	photoIdx := findPhoto(&contours)

	// Created a Rotated Rectangle
	rect := getPhotoRect(&contours, photoIdx)

	// Create a 2d Rotational Matix from Rotated Rectangle
	rotMat := getRotationalMatrix(&img, &rect)

	// Rotate Original Scanned Image with Rotational Matrix
	rotImg := rotateImage(&img, &rotMat)

	// Crop Scan with Rotated Rectangle
	croppedImg := gocv.NewMat()
	szOffset := 0
	// You have to flip Height & Width or it won't crop right
	sz := image.Pt(rect.Height-szOffset, rect.Width-szOffset)
	gocv.GetRectSubPix(rotImg, sz, rect.Center, &croppedImg)

	// Save Cropped Image to File
  outputName := outputFileName(fileName, dstPath)
	saveImg(outputName, &croppedImg)
  fmt.Println("Cropped photo saved to", outputName)

	// Clean Up
	contours.vectors.Close()
	contours.hierarchy.Close()
	rotImg.Close()
	croppedImg.Close()
}

func outputFileName(fileName string, dstPath string) string{
  origFile := filepath.Base(fileName)
  fileExt := filepath.Ext(origFile)
  originalName := strings.Trim(origFile, fileExt)
  outputName := fmt.Sprintf("%v_CROPPED%v", originalName, fileExt)
  output := filepath.Join(dstPath, outputName)
  return output
}

func loadImg(fileName string) gocv.Mat {
	_, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
	}
	img := gocv.IMRead(fileName, gocv.IMReadColor)
	return img
}

func saveImg(fileName string, img *gocv.Mat) {
	gocv.IMWrite(fileName, *img)
}

type Contours struct {
	vectors   gocv.PointsVector
	hierarchy gocv.Mat
}

func getContours(img *gocv.Mat) Contours {
	hierarchy := gocv.NewMat()

	vectors := gocv.FindContoursWithParams(
		*img,
		&hierarchy,
		gocv.RetrievalTree,
		gocv.ChainApproxSimple,
	)

	contours := Contours{
		vectors:   vectors,
		hierarchy: hierarchy,
	}

	return contours
}

func findPhoto(c *Contours) int {
	photos := []int{}

	for i := 0; i < c.hierarchy.Cols(); i++ {
		v := c.hierarchy.GetVeciAt(0, i)
		// Hierarchy contains relation of each contour to one another
		// [next, previous, children, parent]
		_, _, child, parent := v[0], v[1], v[2], v[3]
		contourArea := gocv.ContourArea(c.vectors.At(i))
		if parent == -1 && child != -1 && contourArea > 200 {
			photos = append(photos, i)
		}
	}

	if len(photos) > 1 {
		log.Fatal("Error: Found more than one photo in scan")
	}
	return photos[0]
}

func getPhotoRect(c *Contours, Idx int) gocv.RotatedRect {
	rect := gocv.MinAreaRect(c.vectors.At(Idx))
	return rect
}

func getRotationalMatrix(img *gocv.Mat, rect *gocv.RotatedRect) gocv.Mat {
	center := rect.Center
	var angle float64
	H := img.Rows()
	W := img.Cols()

	if W > H {
		angle = rect.Angle - 90
	} else {
		angle = (90 - rect.Angle) * -1
	}
	// sz := image.Pt(W, H)
	M := gocv.GetRotationMatrix2D(center, angle, 1)

	return M
}

func rotateImage(img *gocv.Mat, rotMat *gocv.Mat) gocv.Mat {
	rotImg := gocv.NewMat()

	H := img.Rows()
	W := img.Cols()

	size := image.Pt(W, H)

	gocv.WarpAffine(*img, &rotImg, *rotMat, size)

	return rotImg
}
