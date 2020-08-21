package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type Resources struct {
	XMLName xml.Name `xml:"resources"`
	Text    string   `xml:",chardata"`
	String  []struct {
		Text         string `xml:",chardata"`
		Name         string `xml:"name,attr"`
		Translatable string `xml:"translatable,attr"`
	} `xml:"string"`
}

func (r *Resources) exportCSV() [][]string {
	headers := []string{"key", "value", "translateable"}
	rows := [][]string{
		headers,
	}
	for _, s := range r.String {
		row := []string{s.Name, s.Text, s.Translatable}
		rows = append(rows, row)
	}
	return rows
}

func parseXML(path *string) Resources {
	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	resources := Resources{}
	err = xml.Unmarshal(data, &resources)
	if err != nil {
		log.Fatal(err)
	}
	resources.exportCSV()
	return resources
}

func main() {
	filePath := flag.String("file", "", "File to parse. (Required)")
	flag.Parse()

	if *filePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	resources := parseXML(filePath)
	writer := csv.NewWriter(os.Stdout)
	writer.WriteAll(resources.exportCSV())
}
