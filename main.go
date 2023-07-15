package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/posts", handleGetList)
	http.HandleFunc("/posts/", handleRequest)

	server.ListenAndServe()
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// postの一覧を取得するハンドラ
func handleGetList(w http.ResponseWriter, r *http.Request) {
	var err error

	posts, err := getPosts(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	output, err := json.MarshalIndent(&posts, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write(output)
}

// GETで指定したidのpostを取得するハンドラ
func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	// Base は，pathの最後の要素を返します．末尾のスラッシュは，最後の要素を抽出する前に削除されます．
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	output, err := json.MarshalIndent(&post, "", "\t")
	if err != nil {
		return
	}

	w.WriteHeader(200)
	w.Write(output)
	return
}

// POSTでpostを新規作成するハンドラ
func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	contentLength := r.ContentLength
	contentBody := make([]byte, contentLength)
	r.Body.Read(contentBody)

	// ここから：以下の追加
	fmt.Printf("%T\n", r.Body)
	// ここまで

	var post Post
	err = json.Unmarshal(contentBody, &post)
	if err != nil {
		return
	}

	err = post.create()
	if err != nil {
		return
	}

	output, err := json.MarshalIndent(&post, "", "\t")
	if err != nil {
		return
	}

	w.WriteHeader(200)
	w.Write(output)
	return
}

// PUTで指定したidのpostを更新するハンドラ
func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	contentLength := r.ContentLength
	contentBody := make([]byte, contentLength)
	r.Body.Read(contentBody)

	// ここから：以下の追加
	fmt.Printf("%T\n", r.Body)
	// ここまで

	err = json.Unmarshal(contentBody, &post)
	if err != nil {
		return
	}

	err = post.update()
	if err != nil {
		return
	}

	output, err := json.MarshalIndent(&post, "", "\t")
	if err != nil {
		return
	}

	w.WriteHeader(200)
	w.Write(output)
	return
}

// DELETEで指定したidのpostを削除するハンドラ

func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	post, err := retrieve(id)
	if err != nil {
		return
	}
	err = post.delete()
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&post, "", "\t")
	if err != nil {
		return
	}
	w.WriteHeader(200)
	w.Write(output)
	return
}
