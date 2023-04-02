package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"path/filepath"
)

type prConf struct {
	organization  string
	origin        string
	featureBranch string
	targetBranch  string
	prTitle       string
	prBody        string
	commitMsg     string
	prReviewer    []string
	token         string
}

func main() {

	// paths of your local repo
	paths := []string{
		"C:\\Users\\myName\\Documents\\GitHub\\my_repo_1",
		"C:\\Users\\myName\\Documents\\GitHub\\my_repo_2",
		"C:\\Users\\myName\\Documents\\GitHub\\my_repo_3"}

	// configuration of the PR
	conf := prConf{
		organization:  "myOrg",
		origin:        "origin/Rel-1.22",
		featureBranch: "myfeatureBranch/Rel-1.22",
		targetBranch:  "Rel-1.22",
		prTitle:       "PR 230329-01 from Rel-1.22 to Rel-1.22", // PR title
		prBody:        "",                                       //PR body
		commitMsg:     "",
		prReviewer:    []string{""}, // add the reviewers name
		token:         "",           //gh token
	}

	//// for dev to stage
	devToStage(paths, conf)

	//// for updating dev
	//devUpdateDev(paths, conf)

}

func devToStage(paths []string, conf prConf) {

	var prLinks []string

	for i := 0; i < len(paths); i++ {

		repo := filepath.Base(paths[i])

		// Change directory to the Git repository
		err := os.Chdir(paths[i])
		if err != nil {
			fmt.Println("Error changing directory:", err)
			return
		}

		// Git fetch origin
		gitFetchCmd := exec.Command("git", "fetch", "origin")
		gitFetchCmd.Stdout = os.Stdout
		gitFetchCmd.Stderr = os.Stderr
		err = gitFetchCmd.Run()
		if err != nil {
			fmt.Println("Error running git fetch:", err)
			return
		}

		// Git checkout -b cf/230323-01-dev-R5.11 origin/dev-R5.11
		gitCheckoutCmd := exec.Command("git", "checkout", "-b", conf.featureBranch, conf.origin)
		gitCheckoutCmd.Stdout = os.Stdout
		gitCheckoutCmd.Stderr = os.Stderr
		err = gitCheckoutCmd.Run()
		if err != nil {
			fmt.Println("Error running git checkout:", err)
			return
		}

		// Git status -sb
		gitStatusCmd := exec.Command("git", "status", "-sb")
		gitStatusCmd.Stdout = os.Stdout
		gitStatusCmd.Stderr = os.Stderr
		err = gitStatusCmd.Run()
		if err != nil {
			fmt.Println("Error running git status:", err)
			return
		}

		// Git merge origin/stage-R5.11
		gitMergeCmd := exec.Command("git", "merge", "origin/"+conf.targetBranch)
		gitMergeCmd.Stdout = os.Stdout
		gitMergeCmd.Stderr = os.Stderr
		err = gitMergeCmd.Run()
		if err != nil {
			fmt.Println("Error running git merge:", err)
			return
		}

		// Git status
		gitStatusCmd = exec.Command("git", "status")
		gitStatusCmd.Stdout = os.Stdout
		gitStatusCmd.Stderr = os.Stderr
		err = gitStatusCmd.Run()
		if err != nil {
			fmt.Println("Error running git status:", err)
			return
		}

		// Git push -u origin cf/230323-01-dev-R5.11
		gitPushCmd := exec.Command("git", "push", "-u", "origin", conf.featureBranch)
		gitPushCmd.Stdout = os.Stdout
		gitPushCmd.Stderr = os.Stderr
		err = gitPushCmd.Run()
		if err != nil {
			fmt.Println("Error running git push:", err)
			return
		}

		// create PR
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: conf.token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		newPR := &github.NewPullRequest{
			Title:               github.String(conf.prTitle),
			Head:                github.String(conf.featureBranch),
			Base:                github.String(conf.targetBranch),
			Body:                github.String(conf.prBody),
			MaintainerCanModify: github.Bool(true),
		}

		pr, _, err := client.PullRequests.Create(ctx, conf.organization, repo, newPR)
		if err != nil {
			fmt.Println(err)
			return
		}
		prLinks = append(prLinks, pr.GetHTMLURL())

		reviewers := github.ReviewersRequest{
			//NodeID:    github.String(*pr.NodeID),
			Reviewers: conf.prReviewer,
			//TeamReviewers: []string{""},
		}

		_, _, err = client.PullRequests.RequestReviewers(ctx, conf.organization, repo, *pr.Number, reviewers)
		if err != nil {
			fmt.Println("Error adding reviewer:", err)
			return
		}

	} //end for loop

	fmt.Println("===========================================================")
	fmt.Println("Here are the PRs link...")
	printArray(prLinks)
}

func printArray(arr []string) {
	// using for loop
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}

func devUpdateDev(paths []string, conf prConf) {

	var prLinks []string

	for i := 0; i < len(paths); i++ {

		//repo := filepath.Base(paths[i])

		// Change directory to the Git repository
		err := os.Chdir(paths[i])
		if err != nil {
			fmt.Println("Error changing directory:", err)
			return
		}

		// Git fetch origin
		gitFetchCmd := exec.Command("git", "fetch", "origin")
		gitFetchCmd.Stdout = os.Stdout
		gitFetchCmd.Stderr = os.Stderr
		err = gitFetchCmd.Run()
		if err != nil {
			fmt.Println("Error running git fetch:", err)
			return
		}

		// Git checkout -b cf/230323-01-dev-R5.11 origin/dev-R5.11
		gitCheckoutCmd := exec.Command("git", "checkout", "-b", conf.featureBranch, conf.origin)
		gitCheckoutCmd.Stdout = os.Stdout
		gitCheckoutCmd.Stderr = os.Stderr
		err = gitCheckoutCmd.Run()
		if err != nil {
			fmt.Println("Error running git checkout:", err)
			return
		}

		// Git push -u origin cf/230323-01-dev-R5.11
		gitPushCmd := exec.Command("git", "push", "-u", "origin", conf.featureBranch)
		gitPushCmd.Stdout = os.Stdout
		gitPushCmd.Stderr = os.Stderr
		err = gitPushCmd.Run()
		if err != nil {
			fmt.Println("Error running git push:", err)
			return
		}

	} //end for loop

	// Wait for user input to continue to the next step
	var input string
	fmt.Println("Type 'ok' to proceed to the merge and create PR")
	fmt.Scanln(&input)

	// initiate github client service
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Only proceed to the second array if the user inputs "ok"
	if input == "ok" {

		fmt.Println("Proceeding to the next step...")
		for i := 0; i < len(paths); i++ {
			repo := filepath.Base(paths[i])
			//	fmt.Println(repo)

			// Change directory to the Git repository
			err := os.Chdir(paths[i])
			if err != nil {
				fmt.Println("Error changing directory:", err)
				return
			}

			// Check current branch
			cmd := exec.Command("git", "branch")
			out, err := cmd.Output()
			if err != nil {
				fmt.Println("error checking current branch:", err)
				return
			}
			fmt.Println(string(out))

			// Switch to feature branch
			cmd = exec.Command("git", "checkout", conf.featureBranch)
			err = cmd.Run()
			if err != nil {
				fmt.Println("error switching branch:", err)
				return
			}

			// Stagsh changes
			cmd = exec.Command("git", "add", ".")
			err = cmd.Run()
			if err != nil {
				fmt.Println("error stashing changes:", err)
				return
			}

			// Commit changes
			cmd = exec.Command("git", "commit", "-m", conf.commitMsg)
			err = cmd.Run()
			if err != nil {
				fmt.Println("error committing changes:", err)
				return
			}

			// Push changes
			cmd = exec.Command("git", "push", "origin", conf.featureBranch)
			err = cmd.Run()
			if err != nil {
				fmt.Println("error pushing changes:", err)
				return
			}

			fmt.Println("Changes committed and pushed to feature branch.")

			// create PR
			newPR := &github.NewPullRequest{
				Title:               github.String(conf.prTitle),
				Head:                github.String(conf.featureBranch),
				Base:                github.String(conf.targetBranch),
				Body:                github.String(conf.prBody),
				MaintainerCanModify: github.Bool(true),
			}

			pr, _, err := client.PullRequests.Create(ctx, conf.organization, repo, newPR)
			if err != nil {
				fmt.Println(err)
				return
			}
			prLinks = append(prLinks, pr.GetHTMLURL())

			reviewers := github.ReviewersRequest{
				//NodeID:    github.String(*pr.NodeID),
				Reviewers: conf.prReviewer,
				//TeamReviewers: []string{""},
			}

			_, _, err = client.PullRequests.RequestReviewers(ctx, conf.organization, repo, *pr.Number, reviewers)
			if err != nil {
				fmt.Println("Error adding reviewer:", err)
				return
			}
		} // end for loop
		fmt.Println("===========================================================")
		fmt.Println("Here are the PRs link...")
		printArray(prLinks)

	} else {

		fmt.Println("Input was not 'ok', deleting new feature branches.")

		// todo   if input not ok or something goes wrong ,should remove the feature branchs

	}

}
