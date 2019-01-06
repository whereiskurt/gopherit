package ui

import (
	"00-newapp-template/pkg/config"
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
)

// CLI makes the text output to the terminal.
type CLI struct {
	Config *config.Config
}

// NewCLI takes a configuration used for describing how to output.
func NewCLI(c *config.Config) (cli CLI) {
	cli.Config = c
	return
}

// DrawGopher outpus a text gopher to stdout
func (cli *CLI) DrawGopher() {
	fmt.Println(Gopher())
	return
}

// Gopher returns a string printable gopher! Thanks to belbomemo!
func Gopher() string {
	gopher := `
	         ,_---~~~~~----._         
	  _,,_,*^____      _____''*g*\"*, 
	 / __/ /'     ^.  /      \ ^@q   f 
	[  @f | @))    |  | @))   l  0 _/  
	 \'/   \~____ / __ \_____/    \   
	  |           _l__l_           I   
	  }          [______]           I  
	  ]            | | |            |  
	  ]             ~ ~             |  
	  |                            |   
	
	[[@https://gist.github.com/belbomemo]]
	`
	return gopher
}

// Render will output the UI templates as per the config bind the data.
func (cli *CLI) Render(name string, data interface{}) (usage string) {
	var raw bytes.Buffer
	var err error

	templateDir := cli.Config.TemplateFolder
	templateDir = strings.TrimSuffix(templateDir, "/")

	t := template.New("")
	t, err = t.Funcs(
		template.FuncMap{
			"Gopher": Gopher,
		},
	).ParseGlob(fmt.Sprintf("%s/ui/*.tmpl", templateDir))

	if err != nil {
		log.Fatalf("couldn't load template: %v", err)
	}

	err = t.ExecuteTemplate(&raw, name, data)
	if err != nil {
		log.Fatalf("error in Execute template: %v", err)
	}

	usage = raw.String()
	return
}
