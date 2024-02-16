package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	cmds := []Command{&GitCommand{}, &TarCommand{}, &S3Command{}, &CleanCommand{}}

	if err := verifyRuntime(cmds); err != nil {
		fmt.Println("[Failed]")
		panic(err)
	}

	options, err := parseFlags()
	if err != nil {
		panic(err)
	}

	for _, sourceRepository := range options.RepositoryList {

		options.Repository = sourceRepository
		for _, cmd := range cmds {

			cmd.WithOptions(options)
			if err := cmd.Exec(); err != nil {
				panic(err)
			}
		}
	}
}

// VerifyRuntime checks if requured external commands are available.
// git, zip, aws cli
func verifyRuntime(cmds []Command) error {

	fmt.Print("\nVerify runtime...")

	for _, cmd := range cmds {
		if err := cmd.Verify(); err != nil {
			return err
		}
	}

	fmt.Println("")
	return nil
}

func parseFlags() (Options, error) {

	options := Options{}

	repositories := flag.String("repositories", "repos.list", "Specify a file with source repositories on GitHub.")
	sourceRepository := flag.String("source-repository", "", "Specify a source repository on GitHub.")
	workDir := flag.String("work-directory", "", "Specify a directorz used for cloning.")

	targetBucket := flag.String("target-bucket", "", "Specify a S3 bucket for upload.")
	targetPath := flag.String("target-path", "", "Specify an optional target path.")

	flag.Parse()

	if sourceRepository != nil && *sourceRepository != "" {
		options.RepositoryList = []string{*sourceRepository}
	} else {
		content, err := os.ReadFile(*repositories)
		if err != nil {
			return options, err
		}
		options.RepositoryList = strings.Split(string(content), "\n")
	}
	options.Workdir = *workDir

	if len(options.RepositoryList) == 0 {
		return options, errors.New("No source repository passed!")
	}

	if targetBucket != nil && len(*targetBucket) > 0 {
		options.TargetBucket = *targetBucket
	} else {
		return options, errors.New("No target bucket specified!")
	}
	if targetPath != nil && len(*targetPath) > 0 {
		options.TargetPath = *targetPath
	}

	return options, nil
}
