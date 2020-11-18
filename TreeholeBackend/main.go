package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/newpost", NewPost)
	http.HandleFunc("/newcomment", NewComment)
	http.HandleFunc("/upgood", UpGood)
	http.HandleFunc("/downbad", DownBad)
	http.HandleFunc("/getpost", GetPost)


	err := http.ListenAndServeTLS(":80","1_tonggege.work_bundle.crt", "2_tonggege.work.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
