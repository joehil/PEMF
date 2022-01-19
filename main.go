package main

import (
    "html/template"
    "net/http"
    "io/ioutil"
    "fmt"
)

type FFiles struct {
    Frfile string
}

type FormsData struct {
    Title string
    Frfiles []FFiles
    Stage string
}

type Answer struct {
    Frmethod string
    Frfile string
    Stage string
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

        answer := Answer{
            Frmethod: r.FormValue("frmethod"),
            Frfile: r.FormValue("frfile"),
            Stage: r.FormValue("stage"),
        }

	fmt.Println(answer)

        data := FormsData{
		Title: "Kein Titel", 
            	Frfiles: listDir("data/"+answer.Stage),
		Stage: "Run",
        }

	fmt.Println(data)

        // do something with details
//        _ = details

        tmpl.Execute(w, data)
    })

    http.ListenAndServe(":8080", nil)
}
