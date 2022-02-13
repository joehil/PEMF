package main

import (
    "html/template"
    "net/http"
    "io/ioutil"
    "fmt"
    "os"
    "syscall"
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
    Pemffactor string
}

type Answer struct {
    Frmethod string
    Frfile string
    Stage string
    Until string
    Pemffactor string
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
var pemffactor string = "100"
var curFile string = ""
var hasEnded bool = false

var lostart int = 0
var loend int = 0
var locnt int = 0
var lotime string = ""
var isLoop bool = false
var stopFlag bool = false

var isRunning bool = false

func main() {
    var data FormsData

    chome := os.Getenv("HOME")

    fmt.Println("HOME: "+chome)

    cpipe := os.Getenv("PIPE")

    fmt.Println("PIPE: "+cpipe)

    if len(os.Args)>1{
	if os.Args[1] == "generator"{
		fmt.Println("Generator started")

		cusb := os.Getenv("USBPORT")

		fmt.Println("USBPORT: "+cusb)

		cspeed := os.Getenv("USBSPEED")

		fmt.Println("USBSPEED: "+cspeed)

		if err := os.Remove(cpipe); err != nil && !os.IsNotExist(err) {
			fmt.Printf("remove: %v\n", err)
		}
		if err := syscall.Mkfifo(cpipe, 0644); err != nil {
			fmt.Printf("mkfifo: %v\n", err)
		}

		f, err := os.OpenFile(cpipe, os.O_RDONLY|syscall.O_NONBLOCK, 0644)
		if err != nil {
			fmt.Printf("open: %v\n", err)
		}
		defer f.Close()

		reader := bufio.NewReader(f)

    		speed, _ := strconv.Atoi(cspeed)

    		c := &serial.Config{Name: "/dev/serial/by-id/"+cusb, Baud: speed}
    		s, err := serial.OpenPort(c)
    		if err != nil {
        		fmt.Println(err)
			os.Exit(-1)
    		} 

		for true {
                	line, err := reader.ReadBytes('\n')
                	if err == nil {
				m := string(line)
				m = strings.ReplaceAll(m, "|", "\n")
                        	fmt.Printf("Message: %s", m)

               			_, err := s.Write([]byte(m+"\n"))
                		if err != nil {
                        		fmt.Println(err)
                		}
	                }
                        time.Sleep(1 * time.Second)
		} 

		os.Exit(0)
	}
    }

    data.Frmethod = "Audio"

    fmt.Println("Frequency server started")

    cport := os.Getenv("WEBPORT")

    fmt.Println("WEBPORT: "+cport)

    cfactor := os.Getenv("GENFACTOR")

    fmt.Println("GENFACTOR: "+cfactor)

    cgenport := os.Getenv("GENPORT")

    fmt.Println("GENPORT: "+cgenport)

    tmpl := template.Must(template.ParseFiles(chome+"/forms.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

        answer := Answer{
            Frmethod: r.FormValue("frmethod"),
            Frfile: r.FormValue("frfile"),
            Stage: r.FormValue("stage"),
	    Until: r.FormValue("loopuntil"),
	    Pemffactor: r.FormValue("pemffactor"),
        }

	if answer.Frfile != "" {
		curFile = answer.Frfile
	}

//	fmt.Println(answer)

        if answer.Frmethod == "" {
		answer.Frmethod = "Audio"
        }

        if answer.Pemffactor == "" {
                answer.Pemffactor = pemffactor
        }

        if answer.Stage == "" {
                answer.Stage = "Initial"
        }

        if answer.Stage == "Stop" {
		stopFlag = true
                answer.Stage = "Run"
		fmt.Println("Abort initiated")
        }

        data = FormsData{
		Title: "Kein Titel",
		Frmethod: answer.Frmethod, 
            	Frfiles: listDir(chome+"/data/"+answer.Frmethod),
		Frfile: answer.Frfile, 
		Stage: answer.Stage,
		Pemffactor: answer.Pemffactor,
        }

//	fmt.Println(data)

	if answer.Stage == "Success" {
		fmt.Println("File "+answer.Frfile+" chosen")
		switch answer.Frmethod {
        		case "Audio":
                        fmt.Println("Audio: "+answer.Frfile)
//        		procAudio("data/"+answer.Frfile)
        		case "FY2300":
			fmt.Println("FY2300: "+answer.Frfile, answer.Until)
        		go procFy2300(chome+"/data/FY2300/"+answer.Frfile,answer.Until,cfactor,answer.Pemffactor,cpipe,cgenport)
                        case "FY6900":
                        fmt.Println("FY6900: "+answer.Frfile, answer.Until)
                        go procFy2300(chome+"/data/FY6900/"+answer.Frfile,answer.Until,cfactor,answer.Pemffactor,cpipe,cgenport)
		        default:
        		fmt.Println("The command is wrong!")
			data.Stage = "Run"
    		}
	}

	if  data.Stage == "Run" || isRunning {
		data.TimeToGo = timeToGo
                data.Frequency = frequency
                data.Amplitude = amplitude
                data.Waveform = waveform
		data.Frfile = curFile
		data.Stage = "Run"
		if hasEnded {
			data.Stage = "Ended"
			hasEnded = false
			stopFlag = false
		}
	}

        tmpl.Execute(w, data)
    })

    http.ListenAndServe(":"+cport, nil)
}

func procFy2300(path string, loopuntil string, cfactor string, pemffactor string, cpipe string, cgenport string){
    var cser string
    var cint string
    var p []string

    isRunning = true
    loop := strings.Replace(loopuntil, ":", ".",-1)
    lines,err := readLines(path)

    if err != nil {
	fmt.Println(err)
    }

    for ind := 0; ind < len(lines); ind++ {
	cmd := lines[ind]
	if cgenport == "P"{
		cser, cint, p = parseFy2300_prim(cmd,cfactor,pemffactor)
	} else {
                cser, cint, p = parseFy2300_sec(cmd,cfactor,pemffactor)
	}
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
	                        if isLoop {
        	                        now := time.Now()
                	                tim := fmt.Sprintf("%02d.%02d",now.Hour(),now.Minute())
                        	        if n % 10 == 0 {
						fmt.Println(tim+" - "+lotime)
					}
                                	if tim == lotime {
                                        	ind = loend
						fmt.Println("Loop finished")
						n = limit + 1
                                	}
                        	}
				timeToGo = fmt.Sprintf("%d",limit-n)
                        	time.Sleep(1 * time.Second)
				if stopFlag {
					n = limit + 1
					ind = len(lines)
					if cgenport == "P"{ 
						cser = "WMN0"
					} else {
                                             	cser = "WFN0"
                                        } 
					fmt.Println("Process aborted")
				}
                	} 
		}
                if pt[0] == "lo" {
			lostart = ind
			fmt.Println("Loop initiated")
		}
                if pt[0] == "un" {
                        locnt++
			limit, _ := strconv.Atoi(pt[1])
			fmt.Println(limit,locnt)
			if limit > locnt {
				ind = lostart
			} else {
				fmt.Println("Loop finished")
			}
                }
                if pt[0] == "ti" {
			loend = ind 
			isLoop = true
			lotime = strings.Replace(pt[1], "<UNTIL>", loop, -1)
                        ind = lostart
                }
	}
	if cser != "" {
		fmt.Println(cser)
		writeGenerator(cser, cpipe)
	}
    }
    hasEnded = true
    isRunning = false
    lostart = 0
    loend = 0
    lotime = ""
    locnt = 0
    isLoop = false
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

func parseFy2300_prim(cmd string, cfactor string, pemffactor string) (string, string, []string) {
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
	freq,_ := strconv.ParseFloat(parts[1], 64)
        fact,_ := strconv.ParseFloat(cfactor, 64)
	if fact != 1 {
		freq *= fact
        	cser = "WMF"+fmt.Sprintf("%.0f",freq)
	} else {
		cser = "WMF"+parts[1]
	}
        case "am":
        ampl,_ := strconv.ParseFloat(parts[1], 64)
        fact,_ := strconv.ParseFloat(pemffactor, 64)
        fact = fact / 100
        ampl *= fact
        cser = "WMA"+fmt.Sprintf("%2.2f",ampl)
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

func parseFy2300_sec(cmd string, cfactor string, pemffactor string) (string, string, []string) {
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
	freq,_ := strconv.ParseFloat(parts[1], 64)
        fact,_ := strconv.ParseFloat(cfactor, 64)
	if fact != 1 {
		freq *= fact
        	cser = "WFF"+fmt.Sprintf("%.0f",freq)
	} else {
		cser = "WFF"+parts[1]
	}
        case "am":
        ampl,_ := strconv.ParseFloat(parts[1], 64)
        fact,_ := strconv.ParseFloat(pemffactor, 64)
        fact = fact / 100
        ampl *= fact
        cser = "WFA"+fmt.Sprintf("%2.2f",ampl)
        case "wv":
        cser = "WFW"+parts[1]
        case "on":
        cser = "WFN1"
        case "of":
        cser = "WFN0"
	case "##":
	break
        default:
        fmt.Println("The command is wrong!")
    }

    return cser, cint, parts
}

func writeGenerator(msg string, cpipe string) {
        f, err := os.OpenFile(cpipe, os.O_WRONLY|syscall.O_NONBLOCK, 0644)
        if err != nil {
                fmt.Printf("open: %v\n", err)
        }
        defer f.Close()

        _, err = f.WriteString(msg+"\n")

        if err != nil {
	        fmt.Println(err)
        }
}
