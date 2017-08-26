package cognitron

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/lazywei/go-opencv/opencv"
)

const basePath string = "../train_base/"
const filtersPath string = "../filters/"

var faceDetectCascade *opencv.HaarCascade

// var faceDet2 *opencv.HaarCascade
// var faceDet3 *opencv.HaarCascade
// var faceDet4 *opencv.HaarCascade
var emotions []string

func init() {
	faceDetectCascade = opencv.LoadHaarClassifierCascade(filtersPath + "f_haarcascade_frontalface_alt.xml")
	// faceDet2 = opencv.LoadHaarClassifierCascade("haarcascade_frontalface_alt2.xml")
	// faceDet3 = opencv.LoadHaarClassifierCascade(filtersPath + "haarcascade_frontalface_alt.xml")
	// faceDet4 = opencv.LoadHaarClassifierCascade(filtersPath + "haarcascade_frontalface_alt_tree.xml")
	emotions = []string{
		"neutral", "anger", "disgust", "fear", "happiness", "sadness", "surprise",
	}
}

func extractFace(path string) (*opencv.IplImage, error) {
	fmt.Println(path)

	image := opencv.LoadImage(path)
	defer image.Release()

	w := image.Width()
	h := image.Height()

	var greyImage = opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
	defer greyImage.Release()

	opencv.CvtColor(image, greyImage, opencv.CV_BGR2GRAY)
	faces := faceDetectCascade.DetectObjects(greyImage)

	for _, value := range faces {
		croppedWidth := value.Width()
		croppedHeight := value.Height()
		croppedFace := opencv.Crop(greyImage, value.X(), value.Y(), croppedWidth, croppedHeight)
		defer croppedFace.Release()

		resizedFace := opencv.Resize(croppedFace, 350, 350, opencv.CV_INTER_CUBIC)
		return resizedFace, nil
	}

	return nil, errors.New("Face not detected")
}

// PreprocessDatabase gets faces from all images in the database
func PreprocessDatabase() {
	for _, emotionsFolder := range emotions {
		files, _ := ioutil.ReadDir(basePath + emotionsFolder)
		for _, f := range files {
			if strings.HasPrefix(f.Name(), "img_") == false {
				continue
			}
			faceFrame, err := extractFace(path.Join(basePath+emotionsFolder, f.Name()))
			if err != nil {
				log.Println("Face not detected")
				continue
			}

			defer faceFrame.Release()

			nameTail := strings.Split(f.Name(), "_")[1]
			newName := "face_" + nameTail

			fmt.Printf("%s saved\n", path.Join(basePath+emotionsFolder, newName))
			opencv.SaveImage(path.Join(basePath+emotionsFolder, newName), faceFrame, nil)
		}
	}

	// def detect_faces(emotion):
	//     files = glob.glob("sorted_set\\%s\\*" %emotion) #Get list of all images with emotion
	//
	//     filenumber = 0
	//     for f in files:
	//
	//
	//
	//         #Cut and save face
	//         for (x, y, w, h) in facefeatures: #get coordinates and size of rectangle containing face
	//             print "face found in file: %s" %f
	//             gray = gray[y:y+h, x:x+w] #Cut the frame to size
	//
	//             try:
	//                 out = cv2.resize(gray, (350, 350)) #Resize face so all images have same size
	//                 cv2.imwrite("dataset\\%s\\%s.jpg" %(emotion, filenumber), out) #Write image
	//             except:
	//                pass #If error, pass file
	//         filenumber += 1 #Increment image number
	//
	// for emotion in emotions:
	//     detect_faces(emotion) #Call functiona
}

// DrawFaceFrame draws a frame around the face on an image
func DrawFaceFrame(incomingImgURL string) {
	image := opencv.LoadImage(incomingImgURL)
	defer image.Release()

	faces := faceDetectCascade.DetectObjects(image)
	color := opencv.NewScalar(0.0, 255.0, 0.0, 255.0)
	for _, value := range faces {
		opencv.Rectangle(image,
			opencv.Point{X: value.X() + value.Width(), Y: value.Y()},
			opencv.Point{X: value.X(), Y: value.Y() + value.Height()},
			color, 2, 1, 0)
	}

	opencv.SaveImage(incomingImgURL, image, nil)
}
