package repository

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

var (
	ErrRetrievingHead       = errors.New("Couldn't bring heads from target branch")
	ErrNotFound             = errors.New("Repository not found in provided path")
	ErrDeleteBranch         = errors.New("Sorry, couldn't delete branch")
	ErrGettingLocalBranches = errors.New("Something bad happened getting local branches")
)

type Repository struct {
	Self *git.Repository
}

// Init returns a repository instance or an error otherwise
func (repository *Repository) Init(folder string) (*git.Repository, error) {
	instance, err := git.PlainOpen(folder)

	if err != nil {
		return nil, ErrNotFound
	}

	return instance, nil
}

// GetMergedBranches returns all those branches that have been merged
func (repository *Repository) GetMergedBranches(branch string) ([]string, error) {
	err := repository.checkout(branch)
	if err != nil {
		panic(err)
	}
	references := repository.GetLocalBranches(branch)

	targetheads := make(map[string]plumbing.Hash)
	commits, err := repository.Self.Log(&git.LogOptions{From: targetheads[branch]})
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
func (repository *Repository) checkout(branch string) error {
	target := fmt.Sprintf("refs/heads/%s", branch)
	b := plumbing.ReferenceName(target)

	wt, err := repository.Self.Worktree()
	if err != nil {
		panic(err)
	}

	return wt.Checkout(&git.CheckoutOptions{Create: false, Force: false, Branch: b})
}

// GetLocalBranches returns a ReferenceIter of all filtered branches
func (repository *Repository) GetLocalBranches(exclude string) storer.ReferenceIter {
	references, err := repository.Self.References()
	if err != nil {
		panic(err)
	}

	iter := storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return strings.Contains(ref.Name().String(), "/heads/") && ref.Name().Short() != "master" && ref.Name().Short() != exclude
	}, references)

	return iter
}

// DeleteBranches deletes all those merged branches from the local repository
func (repository *Repository) DeleteBranches(branches []string) (bool, error) {
	for _, name := range branches {
		branch := fmt.Sprintf("refs/heads/%s", name)
		b := plumbing.ReferenceName(branch)

		err := repository.Self.Storer.RemoveReference(b)
		if err != nil {
			return false, ErrDeleteBranch
		}
	}
	return true, nil
}

func (repository *Repository) GetBranches() []string {
	branches, _ := repository.Self.Branches()
	result := make([]string, 0)
	branches.ForEach(func(branch *plumbing.Reference) error {
		result = append(result, branch.Name().Short())
		return nil
	})

	return result
}
