package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"
)

func main() {
	var err error

	// Flags
	var resumePath = flag.String("resume", "example.yaml", "Path to the resume YAML file")
	var browserRemote = flag.String("browser", "http://127.0.0.1:9222", "Remote instance of the browser (to print the resume as PDF)")
	var templatesPath = flag.String("templates", "templates/tmpl", "Path to the output templates (everything except the file extension)")
	var outputPath = flag.String("output", "output/example", "Path to the output files (everything except the file extension)")
	flag.Parse()

	// Load resume
	r, err := loadResume(*resumePath)
	if err != nil {
		log.Fatal(err)
	}

	// Templates extensions
	// There must exist a template named TemplatesPath+Extension for each extension
	var templatesExtensions = [...]string{".html", ".md", ".txt", ".xml"}

	// Save resume using the templates
	for _, extension := range templatesExtensions {
		err = r.saveAs((*outputPath)+extension, (*templatesPath)+extension)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Save resume as JSON
	err = r.saveAsJSON(*outputPath + ".json")
	if err != nil {
		log.Fatal(err)
	}

	// Save resume as PDF
	err = saveHTMLAsPDF(*browserRemote, *outputPath+".html", *outputPath+".pdf", defaultPrintToPDFParams())
	if err != nil {
		log.Fatal(err)
	}
}

//
// Resume structure
//

// Resume format
type Resume struct {
	Name     string    `json:",omitempty"`
	Title    string    `json:",omitempty"`
	Contact  Contact   `json:",omitempty"`
	Summary  string    `json:",omitempty"`
	Sections []Section `json:",omitempty"`
}

// Contact section
type Contact struct {
	Phone    string `json:",omitempty"`
	Address  string `json:",omitempty"`
	Email    string `json:",omitempty"`
	Webpage  Link   `json:",omitempty"`
	Linkedin Link   `json:",omitempty"`
	Github   Link   `json:",omitempty"`
}

// Link to URL
type Link struct {
	Name string `json:",omitempty"`
	URL  string `json:",omitempty"`
}

// Section of the resume
type Section struct {
	Name    string  `json:",omitempty"`
	Entries []Entry `json:",omitempty"`
}

// Entry of a section
type Entry struct {
	What        string   `json:",omitempty"`
	URL         string   `json:",omitempty"`
	Where       string   `json:",omitempty"`
	When        string   `json:",omitempty"`
	Location    string   `json:",omitempty"`
	Description string   `json:",omitempty"`
	Details     []string `json:",omitempty"`
}

//
// Load from YAML
//

func loadResume(yamlPath string) (*Resume, error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("Open input YAML (%s) failed\n%s", yamlPath, err)
	}
	resume := Resume{}
	err = yaml.Unmarshal(yamlFile, &resume)
	if err != nil {
		return nil, fmt.Errorf("Read input YAML (%s) failed\n%s", yamlPath, err)
	}
	return &resume, nil
}

//
// Save using a template
//

func (r *Resume) saveAs(outputPath string, templatePath string) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("Parse template file (%s) failed\n%s", templatePath, err)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Open output file (%s) failed\n%s", outputPath, err)
	}
	defer outputFile.Close()
	err = tmpl.Execute(outputFile, *r)
	if err != nil {
		return fmt.Errorf("Execute template (%s) failed\n%s", templatePath, err)
	}
	return nil
}

//
// Save as JSON
//

func (r *Resume) saveAsJSON(outputPath string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Open output file (%s) failed\n%s", outputPath, err)
	}
	defer outputFile.Close()
	output, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("Write output file (%s) failed\n%s", outputPath, err)
	}
	outputFile.Write(output)
	return nil
}

//
// Save as PDF
//

type printToPDFParams struct {
	// Paper orientation
	// If not present, defaults to false
	Landscape bool `json:"landscape,omitempty"`
	// Display header and footer
	// If not present, defaults to false
	DisplayHeaderFooter bool `json:"displayHeaderFooter,omitempty"`
	// Print background graphics
	// If not present, defaults to false
	PrintBackground bool `json:"printBackground,omitempty"`
	// Scale of the webpage rendering
	// If not present, defaults to 1
	Scale float64 `json:"scale"`
	// Paper width in inches
	// If not present, defaults to 8.5 inches
	PaperWidth float64 `json:"paperWidth"`
	// Paper height in inches
	// If not present, defaults to 11 inches
	PaperHeight float64 `json:"paperHeight"`
	// Top margin in inches
	// If not present, defaults to 1cm (~0.4 inches)
	MarginTop float64 `json:"marginTop"`
	// Bottom margin in inches
	// If not present, defaults to 1cm (~0.4 inches)
	MarginBottom float64 `json:"marginBottom"`
	// Left margin in inches
	// If not present, defaults to 1cm (~0.4 inches)
	MarginLeft float64 `json:"marginLeft"`
	// Right margin in inches
	// If not present, defaults to 1cm (~0.4 inches)
	MarginRight float64 `json:"marginRight"`
	// Paper ranges to print, e.g., '1-5, 8, 11-13'
	// If not present, defaults to the empty string, which means print all pages
	PageRanges string `json:"pageRanges,omitempty"`
	// Whether to silently ignore invalid but successfully parsed page ranges, such as '3-2'
	// If not present, defaults to false
	IgnoreInvalidPageRanges bool `json:"ignoreInvalidPageRanges,omitempty"`
}

func defaultPrintToPDFParams() *printToPDFParams {
	return &printToPDFParams{
		Landscape:           false,
		DisplayHeaderFooter: false,
		PrintBackground:     true,
		Scale:               1.0,
		PaperWidth:          8.5,
		PaperHeight:         11.0,
		MarginTop:           0.0,
		MarginBottom:        0.0,
		MarginLeft:          0.0,
		MarginRight:         0.0,
		PageRanges:          "",
	}
}

func saveHTMLAsPDF(browserRemote string, inputHTML string, outputPDF string, params *printToPDFParams) error {
	var err error

	// Obtain the address of a tab of the remote browser
	res, err := http.Get(browserRemote + "/json/list")
	if err != nil {
		return fmt.Errorf("Connect to remote browser failed\n%s", err)
	}
	defer res.Body.Close()
	rawTabs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Read list of tabs from remote browser failed\n%s", err)
	}
	var tabs []map[string]string
	err = json.Unmarshal(rawTabs, &tabs)
	if err != nil {
		return fmt.Errorf("Unmarshal response from remote browser failed\n%s", err)
	}
	if len(tabs) == 0 {
		return errors.New("Remote browser has no tabs to connect to")
	}
	wsURL, ok := tabs[0]["webSocketDebuggerUrl"]
	if !ok {
		return errors.New("Remote browser does not allow web socket connection to tab")
	}

	// Connect to the remote browser instance tab over a web socket
	maxBufferSize := 4 * 1024 * 1024
	d := &websocket.Dialer{
		ReadBufferSize:  maxBufferSize,
		WriteBufferSize: maxBufferSize,
	}
	ws, _, err := d.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("Connect to remote browser tab web socket failed\n%s", err)
	}
	defer ws.Close()

	// Navigate to the input html
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Get path of working directory failed\n%s", err)
	}
	inputPath := "file://" + path.Join(dir, inputHTML)
	navigateParams, err := json.Marshal(map[string]string{"url": inputPath})
	if err != nil {
		return fmt.Errorf("Marshal of parameters to navigate to page command failed\n%s", err)
	}
	_, err = wsCommand(ws, 0, "Page.navigate", navigateParams)
	if err != nil {
		return fmt.Errorf("Send request to open input html failed\n%s", err)
	}

	// Save page as PDF
	rawParams, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("Marshal of parameters to print pdf failed\n%s", err)
	}
	data, err := wsCommand(ws, 1, "Page.printToPDF", rawParams)
	if err != nil {
		return fmt.Errorf("Send request to print pdf failed\n%s", err)
	}
	err = ioutil.WriteFile(outputPDF, data, 0644)
	if err != nil {
		return fmt.Errorf("Write pdf failed\n%s", err)
	}
	return nil
}

// Web socket message
type wsMessage struct {
	ID     int64           `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	Result json.RawMessage `json:"result"`
}

func wsCommand(ws *websocket.Conn, id int64, method string, params json.RawMessage) ([]byte, error) {
	var err error

	// Serialize command
	cmd := wsMessage{
		ID:     id,
		Method: method,
		Params: params,
	}
	rawCmd, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("Marshal of command failed\n%s", err)
	}

	// Write command to socket
	err = ws.WriteMessage(websocket.TextMessage, rawCmd)
	if err != nil {
		return nil, fmt.Errorf("Write message to socket failed\n%s", err)
	}
	time.Sleep(time.Second)

	// Receive response from socket
	_, b, err := ws.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("Read message from socket failed\n%s", err)
	}

	// Parse response
	var r wsMessage
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal response from socket failed\n%s", err)
	}
	var res map[string]string
	err = json.Unmarshal(r.Result, &res)
	if err != nil {
		return nil, fmt.Errorf("Unmarshall response data from socket failed\n%s", err)
	}
	return base64.StdEncoding.DecodeString(res["data"])
}
