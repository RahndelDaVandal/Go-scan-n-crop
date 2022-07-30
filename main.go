package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"sort"

	"gocv.io/x/gocv"
)

type Photos []gocv.PointVector

func main() {
	// CornerDebug()
	t := []int{1, 3, 2}
	m := Median(t)
	fmt.Printf("%v: %v\n", m, t[int(m)])
}

type Contour struct {
	i int
	n int32
	l int32
	c int32
	p int32
	a float64
}
type Container []Contour

func (pC *Container) Append(v gocv.Veci, a float64, i int) {
	c := Contour{
		i: i,
		n: v[0],
		l: v[1],
		c: v[2],
		p: v[3],
	}
	*pC = append(*pC, c)
}

func (C *Contour) Show() {
	fmt.Printf("%v [ n: %v, l: %v, c: %v, p: %v ] a: %v\n", C.i, C.n, C.l, C.c, C.p, C.a)
}

func CornerDebug() {
	img := LoadImg("./imgs-test/corner.JPG")
	C := GetContours(&img)
	fmt.Printf("C: %v, H: %v\n", C.C.Size(), C.H.Cols())

	container := Container{}

	for i := 0; i < C.C.Size(); i++ {
		v := C.H.GetVeciAt(0, i)
		A := gocv.ContourArea(C.C.At(i))
		if A > 1000 {
			fmt.Println(A)
		}
		container.Append(v, A, i)
	}

	for i := 0; i < len(container); i++ {
		tmp := container[i]
		if tmp.a > 0 {
			tmp.Show()
		}
	}
}

func Median(values []int) int {
	valuesMap := make(map[int]int)
	for i, v := range values {
		valuesMap[v] = i
	}
	fmt.Println("map: ", valuesMap)

	valuesCopy := make([]int, len(values))
	copy(valuesCopy, values)
	fmt.Println("valuesCopy: ", valuesCopy)
	sort.Ints(valuesCopy)
	fmt.Println("Sorted valuesCopy: ", valuesCopy)

	l := len(valuesCopy)
	fmt.Println("l: ", l)
	var m int

	if l == 0 {
		log.Fatal("len(valuesCopy) == 0")
	} else if l%2 == 0 {
		m = (l / 2) - (l % 2)
		fmt.Println("m: ", m)
	} else {
		m = l / 2
		fmt.Println("m: ", m)
	}

	return valuesMap[valuesCopy[m]]
}

func QuickSort(c Container) Container {
	if len(c) <= 1 {
		return c
	} else {
		mid := int(math.Floor(float64(len(c) / 2)))
		m := map[float64]int{
			c[0].a:      0,
			c[mid].a:    1,
			c[len(c)].a: 2,
		}
		var tmp []float64
		for k, _ := range m {
			tmp = append(tmp, k)
		}
		pivot := c[m[tmp[1]]]
	}
}

func Run() {
	middle := "./imgs-test/middle.JPG"
	mid_rot := "./imgs-test/middle_rotated.JPG"
	crooked := "./imgs-test/crooked.JPG"

	img := LoadImg(middle)
	output := CropImg(&img)
	SaveImg("./cropped/middle_cropped.JPG", &output)

	img = LoadImg(mid_rot)
	output = CropImg(&img)
	SaveImg("./cropped/middle_rotated_cropped.JPG", &output)

	img = LoadImg(crooked)
	output = CropImg(&img)
	SaveImg("./cropped/crooked_cropped.JPG", &output)

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

	sz := image.Pt(rect.Height, rect.Width)
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
	_ = gocv.Threshold(imgGray, &thresh, 127, 255, gocv.ThresholdBinary)

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
		if p == 0 && c != -1 && contourArea > 2500 {
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
