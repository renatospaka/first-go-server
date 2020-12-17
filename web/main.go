package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Post é apenas um post
type Post struct {
	ID    int
	Title string
	Body  string
}

var db, err = sql.Open("mysql", "root:root@/go_course?charset=utf8")

func main() {

	router := mux.NewRouter()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/{id}/view", ViewHandler)

	fmt.Println(http.ListenAndServe(":8080", router))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//GetPosts retrieves all posts
func GetPosts() []Post {
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

//GetPostById retrieves a specific post based on its id
func GetPostById(id string) Post {
	post := Post{}
	row := db.QueryRow("select * from posts where id = ?", id)
	row.Scan(&post.ID, &post.Title, &post.Body)
	return post
}

// HomeHandler handles a route and template
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	if err := t.ExecuteTemplate(w, "index.html", GetPosts()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ViewHandler shows one post
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	t := template.Must(template.ParseFiles("templates/view.html"))
	if err := t.ExecuteTemplate(w, "view.html", GetPostById(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
