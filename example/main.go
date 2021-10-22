package main

import (
	"log"
	"os"

	"github.com/oherych/project"
)

func main() {
	rootFolder, err := getPathToRoot()
	if err != nil {
		log.Fatalln(err)
	}

	jobs := prepareJobs()

	if err := generate(rootFolder, jobs); err != nil {
		log.Fatalln(err)
	}
}

func prepareJobs() []project.Job {
	const perm = 0777

	return project.NewBuilder().
		DirIn("cmd", perm).
		Touch("main.go", perm).
		Done()
}

func generate(rootFolder string, jobs []project.Job) error {
	r := project.NewRunner(jobs, project.FileSystem{Root: rootFolder}, project.ConsoleFeedback{Writer: os.Stdout})

	return r.Run()
}

func getPathToRoot() (string, error) {
	if len(os.Args) >= 2 {
		return os.Args[1], nil
	}

	return os.Getwd()
}
