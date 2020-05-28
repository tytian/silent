package main

import (
	"log"
	"net/http"
)

// Run 路由
func Run() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", IndexView)
	http.HandleFunc("/upload", UploadView)
	http.HandleFunc("/uploadMore", UploadMoreView)
	http.HandleFunc("/list", ListView)
	http.HandleFunc("/detail", DetailView)

	http.HandleFunc("/api/upload", ApiUpload)
	http.HandleFunc("/api/uploadMore", ApiUploadMore)
	http.HandleFunc("/api/list", ApiList)
	http.HandleFunc("/api/drop", ApiDrop)

	log.Println("Run 8080 ...")
	_ = http.ListenAndServe(":8080", nil)
}
