package main

import (
    "html/template"
    "net/http"
    "io/ioutil"
    "fmt"
    "os"
    "time"
    "bufio"
    "strings"
    "github.com/tarm/serial"
)

type FFiles struct {
    Frfile string
}

type FormsData struct {
    Title string
    Frmethod string
    Frfiles []FFiles
    Frfile string
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
    data := FormsData{
	Title: "Kein Titel",
	Frmethod : "Audio",
	Stage: "Run",
    }

    fmt.Println("Frequency server started")

    tmpl := template.Must(template.ParseFiles("forms.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

        answer := Answer{
            Frmethod: r.FormValue("frmethod"),
            Frfile: r.FormValue("frfile"),
            Stage: r.FormValue("stage"),
        }

//	fmt.Println(answer)

        if answer.Frmethod == "" {
		answer.Frmethod = "Audio"
        }

        data = FormsData{
		Title: "Kein Titel",
		Frmethod: answer.Frmethod, 
            	Frfiles: listDir("data/"+answer.Frmethod),
		Frfile: answer.Frfile, 
		Stage: answer.Stage,
        }

//	fmt.Println(data)

	if answer.Stage == "Success" {
		fmt.Println("Datei "+answer.Frfile+" ausgewählt")
		switch answer.Frmethod {
        		case "Audio":
                        fmt.Println("Audio: "+answer.Frfile)
//        		procAudio("data/"+answer.Frfile)
        		case "FY2300":
			fmt.Println("FY2300: "+answer.Frfile)
        		procFy2300("data/FY2300/"+answer.Frfile)
		        default:
        		fmt.Println("The command is wrong!")
    		}
	}

        tmpl.Execute(w, data)
    })

    http.ListenAndServe(":8080", nil)
}

func procFy2300(path string){
    lines,err := readLines(path)

    if err != nil {
	fmt.Println(err)
    }

    c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
    s, err := serial.OpenPort(c)
    if err != nil {
    	fmt.Println(err)
    } 

    for _, cmd := range lines {
	cser, cint := parseFy2300(cmd)
	if cint != "" {
		t,_ := time.ParseDuration(cint+"s")
		time.Sleep(t)
	}
	if cser != "" {
		fmt.Println(cser)
        	_, err := s.Write([]byte(cser+"\n"))
        	if err != nil {
                	fmt.Println(err)
        	}
		time.Sleep(1 * time.Second)
	}
    }

}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func parseFy2300(cmd string) (string, string) {
    var cser string = ""
    var cint string = ""

    parts := strings.Split(cmd, " ")

    switch parts[0] {
        case "do":
        cint = parts[1]
        case "fr":
        cser = "WMF"+parts[1]
        case "am":
        cser = "WMA"+parts[1]
        case "wv":
        cser = "WMW"+parts[1]
        case "on":
        cser = "WMN1"
        case "of":
        cser = "WMN0"
        default:
        fmt.Println("The command is wrong!")
    }

    return cser, cint
}
