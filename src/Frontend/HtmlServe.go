package frontend

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HtmlServe struct {
	Directory string
	views     map[string][]byte
}

func (this *HtmlServe) CacheHtml() {
	this.views = make(map[string][]byte)
	dir, err := ioutil.ReadDir(this.Directory)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(dir); i++ {
		out, err := ioutil.ReadFile(this.Directory + "/" + dir[i].Name())
		this.views[dir[i].Name()] = out
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (this *HtmlServe) Serve(w *http.ResponseWriter, fileName string) {
	//view, _ := ioutil.ReadFile("view/index.html")
	//io.WriteString(*w, string(view))
	//
	io.WriteString(*w, string(this.views[fileName]))
}
