package server

import (
	"bytes"
	"cfg"
	"cognitron"
	"db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"models/comment"
	"models/picture"
	"models/user"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jirfag/gointensive/lec3/2_http_pool/workers"
)

// IPool is a magic interface for WorkerPool
type IPool interface {
	Size() int
	Run()
	AddTaskSyncTimed(f workers.Func, timeout time.Duration) (interface{}, error)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var wp IPool
var requestWaitInQueueTimeout = 0 * time.Microsecond

func init() {
	rand.Seed(time.Now().UnixNano())

	numOfWorkers := cfg.GetApp().Cognitron.MaxJobs
	timeout := time.Duration(cfg.GetApp().Cognitron.Timeout)
	requestWaitInQueueTimeout = timeout * time.Millisecond

	wp = workers.NewPool(numOfWorkers)
	wp.Run()
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

func writeJSONResponse(w http.ResponseWriter, obj interface{}) {
	seq, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal to json")
		return
	}

	if _, err := w.Write(seq); err != nil {
		log.Printf("can't write to connection: %s", err)
	}
}

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
		fmt.Println("Error retrieving image: " + err.Error())
		return ""
	}
	defer file.Close()
	newFilename := "../stored_images/" + randStringBytes(42) + ".jpg"
	f, err := os.OpenFile(newFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error writing image: " + err.Error())
		return ""
	}
	defer f.Close()
	io.Copy(f, file)
	return newFilename
}

func cognitronHandler(w http.ResponseWriter, r *http.Request) {
	_, err := wp.AddTaskSyncTimed(func() interface{} {
		if r.Method != "POST" {
			return nil
		}
		// request read
		r.ParseMultipartForm(32 << 20)
		incomingFilename := dumpJpegImage(r)
		nickname := r.FormValue("user_nickname")
		log.Println("Image name = " + incomingFilename)
		log.Println("User nickname = " + nickname)
		log.Println("Image dumped")

		// recognition/processing
		cognitron.DrawFaceFrame(incomingFilename)
		log.Println("Face recognized")

		usr := user.Get(nickname)
		if usr != nil {
			// storing in database
			p := picture.Picture{
				UserID: db.Int64FK(int64(usr.ID)),
				URL:    incomingFilename,
			}
			p.Save()
			log.Println("Picture URL stored in DB")
		}

		if usr == nil {
			defer os.Remove(incomingFilename)
		}

		// response write
		writeJpegFile(w, incomingFilename)
		return nil
	}, requestWaitInQueueTimeout)

	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s!\n", err), 500)
	}
}

func writeOKMessage(w http.ResponseWriter, msg string) {
	obj := struct {
		Status  string
		Message string
	}{
		Status:  "OK",
		Message: msg,
	}
	writeJSONResponse(w, obj)
}

func writeErrorMessage(w http.ResponseWriter, msg string) {
	obj := struct {
		Status  string
		Message string
	}{
		Status:  "Error",
		Message: msg,
	}
	writeJSONResponse(w, obj)
}

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	imageID, err := strconv.ParseInt(r.FormValue("id"), 10, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s!\n", err), 500)
	}
	pict := picture.Get(int(imageID))
	if pict != nil {
		writeJpegFile(w, pict.URL)
		return
	}

	writeErrorMessage(w, "Picture with such ID not found")
}

func getImagesListHandler(w http.ResponseWriter, r *http.Request) {
	userNickname := r.FormValue("nickname")
	pictures := user.GetPictures(userNickname)
	if pictures != nil {
		log.Printf("Number of pictures = %d\n", len(pictures))
		ids := make([]uint, len(pictures))
		for i, pict := range pictures {
			ids[i] = pict.ID
		}
		writeJSONResponse(w, ids)
		return
	}
	writeErrorMessage(w, "User with such nickname not found")
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	userNickname := r.FormValue("nickname")
	userEmail := r.FormValue("email")
	oldusr := user.Get(userNickname)
	if oldusr == nil {
		newusr := user.User{
			Nickname: userNickname,
			Email:    userEmail,
		}
		newusr.Add()
		resp := struct {
			Status  string
			Message string
			UserID  uint
		}{
			Status:  "OK",
			Message: fmt.Sprintf("User '%s' created successfully", userNickname),
			UserID:  newusr.ID,
		}
		writeJSONResponse(w, resp)
		return
	}
	writeErrorMessage(w, "User with such nickname already exists")
}

func leaveCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	authorNickname := r.FormValue("author_nickname")
	pictureIDStr := r.FormValue("picture_id")

	pictureID, err := strconv.ParseInt(pictureIDStr, 10, 64)
	if err != nil {
		log.Println("Error retrieving picture_id: " + err.Error())
		writeErrorMessage(w, "Error parsing picture_id, comment not left")
		return
	}

	author := user.Get(authorNickname)
	if author == nil {
		log.Println("Error: author of comment not recognized")
		writeErrorMessage(w, "Author unknown, comment not left")
		return
	}

	pict := picture.Get(int(pictureID))
	if pict == nil {
		log.Println("Error: picture with such id not exists")
		writeErrorMessage(w, "No picture with such ID, comment not left")
		return
	}

	commentText := r.FormValue("comment_text")
	comment.Leave(int(author.ID), int(pictureID), commentText)
	writeOKMessage(w, "Comment left successfully")
}

// CommentExpress used to response as JSON
type CommentExpress struct {
	Author string
	Text   string
	Date   string
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	pictureIDStr := r.FormValue("picture_id")

	pictureID, err := strconv.ParseInt(pictureIDStr, 10, 64)
	if err != nil {
		log.Println("Error retrieving picture_id: " + err.Error())
		writeErrorMessage(w, "Error parsing picture_id, comment not left")
		return
	}

	pict := picture.Get(int(pictureID))
	if pict == nil {
		log.Println("Error: picture with such id not exists")
		writeErrorMessage(w, "No picture with such ID, comment not left")
		return
	}

	comments := pict.GetComments()
	resp := make([]CommentExpress, len(comments))
	for i, cmt := range comments {
		author := user.GetByID(int(cmt.AuthorID.Int64))
		resp[i] = CommentExpress{
			Author: author.Nickname, //fmt.Sprintf("%d", cmt.AuthorID.Int64),
			Text:   cmt.Text,
			Date:   string(cmt.CreatedAt.Format("2006-01-02 15:04:05")),
		}
	}
	writeJSONResponse(w, resp)
}

// RunServer launches HTTP server
func RunServer(listenAddr string) {
	log.Println("Starting HTTP server on " + listenAddr)
	//http.HandleFunc("/json/", myJSONHandler)
	// http.HandleFunc("/", handler)

	http.HandleFunc("/detect_face", cognitronHandler)
	http.HandleFunc("/image", getImageHandler)
	http.HandleFunc("/images_list", getImagesListHandler)
	http.HandleFunc("/sign_up", signUpHandler)
	http.HandleFunc("/leave_comment", leaveCommentHandler)
	http.HandleFunc("/get_comments", getCommentsHandler)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
