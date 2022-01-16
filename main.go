package main

import (
    "html/template"
    "net/http"
    "io/ioutil"
)

type FFiles struct {
    Frfile string
}

type FormsData struct {
    Title string
    Frfiles []FFiles
    Success bool
}

func listDir(direc string) (trfiles []FFiles) {
    files, err := ioutil.ReadDir(direc)
    if err == nil {
	    for _, file := range files {
                tfile := FFiles{
			Frfile: file.Name(),
		}
                trfiles = append(trfiles,tfile)
    	    }
   }
   return
}

func main() {
    tmpl := template.Must(template.ParseFiles("forms.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            data := FormsData{
		Title: "Kein Titel", 
            	Frfiles: listDir("data/FY2300"),
		Success: false,
	    }
            tmpl.Execute(w, data)
            return
        }

        data := FormsData{
	    Title: "Titel",
            Frfiles: []FFiles{
                {Frfile: "Audio1"},
                {Frfile: "Audio2"},
            },
            Success: true,
        }

        // do something with details
//        _ = details

        tmpl.Execute(w, data)
    })

    http.ListenAndServe(":8080", nil)
}
