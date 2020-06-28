package local

import (
	"errors"
	"fmt"
	"os/user"
	"sort"
	"strings"

	"github.com/wwleak/tidy/config"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

var (
	ErrLoadingEnvFile       = errors.New("Couldn't load the .env file")
	ErrRetrievingHead       = errors.New("Couldn't bring heads from target branch")
	ErrNotFound             = errors.New("Repository not found in provided path")
	ErrDeleteBranch         = errors.New("Sorry, couldn't delete branch")
	ErrGettingLocalBranches = errors.New("Something bad happened getting local branches")
)

type Local struct {
	Repository *git.Repository
}

// Init returns a repository instance or an error otherwise
func (local *Local) Init() (*git.Repository, error) {
	err := config.Load()
	if err != nil {
		return nil, ErrLoadingEnvFile
	}

	usr, _ := user.Current()
	repository, err := git.PlainOpen(usr.HomeDir + "/" + config.Get("TARGET"))

	if err != nil {
		return nil, ErrNotFound
	}

	return repository, nil
}

// GetMergedBranches returns all those branches that have been merged
func (local *Local) GetMergedBranches() ([]string, error) {
	err := local.checkoutToTarget()
	if err != nil {
		panic(err)
	}
	references := local.getLocalBranches()

	targetheads := make(map[string]plumbing.Hash)
	commits, err := local.Repository.Log(&git.LogOptions{From: targetheads[config.Get("BRANCH")]})
	if err != nil {
		return nil, ErrRetrievingHead
	}

	branches := make(map[string]plumbing.Hash)

	err = references.ForEach(func(branch *plumbing.Reference) error {
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

	sort.Strings(merged)

	return merged, nil
}

// checkoutToTarget switches to the target branch or throws error
func (local *Local) checkoutToTarget() error {
	branch := fmt.Sprintf("refs/heads/%s", config.Get("BRANCH"))
	fmt.Println(branch)
	b := plumbing.ReferenceName(branch)

	wt, err := local.Repository.Worktree()
	if err != nil {
		panic(err)
	}

	return wt.Checkout(&git.CheckoutOptions{Create: false, Force: false, Branch: b})
}

// getLocalBranches returns a ReferenceIter of all filtered branches
func (local *Local) getLocalBranches() storer.ReferenceIter {
	references, err := local.Repository.References()
	if err != nil {
		panic(err)
	}

	iter := storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return strings.Contains(ref.Name().String(), "/heads/") && ref.Name().Short() != "master" && ref.Name().Short() != config.Get("BRANCH")
	}, references)

	return iter
}

// DeleteBranches deletes all those merged branches from the local repository
func (local *Local) DeleteBranches(branches []string) (bool, error) {
	for _, name := range branches {
		branch := fmt.Sprintf("refs/heads/%s", name)
		b := plumbing.ReferenceName(branch)

		err := local.Repository.Storer.RemoveReference(b)
		if err != nil {
			return false, ErrDeleteBranch
		}
		break
	}
	return true, nil
}
