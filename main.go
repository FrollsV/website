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

	"github.com/frollsv/website/pages"
)

type ListOfPages struct {
	Pages []pages.Page
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
	var p []pages.Page
	for _, file := range files {
		page, err := loadPage(file.Name())
		if err != nil {
			panic(err)
		}
		p = append(p, page)
	}
	log.Println("P", p)
	pageTemplate, err := template.New("root").ParseFiles("tmpl/root.html", "tmpl/header.html", "tmpl/body.html", "tmpl/page.html", "tmpl/paragraph.html")
	if err != nil {
		panic(err)
	}

	err = pageTemplate.Execute(w, &ListOfPages{p})

	if err != nil {
		panic(err)
	}
}

func loadPage(path string) (pages.Page, error) {
	filename := "data/" + path
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return pages.Page{}, err
	}
	return pages.LoadArticle(string(body))
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

	//log.Printf("p: :%s, :%s, :%s", p.Name, p.Body, p.Title)
	//renderTemplate(w, "edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request, filename string) {
	log.Print("Loading ", filename)

	filename = "data/" + filename + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	p, err := pages.LoadArticle(string(body))
	if err != nil {
		log.Print(err)
		return
	}
	pageTemplate, err := template.New("root").ParseFiles("tmpl/root.html", "tmpl/header.html", "tmpl/view.html", "tmpl/page.html", "tmpl/paragraph.html")
	if err != nil {
		log.Print(err)
		return
	}
	pageTemplate.Execute(w, p)
	//err = templates.ExecuteTemplate(w, "page.html", p)

	if err != nil {
		log.Print(err)
		return
	}

	/*
		p, err := loadPage(filename + ".txt")
		if err != nil {
			return
		}
		templates.ExecuteTemplate(w, "view.html", p)
	*/
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

var templates = template.Must(template.New("tmpl/page.html").ParseFiles("tmpl/edit.html", "tmpl/page.html", "tmpl/paragraph.html", "tmpl/view.html", "tmpl/root.html", "tmpl/header.html", "tmpl/body.html"))

func main() {
	numbPtr := flag.Int("port", 8080, "server port value")
	flag.Parse()
	if *numbPtr <= 0 {
		log.Printf("you smart ass... negative TCP port?")
		return
	}
	if *numbPtr > 0 {
		var port = *numbPtr
		log.Printf("Starting server on :%d port", port)
		http.ListenAndServe(fmt.Sprintf(":%d", *numbPtr), GetMux())
	} else {
		log.Println("Something went comepletely wrong, sorry :(")
	}
}
