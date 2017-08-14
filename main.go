package main

import (
	"log"
	"net/http"
	"html/template"
	"io"
	"io/ioutil"
	"flag"
	"fmt"
	"github.com/dapus/dirlist"
	"github.com/dapus/www/gitrepos"
	"github.com/russross/blackfriday"
)

const tpl = `
<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Host }}</title>
    <style>
      body {
        font-family: sans-serif;
		color: rgb(51, 51, 51);
      }

      a, a:hover, a:visited {
        color: inherit;
        text-decoration: underline;
      }

	  .section-sep {
		  font-size: 200%;
		  margin-top: 1em;
		  margin-bottom: 1em;
	  }
    </style>
  </head>
  <body>
	<h1>{{ .Url.Path }}</h1>

	<div>{{ with .Index }}{{ markdown . }}{{ end }}</div>

	<div class="section-sep">🕴</div>

    {{ range .Files -}}
      {{ if ne (indexstr .Name 0) "." -}}
      <a href="{{ $.Url.Path }}{{ .Name }}{{ if .IsDir }}/{{ end }}">
        {{ .Name }}{{ if .IsDir }}/{{ end }}
      </a><br>
      {{ end }}
    {{- end }}
  </body>
</html>
`

var tplFuncs = template.FuncMap{
	"indexstr": func(str string, idx int) string {
		return string(str[idx])
	},
	"markdown": func(r io.Reader) template.HTML {
		data, _ := ioutil.ReadAll(r)

		return template.HTML(blackfriday.MarkdownCommon(data))
	},
}

var listenAddr string
var serveDir string
var serveGit string

func init() {
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8080", "Address to listen to")
	flag.StringVar(&serveDir, "dir", "", "Directory to serve")
	flag.StringVar(&serveGit, "git", "", "Git directory to serve")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path;
		rw := &ResponseWriter{w, 0}
		handler.ServeHTTP(rw, r);
		log.Printf("%d %s %s", rw.Status(), r.Method, url)
	})
}

func main() {
	flag.Parse()

	if serveDir == "" {
		fmt.Print("-dir required")
		return
	}

	tpl, err := template.New("dir").Funcs(tplFuncs).Parse(tpl)

	if err != nil {
		log.Fatal(err)
	}

	dirHandle := &dirlist.DirList{
		http.Dir(serveDir),
		tpl,
		[]string{"index.md"},
	}

	http.Handle("/", logRequest(dirHandle))

	if serveGit != "" {
		gitHandle := &dirlist.DirList{
			gitrepos.GitRepos(serveGit),
			tpl,
			[]string{"README.md", "README"},
		}

		http.Handle("/git/", logRequest(http.StripPrefix("/git", gitHandle)))
	}

	log.Printf("Listening to %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil));
}
