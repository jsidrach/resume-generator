package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Resume format
type Resume struct {
	Name     string
	Title    string    `json:",omitempty"`
	Contact  Contact   `json:",omitempty"`
	Summary  string    `json:",omitempty"`
	Sections []Section `json:",omitempty"`
}

type Contact struct {
	Phone    string `json:",omitempty"`
	Address  string `json:",omitempty"`
	Email    string `json:",omitempty"`
	Webpage  Link   `json:",omitempty"`
	Linkedin Link   `json:",omitempty"`
	Github   Link   `json:",omitempty"`
}

type Link struct {
	Name string `json:",omitempty"`
	Url  string `json:",omitempty"`
}

type Section struct {
	Name    string  `json:",omitempty"`
	Entries []Entry `json:",omitempty"`
}

type Entry struct {
	What        string   `json:",omitempty"`
	Url         string   `json:",omitempty"`
	Where       string   `json:",omitempty"`
	When        string   `json:",omitempty"`
	Location    string   `json:",omitempty"`
	Description string   `json:",omitempty"`
	Details     []string `json:",omitempty"`
}

const (
	// Templates path
	TemplatesPath = "templates/tmpl."
	// Default input file
	DefaultInputFileName = "example.yaml"
	// Default output file(s), without extension
	DefaultOutputPath = "output/example"
)

// Templates extensions
// There must exist a template named TemplatesPath+Extension for each extension
var TemplatesExtensions = [...]string{"html", "md", "txt", "xml"}

func main() {
	// Logger
	l := log.New(os.Stderr, "", 0)

	// Check command line arguments
	lenArgs := len(os.Args)
	if lenArgs > 3 {
		l.Println("Invalid number of arguments")
		l.Println("Run using:")
		l.Fatalln("\tgo run resume-generator.go [<yaml-input> [<output-path>]]")
	}

	// Assign input and output file paths
	yamlFileName := DefaultInputFileName
	outputPath := DefaultOutputPath
	if lenArgs >= 2 {
		yamlFileName = os.Args[1]
	}
	if lenArgs == 3 {
		outputPath = os.Args[2]
	}

	// Read resume data, in YAML format
	yamlFile, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		l.Println("Invalid input filename: " + yamlFileName)
		l.Fatal(err)
	}
	resume := Resume{}
	if err := yaml.Unmarshal(yamlFile, &resume); err != nil {
		l.Println("Invalid input yaml: " + yamlFileName)
		l.Fatal(err)
	}

	// Generate resumes from templates
	for _, extension := range TemplatesExtensions {
		templateFileName := TemplatesPath + extension
		template, err := template.ParseFiles(templateFileName)
		if err != nil {
			l.Println("Invalid template filename: " + templateFileName)
			l.Fatal(err)
		}
		outputFileName := outputPath + "." + extension
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			l.Println("Invalid output filename: " + outputFileName)
			l.Fatal(err)
		}
		defer outputFile.Close()
		if err = template.Execute(outputFile, resume); err != nil {
			l.Println("Invalid template execution: " + templateFileName)
			l.Fatal(err)
		}
	}

	// Generate resume in json
	outputFileName := outputPath + ".json"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		l.Println("Invalid output filename: " + outputFileName)
		l.Fatal(err)
	}
	defer outputFile.Close()
	output, err := json.MarshalIndent(&resume, "", "  ")
	if err != nil {
		l.Println("Invalid json conversion")
		l.Fatal(err)
	}
	outputFile.Write(output)
}
