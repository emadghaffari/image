package main

import (
	"fmt"
	"image"
	gif "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"

	"github.com/lazywei/go-opencv/opencv"
)

func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s <imagefilename>\n", os.Args[0])
		os.Exit(0)
	}

	imageFileName := os.Args[1]

	// we will use Go's method of load images
	// instead of openCV.LoadImage
	// because we want to detect if the user supplies animated GIF or not
	imageFile, err := os.Open(imageFileName)

	errCheck(err)

	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	errCheck(err)

	buffer := make([]byte, 512)
	imageFile.Seek(0, 0) // reset reader
	_, err = imageFile.Read(buffer)
	errCheck(err)

	filetype := http.DetectContentType(buffer)
	// check if image is GIF and if yes, check to see if it is animated GIF by
	// counting the LoopCount number
	fmt.Println("Analyzing image type : ", filetype)

	if filetype == "image/gif" {
		imageFile.Seek(0, 0)
		// warn if image is animated GIF
		gif, err := gif.DecodeAll(imageFile)
		errCheck(err)
		if gif.LoopCount != 0 {
			fmt.Println("Animated gif detected. Will only scan for faces in the 1st frame.")
		}
	}

	// convert Go's image.Image type to OpenCV's IplImage(Intel Image Processing Library)
	openCVImg := opencv.FromImage(img)
	defer openCVImg.Release()

	if openCVImg != nil {
		fmt.Println("Converting [" + imageFileName + "] to greyscale image......")

		w := openCVImg.Width()
		h := openCVImg.Height()

		// create an IplImage with 1 channel(grey)
		greyImg := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
		defer greyImg.Release()

		// convert to greyscale
		opencv.CvtColor(openCVImg, greyImg, opencv.CV_BGR2GRAY)

		// from http://www.shervinemami.info/faceRecognition.html

		// make the image a fixed size
		// CV_INTER_CUBIC or CV_INTER_LINEAR is good for enlarging, and
		// CV_INTER_AREA is good for shrinking / decimation, but bad at enlarging.

		histoEqImg := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
		defer histoEqImg.Release()
		histoEqImg = opencv.Resize(greyImg, w, h, opencv.CV_INTER_LINEAR)

		fmt.Println("Applying Histogram Equalization to [" + imageFileName + "]......")

		// standard brightness and contrast
		greyImg.EqualizeHist(histoEqImg)

		// save to file
		fmt.Println("Saving results to [grey.jpg] and [equalhisto.jpg]")
		opencv.SaveImage("./grey.jpg", greyImg, opencv.CV_IMWRITE_JPEG_QUALITY)
		opencv.SaveImage("./equalhisto.jpg", histoEqImg, opencv.CV_IMWRITE_JPEG_QUALITY)

	} else {
		panic("OpenCV FromImage error")
	}

}
