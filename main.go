package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"scan-n-crop/cropper"
	"strings"
)

func main() {
  var dstPath string
  path := os.Args[1]
  if len(os.Args) > 2{
    dstPath = os.Args[2]
  } else {
    dstPath = getDstPath(path)
  }

  dir, err := os.Stat(dstPath)
  if err != nil {
    log.Fatal(err)
  }
  if !dir.IsDir(){
    fmt.Printf("Ooops! The destination Directory provided ")
    fmt.Printf("doesn't appear to be a Directory!\n")
  }

  file, err := os.Stat(path)
  if err != nil {
    log.Fatal(err)
  }
  if file.IsDir(){
    BatchCropImgs()
  } else if ValidFileType(path){
    cropper := cropper.Cropper{}
    cropper.Crop(path, dstPath)
  } else {
    fmt.Printf("Ooops! That doesn't appear ")
    fmt.Printf("to be a valid filetype!\n")
  }

}

func ValidFileType(path string) bool {
  // TODO - Theres probably a better way....
  validFileTypes := []string{
    ".JPG",
    ".png",
    ".jpeg",
  }
  fileExtension := filepath.Ext(path)

  isValid := false
  for _, v := range validFileTypes{
    if fileExtension == v{
      isValid = true
    }
  }
  return isValid
}

func getDstPath(path string) string{
  fileName := filepath.Base(path)
  return strings.Trim(path, fileName)
}

func BatchCropImgs(){
  fmt.Printf("Ooops! The path you provided ")
  fmt.Printf("appears to be a Directory! Batch ")
  fmt.Printf("cropping for a Directory is coming ")
  fmt.Printf("soon so stay tuned!\n")
}
