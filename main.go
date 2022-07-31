package main

import "scan-n-crop/cropper"

func main() {
	cropper := cropper.Cropper{}
	fileName := "./test-imgs/corner.JPG"
	cropper.Crop(fileName)
}
