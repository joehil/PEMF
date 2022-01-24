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
    "strconv"
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
    Frequency string
    Amplitude string
    Waveform string
    TimeToGo string
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

var timeToGo string = "0"
var frequency string = "0"
var amplitude string = "0"
var waveform string = "0"
var hasEnded bool = false

var lostart int = 0
var locnt int = 0

func main() {
    var data FormsData

    data.Frmethod = "Audio"

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

        if answer.Stage == "" {
                answer.Stage = "Initial"
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
        		go procFy2300("data/FY2300/"+answer.Frfile)
		        default:
        		fmt.Println("The command is wrong!")
			data.Stage = "Run"
    		}
	}

	if  data.Stage == "Run" {
		data.TimeToGo = timeToGo
                data.Frequency = frequency
                data.Amplitude = amplitude
                data.Waveform = waveform
		if hasEnded {
			data.Stage = "Ended"
			hasEnded = false
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

    for ind := 0; ind < len(lines); ind++ {
	cmd := lines[ind]
	cser, cint, p := parseFy2300(cmd)
	switch p[0] {
		case "fr":
		frequency = p[1]
		case "am":
		amplitude = p[1]
		case "wv":
		waveform = p[1]
		default:
	}
	if cint != "" {
		pt := strings.Split(cint, ":")
		if pt[0] == "do" {
			timeToGo = cint
			limit, err := strconv.Atoi(pt[1])
   			if err != nil {
        			fmt.Println(err)
	    		} 
			for n:=0; n<limit; n++ {
				timeToGo = fmt.Sprintf("%d",limit-n)
                        	time.Sleep(1 * time.Second)
                	} 
		}
                if pt[0] == "lo" {
			lostart = ind
		}
                if pt[0] == "un" {
                        locnt++
			limit, _ := strconv.Atoi(pt[1])
			fmt.Println(limit,locnt)
			if limit > locnt {
				ind = lostart
			}
                }
                if pt[0] == "ti" {
                        now := time.Now()
			tim := fmt.Sprintf("%02d:%02d",now.Hour(),now.Minute())
			fmt.Println(tim+" - "+p[1])
                        if strings.Compare(pt[1], tim) < 0 {
                                ind = lostart
                        }
                }
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
    hasEnded = true
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

func parseFy2300(cmd string) (string, string, []string) {
    var cser string = ""
    var cint string = ""

    parts := strings.Split(cmd, " ")

    switch parts[0] {
        case "do":
        cint = "do:"+parts[1]
        case "lo":
        cint = "lo"
        case "un":
        cint = "un:"+parts[1]
        case "ti":
        cint = "ti:"+parts[1]
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
	case "##":
	break
        default:
        fmt.Println("The command is wrong!")
    }

    return cser, cint, parts
}
