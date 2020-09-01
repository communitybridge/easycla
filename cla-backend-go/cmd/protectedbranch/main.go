package main

import (
	"context"
	"errors"
	"flag"
	github2 "github.com/communitybridge/easycla/cla-backend-go/github"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils/github"
	"os"
)

var repoName = flag.String("repo", "", "the repo name to work on")
var branch = flag.String("branch", "", "the owner of the repo")
var orgName = flag.String("org", "", "organization name")
var enable = flag.Bool("enable", false, "if not enabled tries to enable")

func main() {
	flag.Parse()

	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if githubAccessToken == "" {
		log.Fatalf("GITHUB_ACCESS_TOKEN is required")
	}

	if *repoName == "" || *orgName == "" {
		log.Fatal("repo, owner and org are required")
	}

	ctx := context.Background()

	client := github.NewClient(githubAccessToken)
	//org, err := github.GetOrganization(ctx, client, *orgName)
	//if err != nil {
	//	log.Fatalf("org name fetch : %v", err)
	//}

	//if err := github.GetRepos(ctx, client, *orgName); err != nil {
	//	log.Fatalf("fetching repos for org : %v", err)
	//}

	//log.Printf("the org info is like : %+v\n", org)

	owner, err := github2.GetOwnerName(ctx, client, *orgName, *repoName)
	if err != nil {
		log.Fatalf("fetching the owner name : %v", err)
	}

	branchName := *branch

	if branchName == ""{
		branchName, err = github2.GetDefaultBranchForRepo(ctx, client, owner, *repoName)
		if err != nil {
			log.Fatalf("default branch : %v", err)
		}
	}

	protected := true
	protectedBranch, err := github2.GetProtectedBranch(ctx, client, owner, *repoName, branchName)
	if err != nil {
		if !errors.Is(err, github2.BranchNotProtectedError){
			log.Fatalf("fetching the protected branch : %v", err)
		}
		protected = false
	}
	if protected{
		log.Println("the branch is protected : ", branchName)
	}else{
		log.Println("the branch is not protected : ", branchName)
	}

	if protected{
		if github2.IsEnforceAdminEnabled(protectedBranch){
			log.Println("enforce admin is enabled")
		}else{
			log.Println("enforce admin is disabled")
		}

		if github2.AreStatusChecksEnabled(protectedBranch, []string{"ci/circleci: run"}){
			log.Println("status checks are enabled")
		}else{
			log.Println("status checks are disabled or not all of them enabled")
		}
	}

	if *enable{
		err = github2.EnableBranchProtection(ctx, client, owner, *repoName, branchName, true, []string{"ci/circleci: run"})
		if err != nil {
			log.Fatalf("enabling branch protection failed : %v", err)
		}
		log.Println("branch protection with all the options enabled")
	}


}
