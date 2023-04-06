package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Root struct {
	XMLName xml.Name `xml:"tt"`
	Body    Body     `xml:"body"`
}

type Body struct {
	XMLName xml.Name `xml:"body"`
	DIV     Div      `xml:"div"`
}

type Div struct {
	XMLName xml.Name `xml:"div"`
	P       []p      `xml:"p"`
}

type p struct {
	XMLName xml.Name `xml:"p"`
	Begin   string   `xml:"begin,attr"`
	End     string   `xml:"end,attr"`
	Content string   `xml:",innerxml"`
}

func FormatNetflixDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d",
		int(d.Seconds()/3600),
		int(d.Seconds()/60)%60,
		int(d.Seconds())%60,
		int(d.Milliseconds())%1000)
}

func main() {
	xmlFile, err := os.Open("in.xml")
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var root Root
	xml.Unmarshal(byteValue, &root)

	f, err := os.Create("out.srt")
	if err != nil {
		fmt.Println(err)
	}

	lineCount := 1
	for _, subLine := range root.Body.DIV.P {

		beginTime := ""
		endTime := ""

		if strings.Contains(subLine.Begin, "t") {
			beginTime = strings.ReplaceAll(subLine.Begin, "t", "")
			endTime = strings.ReplaceAll(subLine.End, "t", "")
			beginTime = beginTime[:len(beginTime)-1]
			endTime = endTime[:len(endTime)-1]
			dBeginTime, _ := time.ParseDuration(fmt.Sprintf("%sus", beginTime))
			dEndTime, _ := time.ParseDuration(fmt.Sprintf("%sus", endTime))

			beginTime = FormatNetflixDuration(dBeginTime)
			endTime = FormatNetflixDuration(dEndTime)

		} else {
			beginTime = strings.ReplaceAll(subLine.Begin, ".", ",")
			endTime = strings.ReplaceAll(subLine.End, ".", ",")
		}

		f.WriteString(fmt.Sprintf("%d\r\n", lineCount))
		f.WriteString(fmt.Sprintf("%s --> %s\r\n", beginTime, endTime))
		f.WriteString(fmt.Sprintf("%s\r\n\r\n", strings.ReplaceAll(subLine.Content, "<br />", "\n")))
		lineCount++
	}
	f.Close()

	fmt.Printf("\n")

}
