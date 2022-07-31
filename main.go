package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"gocv.io/x/gocv"
)

type Photos []gocv.PointVector

func main() {
	Run()
}

func Run() {
	middle := "./imgs-test/middle.JPG"
	mid_rot := "./imgs-test/middle_rotated.JPG"
	crooked := "./imgs-test/crooked.JPG"
	corner := "./imgs-test/corner.JPG"

	img := LoadImg(middle)
	output := CropImg(&img)
	SaveImg("./cropped/middle_cropped.JPG", &output)

	img = LoadImg(mid_rot)
	output = CropImg(&img)
	SaveImg("./cropped/middle_rotated_cropped.JPG", &output)

	img = LoadImg(crooked)
	output = CropImg(&img)
	SaveImg("./cropped/crooked_cropped.JPG", &output)

	img = LoadImg(corner)
	output = CropImg(&img)
	SaveImg("./cropped/corner_cropped.JPG", &output)

	img.Close()
	output.Close()
}

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

func CropImg(img *gocv.Mat) gocv.Mat {
	output := gocv.NewMat()
	straight := gocv.NewMat()

	contours := GetContours(img)
	photoIdx := FindPhoto(&contours)
	rect := GetPhotoRect(&contours, photoIdx[0])

	contours.C.Close()
	contours.H.Close()

	straight = StraightenImg(img, &rect)

	szOffset := 30
	sz := image.Pt(rect.Height-szOffset, rect.Width-szOffset)
	fmt.Println("sz: ", sz)
	gocv.GetRectSubPix(straight, sz, rect.Center, &output)

	straight.Close()

	return output
}

func StraightenImg(img *gocv.Mat, rect *gocv.RotatedRect) gocv.Mat {
	output := gocv.NewMat()
	center := rect.Center
	var angle float64
	H := img.Rows()
	W := img.Cols()

	if W > H {
		angle = rect.Angle - 90
	} else {
		angle = (90 - rect.Angle) * -1
	}

	fmt.Printf("H: %v, W: %v, C: %v, A: %v\n", H, W, center, angle)
	sz := image.Pt(W, H)

	M := gocv.GetRotationMatrix2D(center, angle, 1)

	gocv.WarpAffine(*img, &output, M, sz)

	return output
}

type Contours struct {
	C gocv.PointsVector
	H gocv.Mat
}

func GetContours(img *gocv.Mat) Contours {
	C := Contours{}
	C.C = gocv.NewPointsVector()
	C.H = gocv.NewMat()

	imgGray := gocv.NewMat()
	thresh := gocv.NewMat()

	gocv.CvtColor(*img, &imgGray, gocv.ColorBGRToGray)
	// _ = gocv.Threshold(imgGray, &thresh, 127, 255, gocv.ThresholdBinary)
	_ = gocv.Threshold(imgGray, &thresh, 245, 255, gocv.ThresholdToZeroInv)

	C.C = gocv.FindContoursWithParams(
		thresh,
		&C.H,
		gocv.RetrievalTree,
		gocv.ChainApproxSimple,
	)

	imgGray.Close()
	thresh.Close()

	return C
}

func FindPhoto(C *Contours) []int {
	photoIdx := []int{}

	for i := 0; i < C.H.Cols(); i++ {
		v := C.H.GetVeciAt(0, i)
		_, _, c, p := v[0], v[1], v[2], v[3]
		contourArea := gocv.ContourArea(C.C.At(i))
		// if p == 0 && c != -1 && contourArea > 2500 {
		if p == -1 && c != -1 && contourArea > 200 {
			photoIdx = append(photoIdx, i)
		}
	}
	fmt.Printf("len(photoIdx): %v\n", len(photoIdx))
	return photoIdx
}

func GetPhotoRect(S *Contours, Idx int) gocv.RotatedRect {
	rect := gocv.MinAreaRect(S.C.At(Idx))
	return rect
}

func DrawBoundingBox(rect gocv.RotatedRect, img *gocv.Mat) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	pts := rect.Points

	gocv.Line(img, pts[0], pts[1], red, 10)
	gocv.Line(img, pts[1], pts[2], red, 10)
	gocv.Line(img, pts[2], pts[3], red, 10)
	gocv.Line(img, pts[3], pts[0], red, 10)
}
