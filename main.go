package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type ListOfPages struct {
	Pages []Page
}
type Page struct {
	Name  string
	Path  string
	Title string
	Body  string
}

type Service struct {
	h func(http.ResponseWriter, *http.Request) error
}

func (p *Page) save() error {
	os.Mkdir("data", 0777)
	filename := "data/" + p.Name
	p.Body = p.Title + "\n" + p.Body + "\n"
	return ioutil.WriteFile(filename, []byte(p.Body), 0600)
}

func (h Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.h(w, r)
	if err != nil {
		log.Println("error:", err)
	}
}

func GetResponse(w http.ResponseWriter, r *http.Request) error {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		err := fmt.Errorf("Not Found")
		return err
	}
	some := strings.SplitAfter(parsedURL.Path, "/")[1]
	if some == "edit/" {
		editHandler(w, r, parsedURL.Path)
	}
	if some == "save/" {
		saveHandler(w, r, parsedURL.Path)
	} else if some == "" {
		rootHandler(w, r)
	}
	viewHandler(w, r, parsedURL.Path)
	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("data/")
	if err != nil {
		log.Fatal(err)
	}
	var p []Page
	for _, file := range files {
		page, _ := loadPage(file.Name())
		p = append(p, *page)
	}
	err = templates.ExecuteTemplate(w, "root.html", ListOfPages{Pages: p})
	if err != nil {
		panic(err)
	}
}

func loadPage(path string) (*Page, error) {
	filename := "data/" + path
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	title := strings.Trim(strings.SplitAfter(string(body), "\n")[0], "\n")
	pageBody := strings.SplitAfter(string(body), title+"\n")[1]
	filePath := strings.Trim(string(path), ".txt")
	return &Page{Name: path, Path: filePath, Title: title, Body: pageBody}, nil
}

func saveHandler(w http.ResponseWriter, r *http.Request, fullPath string) {
	body := r.FormValue("body")
	title := r.FormValue("title")
	name := strings.SplitAfter(fullPath, "/")[2]
	path := strings.Trim(string(name), ".txt")
	p := &Page{Name: name, Path: path, Title: title, Body: body}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/"+path, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, fullPath string) {
	path := strings.SplitAfter(fullPath, "/")[2]
	p, err := loadPage(path + ".txt")
	if err != nil {
		log.Println(err)
		p = &Page{Name: path}
	}
	log.Printf("p: :%s, :%s, :%s", p.Name, p.Body, p.Title)
	renderTemplate(w, "edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request, filename string) {
	p, err := loadPage(filename + ".txt")
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "view.html", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", Service{GetResponse})
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	return mux
}

var linkRegexp = regexp.MustCompile("\\[([a-zA-Z0-9]+)\\]")
var templates = template.Must(template.New("body.html").ParseFiles("tmpl/edit.html", "tmpl/view.html", "tmpl/root.html", "tmpl/header.html", "tmpl/body.html"))

func main() {
	numbPtr := flag.Int("port", 8080, "server port value")
	flag.Parse()
	if *numbPtr != 0 {
		var port = *numbPtr
		log.Printf("Starting server on :%d port", port)
		http.ListenAndServe(fmt.Sprintf(":%d", *numbPtr), GetMux())
	} else {
		log.Println("Something went comepletely wrong, sorry :(")
	}
}
