package main

import (
	"flag"
	"log"
	"models/user"
)

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", ":12177",
		"address to listen")
	flag.Parse()
	log.Printf("listening on %q", listenAddr)

	var usr user.User = user.User{
		Nickname: "JSmith",
		Email:    "johns@foo.com",
	}

	err := usr.AddUser()
	if err != nil {
		log.Fatal(err)
	}
	//db.Get().Save(&usr)
	// usr.GetOrCreate()
	// fmt.Printf("%d\n", usr.UserID)
}
