package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"html/template"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)
type content struct {
	Title string 
	Body template.HTML
}
const (
	defaultTemplate = `<!DOCTYPE html>
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>{{ .Title }}</title>
	</head>
	<body>
	{{ .Body }}
	</body>
	</html>
	`
)
	
const (
	header = `<!DOCTYPE html>
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>Markdown Preview Tool</title>
	</head>
	<body>
	`

	footer = `
	</body>
	</html>
	`
)
func main(){
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(file string, tFname string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(file)
	if err != nil {
		return err 
	}
	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err 
	}
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return nil 
	}
	if err := temp.Close(); err != nil {
		return err 
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)
	if err := saveHtml(outName, htmlData); err != nil {
		return err 
	}
	if skipPreview {
		return nil 
	}
	defer os.Remove(outName)
	return preview(outName)
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}
	c := content {
		Title: "Markdown Preview Tool",
		Body: template.HTML(body),
	}
	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.WriteString(string(body))
	buffer.WriteString(footer)
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func saveHtml(outfile string, data []byte) error {
	return os.WriteFile(outfile, data, 0644)
}

func preview(fname string) error {
	cName := "cmd.exe"
	cParams := []string{}
	cParams = append(cParams, fname)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err 
	}
	return exec.Command(cPath, cParams...).Run()
}

