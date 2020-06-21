package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

var (
	repository *git.Repository
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	repository, err := git.PlainOpen(os.Getenv("PATH"))
	if err != nil {
		panic(err)
	}

}

func GetMergedBranches() []string {
	checkoutToTarget()
	local := getLocalBranches()

	targetheads := make(map[string]plumbing.Hash)
	commits, err := repository.Log(&git.LogOptions{From: targetheads[os.Getenv("BRANCH")]})
	if err != nil {
		panic(err)
	}

	branches := make(map[string]plumbing.Hash)

	err = local.ForEach(func(branch *plumbing.Reference) error {
		// The last commit of every branch is stored in branches slice
		name := branch.Name().Short()
		head := branch.Hash()
		branches[name] = head

		return nil
	})
	if err != nil {
		panic(err)
	}

	merged := make([]string, 0)
	err = commits.ForEach(func(commit *object.Commit) error {
		for name, head := range branches {
			if head.String() == commit.Hash.String() {
				merged = append(merged, name)
			}
		}
		return nil
	})

	return merged
}

/**

 */
func checkoutToTarget() error {
	branch := fmt.Sprintf("refs/heads/%s", os.Getenv("BRANCH"))
	fmt.Println(branch)
	b := plumbing.ReferenceName(branch)

	wt, err := repository.Worktree()
	if err != nil {
		panic(err)
	}

	err = wt.Checkout(&git.CheckoutOptions{Create: false, Force: false, Branch: b})
	if err != nil {
		panic(err)
	}
}

/**

 */
func getLocalBranches() storer.ReferenceIter {
	references, err := repository.References()
	if err != nil {
		panic(err)
	}

	iter := storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return strings.Contains(ref.Name().String(), "/heads/") && ref.Name().Short() != "master" && ref.Name().Short() != os.Getenv("BRANCH")
	}, references)

	return iter
}

func DeleteBranches(branches []string) bool {
	for _, branch := range branches {
		err := repository.DeleteBranch(branch)
		if err != nil {
			panic(err)
		}
	}
	return true
}
