# Project

This is a pretty easy and flexible library for generating a project structure. 
Can be used in own command-line tools for starting new projects with the same structures. 
Supports customizations and is transportable to different platforms.

## Overview
All components can be replaced or extended.

**Job** - the main work item. Contains function that will be executed and interactive description.

**Builder** - tool for generation jobs list in comfortable way

**Runner** - responsible for run provided jobs 

**Output** - responsible for displaying progress of work

**Driver** - describe relation with file system

## Example

```go
package main

import (
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/oherych/project"
)

var templateMain = template.Must(template.New("main").Parse(`package main

import (
	"fmt"
	"net/http"
)

const (
	serviceName = "{{.serviceName}}"
	port        = "{{.port}}"
)

func main() {

	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintln(w, serviceName)
	})

	fmt.Println("starting " + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println(err)
		return
	}
}

`))

func main() {

	flag.Parse()

	rootFolder, err := getPathToRoot()
	if err != nil {
		log.Fatalln(err)
	}

	serviceName := flag.String("name", "demo", "service name")

	jobs := prepareJobs(*serviceName)

	if err := generate(rootFolder, jobs); err != nil {
		log.Fatalln(err)
	}
}

func prepareJobs(serviceName string) []project.Job {
	const perm = 0777

	b := project.NewBuilder()
	b.DirIn("cmd", perm)
	b.Template("main.go", perm, templateMain, map[string]interface{}{
		"serviceName": serviceName,
		"port":        ":8080",
	})

	return b.Done()
}

func generate(rootFolder string, jobs []project.Job) error {
	r := project.NewRunner(jobs, project.FileSystem{Root: rootFolder}, project.ConsoleOutput{Writer: os.Stdout})

	return r.Run()
}

func getPathToRoot() (string, error) {
	if len(os.Args) >= 2 {
		return os.Args[1], nil
	}

	return os.Getwd()
}

```