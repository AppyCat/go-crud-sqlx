package main

import (
	"assets"
	"fmt"
	"github.com/go-zoo/bone"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"log"
	"net/http"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=laptop dbname=estelle_test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	r := render.New(render.Options{
		Directory:  "views",
		Extensions: []string{".html"},
	})

	mux := bone.New()
	ServeResource := assets.ServeResource
	mux.HandleFunc("/img/", ServeResource)
	mux.HandleFunc("/css/", ServeResource)
	mux.HandleFunc("/js/", ServeResource)

	mux.HandleFunc("/foofoo", func(w http.ResponseWriter, req *http.Request) {
		rows, err := db.Queryx("SELECT id, title FROM pages")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		type Pages struct {
			Id    int
			Title string
		}

		Data := Pages{}
		for rows.Next() {
			err := rows.StructScan(&Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%#v\n", Data)
		}
		r.HTML(w, http.StatusOK, "foofoo", Data)
	})

	mux.HandleFunc("/bar", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "bar", nil)
	})

	mux.HandleFunc("/home/:id", func(w http.ResponseWriter, req *http.Request) {
		id := bone.GetValue(req, "id")
		r.HTML(w, http.StatusOK, "index", id)
	})

	mux.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "foo", nil)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "index", nil)
	})

	http.ListenAndServe(":8080", mux)
}
