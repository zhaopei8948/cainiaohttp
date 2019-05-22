package main

import (
	"bytes"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
	"strings"
	"github.com/kataras/iris"
)

var (
	messageDir string
	suffix = "xml"
)

func handleMessage(message string) {
	buff := bytes.Buffer{}
	uid := uuid.NewV4()
	nowt := time.Now()
	strtime := nowt.Format("20060102150405")
	nano := nowt.UnixNano() % nowt.Unix()

	buff.WriteString(strtime)
	buff.WriteString(strconv.FormatInt(nano, 10))
	strtime = buff.String()

	buff.Reset()
	buff.WriteString(messageDir)
	buff.WriteByte(os.PathSeparator)
	// buff.WriteString("/")
	buff.WriteString(uid.String())
	buff.WriteString("_")
	buff.WriteString(strtime)
	buff.WriteString(".writing")
	fileName := buff.String()

	buff.Reset()
	buff.WriteString(messageDir)
	buff.WriteByte(os.PathSeparator)
	// buff.WriteString("/")
	buff.WriteString(uid.String())
	buff.WriteString("_")
	buff.WriteString(strtime)
	buff.WriteString(suffix)
	finalFileName := buff.String()

	// log.Printf("filename = [%s]\n", fileName)

	err := ioutil.WriteFile(fileName, []byte(message), 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.Rename(fileName, finalFileName)
	if err != nil {
		log.Fatal(err)
	} else {
		// log.Printf("success create file %s\n", finalFileName)
	}
}

func makeMessageDir() {
	if _, err := os.Stat(messageDir); os.IsNotExist(err) {
		os.MkdirAll(messageDir, 0755)
	}
}

func main() {
	arguments := os.Args
	port := "8080"
	if len(arguments) < 2 {
		log.Printf("usage: %s <message dir> [<port>] ", arguments[0])
		return
	}

	messageDir = arguments[1]
	makeMessageDir()

	if len(arguments) > 2 {
		port = arguments[2]
	}

	log.Printf("messageDir: %s, port: %s", messageDir, port)

	app := iris.Default()

	app.Post("/EcssTran/httpRequest/messagesHttpRequest", func(ctx iris.Context) {
		result := "<response><success>true</success></response>"
		message := ctx.FormValue("logistics_interface")
		// log.Printf("logistics_interface=[%s]\n", message)
		if strings.Index(message, "<Data>") == -1 ||
		strings.Index(message, "</Data>") == -1 {
			result = "<response><success>false</success><errorCode>S07</errorCode></response>"
		} else {
			handleMessage(message)
		}
		ctx.WriteString(result)
	})

	app.Run(iris.Addr(":" + port))
}
