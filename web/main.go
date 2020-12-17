package main

import (
	"html/template"
	"fmt"
	"net/http"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/gorilla/mux"
)

// Post é apenas um post
type Post struct {
	ID int
	Title string
	Body string
}

var db, err = sql.Open("mysql", "root:root@/go_course?charset=utf8")

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)

	fmt.Println(http.ListenAndServe(":8080", router))
}

func checkErr (err error) {
	if err != nil {
		panic(err)
	}
}

func listPosts() []Post {
	rows, err := db.Query("select * from posts")
	checkErr(err)

	items := []Post{}

	for rows.Next() {
		post := Post{}

		rows.Scan(&post.ID, &post.Title, &post.Body)
		items = append(items, post)
	}

	return items
}

// HomeHandler handles a route and template
func HomeHandler (w http.ResponseWriter, r *http.Request) {
	t :=  template.Must(template.ParseFiles("templates/index.html"))
	if err := t.ExecuteTemplate(w, "index.html", listPosts()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}