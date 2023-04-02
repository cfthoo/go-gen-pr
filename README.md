This  script is useful when you have to create PRs for multiple repos.
It will automate the process of create PRs.

What it does is it will execute the git command in go.
The git command here are
-git fetch origin
-git checkout
-git add
-git commit
-git push
-add PR

There are two functions for different scenario.

Let's say if you have 3 env , dev / stage /prod.
First function is PR for dev.
Second function is to rebase dev to stage.

Simply uncomment "devToStage" or "devToDev".
Then fill in the PR detail in the conf struct and run "go run main.go" 
