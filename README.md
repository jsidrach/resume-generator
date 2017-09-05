# Resume Generator

Generate elegant resumes from a single YAML data file.
Export them to HTML, JSON, Markdown, PDF, PNG, Text or XML.

Examples
----
* [Original YAML](example.yaml)
* [Generated HTML](output/example.html)
* [Generated JSON](output/example.json)
* [Generated Markdown](output/example.md)
* [Generated PDF](output/example.pdf)
* [Generated PNG](output/example.png)
* [Generated Text](output/example.txt)
* [Generated XML](output/example.xml)

Screenshot
![](output/example.png)

Prequisites
----
* [The Go Programming Language](https://golang.org/)
  * [go-yaml](https://github.com/go-yaml/yaml): install running `go get -u gopkg.in/yaml.v2`
  * [gorilla-websocket](https://github.com/gorilla/websocket): install running `go get -u github.com/gorilla/websocket`
* [Imagemagick](https://www.imagemagick.org/): required for PNG output
* Headless browser with remote debugging (e.g. [Chrome](https://www.google.com/chrome/browser/index.html)): required for PDF output

Run
----
* Make sure your data file is named `resume.yaml`, in the project's root folder (use `example.yaml` as a starting point).
* Generate the resume by executing `make resume`.
  * Edit the [Makefile](makefile) and substitute `chromium` with the binary for the headless browser.

License
----
[MIT](LICENSE) - Feel free to use and edit.

Tech
----
* [The Go Programming Language](https://golang.org/)
  * [go-yaml](https://github.com/go-yaml/yaml) - yaml parsing package
  * [gorilla-websocket](https://github.com/gorilla/websocket) - web sockets
* [Imagemagick](https://www.imagemagick.org/) - PDF to PNG converter
* [Font Awesome](https://fortawesome.github.io/Font-Awesome) - icons
