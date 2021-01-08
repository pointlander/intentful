// Copyright 2021 The Intentful Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/pointlander/calc"
	"github.com/pointlander/dbnary"
)

// Index is the index page
const Index = `<html>
  <head><title>Intentful</title></head>
  <body>
    <h3>Dictionary</h3>
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
    <h3>Calculate</h3>
    <form action="/calculate" method="post">
      <input type="text" id="expression" name="expression">
      <input type="submit" value="Submit">
    </form>
  </body>
</html>
`

// Result is the result page
const Result = `<html>
  <head><title>Result</title></head>
  <body>
    <b>%s</b>
  </body>
</html>
`

// Interface outputs the search interface
func Interface(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(Index))
}

// Data is data for endpoints
type Data struct {
	DB *dbnary.DB
}

// Search redirects to the word search
func (d *Data) Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	lang, query := r.Form["language"][0], r.Form["query"][0]
	word, err := d.DB.LookupWordForLanguage(query, lang)
	if err != nil || (len(word.Relations) == 0 && len(word.Parts) == 0) {
		http.Redirect(w, r, fmt.Sprintf("/word-search/%s/%s", lang, query), http.StatusMovedPermanently)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/word/%s/%s", lang, query), http.StatusMovedPermanently)
}

// Calculate calculates an expression
func Calculate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	expression := r.Form["expression"][0]
	cal := &calc.Calculator{Buffer: expression}
	cal.Init()
	if err := cal.Parse(); err != nil {
		fmt.Println(err)
		return
	}
	result, html := cal.Eval(), ""
	if result.Matrix != nil {
		html = fmt.Sprintf(Result, result.Matrix.String())
	} else {
		html = fmt.Sprintf(Result, result.Expression.String())
	}
	w.Write([]byte(html))
}

// Address is the address and port of the server
var Address = flag.String("address", ":80", "the address and port of the server")

func main() {
	flag.Parse()

	db := dbnary.OpenDB("dbnary.db", true)
	defer db.Close()
	data := Data{
		DB: db,
	}
	router := httprouter.New()
	router.GET("/", Interface)
	router.POST("/search", data.Search)
	router.POST("/calculate", Calculate)
	dbnary.Server(db, router)
	server := http.Server{
		Addr:    *Address,
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
