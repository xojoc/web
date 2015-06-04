// Private package.
package web

import (
	htpl "html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

var Pages *htpl.Template

func CreateLog(file string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Print(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))
}

type NetError struct {
	Code    int
	Message string
}

type myHandler func(http.ResponseWriter, *http.Request) *NetError

func ErrorHandler(h myHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nerr := h(w, r)
		if nerr != nil {
			if nerr.Code == 404 {
				log.Printf("Path %q not found: %s", r.URL.Path, nerr.Message)
				w.WriteHeader(404)
				err := Pages.ExecuteTemplate(w, "404.html", nil)
				if err != nil {
					http.NotFound(w, r)
				}
			} else if nerr.Code == 500 {
				log.Printf("Path %q error: %s", r.URL.Path, nerr.Message)
				s := make([]byte, 10000)
				runtime.Stack(s, false)
				log.Printf("%s", s)
				w.WriteHeader(500)
				err := Pages.ExecuteTemplate(w, "500.html", nerr.Message)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.Error(w, nerr.Message, nerr.Code)
			}
		}
	}
}

func ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) *NetError {
	err := Pages.ExecuteTemplate(w, name, data)
	if err != nil {
		return &NetError{500, err.Error()}
	}
	return nil
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

/*

type TweetButton struct {
	Text    string
	Via     string
	Related string
	Url     string
	Count   string
}
*/

/*
func (j *Joke) TweetButton() *TweetButton {
	t := &TweetButton{}
	t.Via = "barzedette"
	t.Related = "penpoe"
	if len(j.Joke) <= 140-len(" via @barzedette") {
		t.Text = j.Joke
		t.Count = "none"
		return t
	}

	suffix := "... via " + j.AbsUrl()
	t.Text = j.Joke[:min(len(j.Joke), 140-len(suffix))] + "..."
	t.Url = j.AbsUrl()
	return t
}
*/
