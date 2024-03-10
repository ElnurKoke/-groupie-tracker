package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandlerError(w, http.StatusMethodNotAllowed)
		return
	}

	artists, err := ParseJson(w, r)
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	if r.URL.Path != "/" {
		HandlerError(w, 404)
		return
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		HandlerError(w, 500)
		return
	}
}

func Artist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandlerError(w, http.StatusMethodNotAllowed)
		return
	}

	artists, err := ParseJson(w, r)

	pattern := `/artists/(\d{1,2}\z)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(r.URL.Path)

	if len(matches) != 2 {
		HandlerError(w, 404)
		return
	}
	ID, err := strconv.Atoi(matches[1])
	if err != nil {
		HandlerError(w, 404)
		return
	}
	if ID > len(artists) || ID <= 0 {
		HandlerError(w, 404)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/artist.html"))
	err = tmpl.Execute(w, artists[ID-1])
	if err != nil {
		HandlerError(w, 404)
		return
	}
}

func HandlerError(w http.ResponseWriter, code int) {
	pageError := struct {
		ErrorNum     int
		ErrorMessage string
	}{
		ErrorNum:     code,
		ErrorMessage: http.StatusText(code),
	}
	// w.WriteHeader(san)
	temp, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	temp.Execute(w, pageError)
}

func HandleRequest() {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", Home)
	http.HandleFunc("/artists/", Artist)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
