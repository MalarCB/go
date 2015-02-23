// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"html/template"
	"path"
	"database/sql"
	"log"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/russross/blackfriday"
	"os"
	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"

	_ "github.com/mattn/go-sqlite3"
	)

type Page struct {
	Title string
	Body  []byte
}

type Post struct {
    Title  string
		Link string
    Author string
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
	port := os.Getenv("PORT")
    if port == "" {
        port = "3030"
    }
	//http.HandleFunc("/", ShowPosts(db))
	//http.ListenAndServe(":8080", ShowPosts(db))

		r := mux.NewRouter().StrictSlash(false)
  	r.HandleFunc("/", ShowPosts)
	// Posts collection
	    posts := r.Path("/posts").Subrouter()
	    posts.Methods("GET").HandlerFunc(PostsIndexHandler)
	    posts.Methods("POST").HandlerFunc(PostsCreateHandler)

	    // Posts singular
	    post := r.PathPrefix("/posts/{id}").Subrouter()
	    post.Methods("GET").Path("/edit").HandlerFunc(PostEditHandler)
	    post.Methods("GET").HandlerFunc(PostShowHandler)
	    post.Methods("PUT", "POST").HandlerFunc(PostUpdateHandler)
	    post.Methods("DELETE").HandlerFunc(PostDeleteHandler)

 		r.HandleFunc("/markdown", GenerateMarkdown)
	    r.Handle("/", http.FileServer(http.Dir("public")))
		http.ListenAndServe(":"+port, r)

	//middleware stack
	n:= negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(MyMiddleware),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
		)

	n.Run(":"+port)
}


func NewDB() *sql.DB {
    db, err := sql.Open("sqlite3", "example.sqlite")
    if err != nil {
        panic(err)
    }

    _, err = db.Exec("create table if not exists posts(title text, link text, author text)")
    if err != nil {
        panic(err)
    }

    return db
}
func ShowPosts(rw http.ResponseWriter, r *http.Request ) {

			db := NewDB()
        var title, link, author string
        rows, err := db.Query("select title, link, author from posts")//.Scan(&title, &link, &author)
        if err != nil {
            panic(err)
        }



				fp := path.Join("public", "posts.html")
				tmpl, err := template.ParseFiles(fp)

				for rows.Next() {
					err := rows.Scan(&title, &link, &author)
					if err != nil {
						log.Fatal(err)
					}
					post := Post{title, link, author}

					if err != nil {
							http.Error(rw, err.Error(), http.StatusInternalServerError)
							return
					}

					if err := tmpl.Execute(rw, post); err != nil {
							http.Error(rw, err.Error(), http.StatusInternalServerError)
					}
				}

				fp1 := path.Join("public", "create.html")
				tmpl1, err1 := template.ParseFiles(fp1)
				if err := tmpl1.Execute(rw, ""); err != nil {
						http.Error(rw, err1.Error(), http.StatusInternalServerError)
				}
}


func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Logging on the way there...")

	if r.URL.Query().Get("password") == "secret123" {
			next(rw, r)
	} else {
			http.Error(rw, "Not Authorized", 401)
	}

	log.Println("Logging on the way back...")
}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "Home")
}

func PostsIndexHandler(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "posts index")
}

func PostsCreateHandler(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "posts create")
		link := blackfriday.MarkdownCommon([]byte(r.FormValue("link")))
    rw.Write(link)
}

func PostShowHandler(rw http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    fmt.Fprintln(rw, "showing post", id)
}

func PostUpdateHandler(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "post update")
}

func PostDeleteHandler(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "post delete")
}

func PostEditHandler(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "post edit")
}

func GenerateMarkdown(rw http.ResponseWriter, r *http.Request) {
    markdown := blackfriday.MarkdownCommon([]byte(r.FormValue("body")))
    rw.Write(markdown)
}
