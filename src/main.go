package main

import "server"

func main() {
	server.RunServer(":12177")
	// var listenAddr string
	// flag.StringVar(&listenAddr, "listen-addr", ":12177",
	// 	"address to listen")
	// flag.Parse()
	// log.Printf("listening on %q", listenAddr)

	// var usr user.User = user.User{
	// 	Nickname: "JSmith",
	// 	Email:    "johns@foo.com",
	// }
	//
	// err := usr.AddUser()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//image_path := "../train_base/happiness/img_0037.jpg"
	//cognitron.DrawFaceFrame(image_path)
	//cognitron.PreprocessDatabase()
}
