package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

func main() {
	switch len(os.Args) {
	case 2:
		if os.Args[1] != "init" {
			PrintUsage()
			return
		}
		InitProject();

	case 4:
		switch os.Args[1] {
		case "add":
			switch os.Args[2] {
			case "view": AddView(os.Args[3])
			case "model": AddModel(os.Args[3])
			case "font": AddFont(os.Args[3])
			default: PrintUsage()
			}
		case "remove":
			switch os.Args[2] {
			case "view": RemoveView(os.Args[3])
			case "model": RemoveModel(os.Args[3])
			case "font": RemoveFont(os.Args[3])
			default: PrintUsage()
			}
		}

	default:
		PrintUsage()
		return
	}
}

// PrintUsage prints some help text to the console.
func PrintUsage() {
	fmt.Printf(`
  Usage: %s <command>

  command:

  init - create a new requirejs/backbone/jquery/sass app

  add <subcommand> <args>

    subcommand:

    view <view name> - create a new view with the given name
    model <model name> - create a new model with the given name
    font <font string> - add a font from Google Web Fonts

  remove <subcommand> <args>

    subcommand:

    view <view name> - remove the view with the given name
    model <model name> - remove the model with the given name
    font <font string> - remove a font

`, path.Base(os.Args[0]))
}

// InitProject downloads the files from github.com/ryanlbrown/spapp
// into the current directory.
func InitProject() {
	fmt.Println("Downloading template...")
	resp, err := http.Get("https://github.com/ryanlbrown/spapp/archive/master.zip")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		panic(err)
	}
	for _, f := range r.File {
		// Ignore dot files.
		if path.Base(f.Name)[0] == '.' {
			continue
		}
		pathTokens := strings.Split(f.Name, "/")
		subName := strings.Join(pathTokens[1:], "/")
		if subName == "" || subName == "README.md" {
			continue
		}
		if subName[len(subName)-1] == '/' {
			dirName := subName[:len(subName)-1]
			fmt.Println("Created:", dirName)
			err = os.Mkdir(dirName, 0755)
			if err != nil && !os.IsExist(err) {
				panic(err)
			}
		} else {
			fa, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer fa.Close()
			fb, err := os.OpenFile(subName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
			if err != nil {
				panic(err)
			}
			defer fb.Close()
			_, err = io.Copy(fb, fa)
			if err != nil {
				panic(err)
			}
			fmt.Println("Created:", subName)
		}
	}
}

// AddView creates a new backbone view flie, a new sass file, a new
// template file, and adds a link tag to index.html.
func AddView(name string) {
	viewName := ViewName(name)
	contents := fmt.Sprintf(viewTemplate, TemplateFilename(name), viewName, ClassName(name), viewName)
	WriteFile(ViewFilename(name), contents, false);
	WriteFile(SassFilename(name), "", false);
	WriteFile(TemplateFilename(name), "", false);
	AddLinkTag(name);
}

// AddMovel creates a new backbone model file.
func AddModel(name string) {
	modelName := ModelName(name)
	contents := fmt.Sprintf(modelTemplate, modelName, modelName)
	WriteFile(ModelFilename(name), contents, false);
}

func AddFont(name string) {
	fontTag := FontTag(name)
	in := ReadLines("index.html")
	out := make([]string, 0)
	for _, line := range in {
		out = append(out, line)
		if strings.Contains(line, "<title") {
			out = append(out, fontTag)
		}
	}
	WriteFile("index.html", strings.Join(out, "\n") + "\n", true)
}

// RemoveView removes a backbone view file, its sass file, its
// template file, and the link tag in index.html.
func RemoveView(name string) {
	RemoveFile(ViewFilename(name))
	RemoveFile(SassFilename(name))
	RemoveFile(TemplateFilename(name))
	RemoveFile(CssFilename(name))
	RemoveLinkTag(name)
}

// RemoveModel removes a backbone model file.
func RemoveModel(name string) {
	RemoveFile(ModelFilename(name))
}

func RemoveFont(name string) {
	fontTag := FontTag(name)
	in := ReadLines("index.html")
	out := make([]string, 0)
	for _, line := range in {
		if strings.Contains(line, fontTag) {
			continue
		}
		out = append(out, line)
	}
	WriteFile("index.html", strings.Join(out, "\n") + "\n", true)
}

func AddLinkTag(name string) {
	linkTag := LinkTag(name)
	in := ReadLines("index.html")
	out := make([]string, 0)
	for _, line := range in {
		if strings.Contains(line, "<script") {
			out = append(out, linkTag)
		}
		out = append(out, line)
	}
	WriteFile("index.html", strings.Join(out, "\n") + "\n", true)
}

func RemoveLinkTag(name string) {
	linkTag := LinkTag(name)
	in := ReadLines("index.html")
	out := make([]string, 0)
	for _, line := range in {
		if strings.Contains(line, linkTag) {
			continue
		}
		out = append(out, line)
	}
	WriteFile("index.html", strings.Join(out, "\n") + "\n", true)
}

func LinkTag(name string) string {
	return fmt.Sprintf(`    <link rel="stylesheet" href="%s" />`, CssFilename(name))
}

func FontTag(name string) string {
	return fmt.Sprintf(`    <link rel="stylesheet" href="http://fonts.googleapis.com/css?family=%s" />`, name)
}

func ViewName(name string) string {
	return fmt.Sprintf("%sView", name)
}

func ModelName(name string) string {
	return fmt.Sprintf("%sModel", name)
}

func ViewFilename(name string) string {
	return fmt.Sprintf("js/%s.js", LowerName(ViewName(name)))
}

func ModelFilename(name string) string {
	return fmt.Sprintf("js/%s.js", LowerName(ModelName(name)))
}

func TemplateFilename(name string) string {
	return fmt.Sprintf("tmpl/%s.html", LowerName(name))
}

func SassFilename(name string) string {
	return fmt.Sprintf("sass/%s.sass", LowerName(name))
}

func CssFilename(name string) string {
	return fmt.Sprintf("css/%s.css", LowerName(name))
}

func ClassName(name string) string {
	return strings.Replace(LowerName(name), "_", "-", -1)
}

func LowerName(name string) string {
	b := bytes.NewBuffer([]byte{})
	for i, r := range name {
		if i == 0 {
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			b.WriteRune('_')
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func ReadLines(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	lines := make([]string, 0)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	return lines
}

func WriteFile(filename, contents string, exists bool) {
	flag := os.O_WRONLY|os.O_CREATE
	if exists {
		flag |= os.O_TRUNC
	} else {
		flag |= os.O_EXCL
	}
	f, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader([]byte(contents)))
	if err != nil {
		panic(err)
	}
	if exists {
		fmt.Println("Modified:", filename)
	} else {
		fmt.Println("Created:", filename)
	}
}

func RemoveFile(filename string) {
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	fmt.Println("Removed file:", filename)
}

var viewTemplate = `define([

  'jquery', 'underscore', 'base_view', 'text!%s',

], function($, _, BaseView, tmplText) {

  'use strict';

  var %s = BaseView.extend({

    tmpl: _.template(tmplText),
    tagName: 'div',
    className: '%s',

    initialize: function(options) {

    },

    render: function() {
      this.$el.html(this.tmpl());
      return this;
    },

  });

  return %s;

});
`

var modelTemplate = `define([

  'underscore', 'base_model',

], function(_, BaseModel) {

  'use strict';

  var %s = BaseModel.extend({

    defaults: {

    },

  });

  return %s;

});
`