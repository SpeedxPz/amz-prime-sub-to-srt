package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
		f.WriteString(fmt.Sprintf("%d\r\n", lineCount))
		f.WriteString(fmt.Sprintf("%s --> %s\r\n", strings.ReplaceAll(subLine.Begin, ".", ","), strings.ReplaceAll(subLine.End, ".", ",")))
		f.WriteString(fmt.Sprintf("%s\r\n\r\n", strings.ReplaceAll(subLine.Content, "<br />", "\n")))
		lineCount++
	}
	f.Close()

	fmt.Printf("%+v\n")

}
