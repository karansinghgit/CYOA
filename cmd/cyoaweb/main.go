package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/karansinghgit/cyoa"
)

var storyViewTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>
		CYOA
	</title>
</head>
<body>
	<section class="page">
		<h1>{{.Title}}</h1>
		{{range .Paragraphs }}
			<p>{{.}}</p>
		{{end}}
		<hr>
		<ul>
		{{range .Options}}
			<li><a href="/{{.Arc}}">{{.Text}}</a></li>
		{{end}}
		</ul>
	</section>
	<style>
		h1 {
			text-align:center;
		}
		.page {
			width: 40%;
			background: #FFFDA2;
			margin: auto;
			font-size: 20px;
			padding: 40px;
		}
	</style>
</body>
</html>
`
var tpl *template.Template
var mux *http.ServeMux

type Story map[string]cyoa.Chapter

func init() {
	tpl = template.Must(template.New("").Parse(storyViewTemplate))
	mux = http.NewServeMux()
}

func myhandlefunc(s Story) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "" {
			path = "/intro"
		}
		path = path[1:]
		if arc, ok := s[path]; ok {
			err := tpl.Execute(w, arc)
			if err != nil {
				http.Error(w, "Execute nahi hua", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Arc Not Found", http.StatusNotFound)
		}
	})
}

func main() {
	port := flag.Int("port", 3000, "Port number to open the server on")
	filename := flag.String("filename", "story.json", "JSON file which contains story")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	var story Story
	json.NewDecoder(f).Decode(&story)

	// This for loop can be ignored, but it may be good to register the routes from before. IDK.
	// for chapterName := range story {
	// 	url := "/" + chapterName
	// 	mux.HandleFunc(url, myhandlefunc(story))
	// }

	mux.HandleFunc("/", myhandlefunc(story))
	fmt.Println("Starting SERVER ON PORT 3000")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), mux))
}
