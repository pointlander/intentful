// Copyright 2021 The Intentful Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/pointlander/dbnary"
)

// Index is the index page
const Index = `<html>
  <head><title>Intentful</title></head>
  <body>
    <form action="/search" method="post">
      <select id="language" name="language">
        <option value="bul">Bulgarian</option>
        <option value="ell">Greek</option>
        <option value="eng" selected>English</option>
        <option value="spa">Spanish</option>
        <option value="fin">Finnish</option>
        <option value="fra">French</option>
        <option value="ind">Indonesian</option>
        <option value="ita">Italian</option>
        <option value="jpn">Japanese</option>
        <option value="lat">Latin</option>
        <option value="lit">Lithuanian</option>
        <option value="mlg">Malagasy</option>
        <option value="nld">Dutch</option>
        <option value="nor">Norwegian</option>
        <option value="pol">Polish</option>
        <option value="por">Portuguese</option>
        <option value="rus">Russian</option>
        <option value="hbs">Serbo-Croatian</option>
        <option value="swe">Swedish</option>
        <option value="tur">Turkish</option>
      </select>
      <input type="text" id="query" name="query">
      <input type="submit" value="Submit">
    </form>
  </body>
</html>
`

// Interface outputs the search interface
func Interface(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(Index))
}

// Search redirects to the word search
func Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	lang, word := r.Form["language"][0], r.Form["query"][0]
	http.Redirect(w, r, fmt.Sprintf("/word-search/%s/%s", lang, word), http.StatusMovedPermanently)
}

func main() {
	db := dbnary.OpenDB("dbnary.db", true)
	defer db.Close()

	router := httprouter.New()
	router.GET("/", Interface)
	router.POST("/search", Search)
	dbnary.Server(db, router)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
