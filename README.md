**Project** is a very easy and flexible solution for building your own project generation tool. It will be helpful if have want to create a lot of services with the same file structure.

The main conсept based on the idea that you have a list of Job and Runner.

**Job** is the command responsible for creating files or directories. You can easily create your own command.

The **Runner** is functionality responsible for running all jobs one by one and error handling. It needs a list of jobs, an adapter for a file system, and adapter output. Runner also can be replaced by own implementation.

**Builder** is a tool for Job generation. It provides a more comfortable interface. Also, support traversing by folders and relative path for them.

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