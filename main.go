package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")

	if githubAccessToken == "" {
		panic("No github auth")
	}

	var versionFlag bool
	flag.BoolVar(&versionFlag, "v", false, "Alias for -version")
	flag.BoolVar(&versionFlag, "version", false, "Print version info and exit")

	var filename string
	flag.StringVar(&filename, "f", "", "Alias for -filename")
	flag.StringVar(&filename, "filename", "", "Name for the gist file to be created. Required.")

	var silentFlag bool
	flag.BoolVar(&silentFlag, "s", false, "Alias for -silent")
	flag.BoolVar(&silentFlag, "silent", false, "Don't produce any output")

	flag.Parse()

	if versionFlag {
		fmt.Printf("gist v0.0.1\n")
		os.Exit(0)
	}

	if filename == "" {
		fmt.Printf("-filename is required!\n")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	input, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		panic(err)
	}

	contents := string(input)

	gistfiles := map[github.GistFilename]github.GistFile{}
	gistfiles[github.GistFilename(filename)] = github.GistFile{
		Filename: &filename,
		Content:  &contents,
	}

	gist := github.Gist{
		Files: gistfiles,
	}

	resultGist, _, err := client.Gists.Create(ctx, &gist)

	if err != nil {
		panic(err)
	}

	if silentFlag == false {
		fmt.Printf("%s available at %s\n", filename, resultGist.GetHTMLURL())
	}
}
