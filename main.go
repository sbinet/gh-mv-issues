package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %[1]s:
$ %[1]s <src-repo> <dst-repo>

example:

$ %[1]s github.com/user/my-repo github.com/orga/repo

`,
			os.Args[0],
		)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}

	src := flag.Arg(0)
	dst := flag.Arg(1)

	fmt.Printf(
		"gh: moving issues...\ngh:   from [%s]\ngh:   to   [%s]\n",
		src,
		dst,
	)

	client := github.NewClient(nil)

	var srcIssues []github.Issue
	{
		toks := strings.SplitN(src, "/", 3)
		_, resp, err := client.Repositories.Get(toks[1], toks[2])
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
		log.Printf("resp: %#v\n", resp)
		issues, _, err := client.Issues.ListByRepo(
			toks[1], toks[2],
			&github.IssueListByRepoOptions{State: "all"},
		)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
		srcIssues = issues
	}

	for i, issue := range srcIssues {
		fmt.Printf("[%3d/%3d] #%d %v (%s)\n", i+1, len(srcIssues), *issue.Number, *issue.Title, *issue.State)
	}
}
