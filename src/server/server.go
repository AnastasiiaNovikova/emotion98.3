package server

import (
	"bytes"
	"cognitron"
	"db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"models/picture"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	if err != nil {
		log.Printf("can't print to connection: %s", err)
	}
}

func myJSONHandler(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Code    string
		Message string
	}{
		Code:    "OK",
		Message: fmt.Sprintf("Hi there, I love %s!", strings.TrimPrefix(r.URL.Path, "/json/")),
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("can't json marshal %+v: %s", resp, err)
		return
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Printf("can't write to connection: %s", err)
	}
}

func writeJpegFile(w http.ResponseWriter, path string) {
	buffer := new(bytes.Buffer)
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln("Image file not exists")
	}
	io.Copy(buffer, f)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func dumpJpegImage(r *http.Request) string {
	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()
	newFilename := "../stored_images/" + randStringBytes(42) + ".jpg"
	f, err := os.OpenFile(newFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer f.Close()
	io.Copy(f, file)
	return newFilename
}

func cognitronHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	// request read
	r.ParseMultipartForm(32 << 20)
	incomingFilename := dumpJpegImage(r)
	strUserID := r.FormValue("user_id")
	log.Println("UserID = " + strUserID)
	userID, err := strconv.ParseInt(strUserID, 10, 64)
	if err != nil {
		log.Println("Error retrieving user_id: " + err.Error())
		userID = 100501
	}
	log.Println("Image dumped")

	// recognition/processing
	cognitron.DrawFaceFrame(incomingFilename)
	log.Println("Face recognized")

	// storing in database
	p := picture.Picture{
		UserID: db.Int64FK(int64(userID)),
		URL:    incomingFilename,
	}
	p.Save()
	log.Println("Picture URL stored in DB")

	// response write
	writeJpegFile(w, incomingFilename)
}

// RunServer launches HTTP server
func RunServer(listenAddr string) {
	log.Println("Starting HTTP server on " + listenAddr)
	// http.HandleFunc("/json/", myJSONHandler)
	// http.HandleFunc("/", handler)

	http.HandleFunc("/detect_face", cognitronHandler)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}