package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/src-d/go-git.v4"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"

	"os"
	"os/user"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	. "gopkg.in/src-d/go-git.v4/_examples"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	fmt.Println("Welcome to tidy")

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	r, err := git.PlainOpen(usr.HomeDir + "/shitchat")
	CheckIfError(err)

	// func (r *Repository) References() (storer.ReferenceIter, error)
	references, err := r.References()
	CheckIfError(err)

	Info("Checkout to master branch")

	branch := fmt.Sprintf("refs/heads/%s", os.Getenv("BRANCH"))
	fmt.Println(branch)
	b := plumbing.ReferenceName(branch)

	wt, err := r.Worktree()
	CheckIfError(err)
	err = wt.Checkout(&git.CheckoutOptions{Create: false, Force: false, Branch: b})
	CheckIfError(err)

	Info("git branch")

	// NewReferenceFilteredIter returns a reference iterator for the given reference Iterator.
	// This iterator will iterate only references that accomplish the provided function.
	iter := storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return strings.Contains(ref.Name().String(), "/heads/") && ref.Name().Short() != "master"
	}, references)

	Info("Get master commits")

	masterheads := make(map[string]plumbing.Hash)
	commits, err := r.Log(&git.LogOptions{From: masterheads["master"]})
	CheckIfError(err)

	Info("Get branches merged to master")

	branches := make(map[string]plumbing.Hash)

	err = iter.ForEach(func(branch *plumbing.Reference) error {
		// The last commit of every branch is stored in branches slice
		name := branch.Name().Short()
		head := branch.Hash()
		branches[name] = head

		return nil
	})
	CheckIfError(err)

	merged := make([]string, 0)
	err = commits.ForEach(func(commit *object.Commit) error {
		for name, head := range branches {
			if head.String() == commit.Hash.String() {
				fmt.Printf("Branch %s head (%s) was found in master, so has been merged!\n", name, head)
				merged = append(merged, name)
			}
		}
		return nil
	})

	// Delete repository branch
	/* for _, branch := range merged {
		fmt.Println("Deleting ", branch, " branch")
		err := r.DeleteBranch(branch)
		if err != nil {
			panic(err)
		}
	} */

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Unable to create window:", err)
	}
	win.SetTitle("Tidy")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new label widget to show in the window.
	/* l, err := gtk.LabelNew("Keep your local git repo clean once and for all!")
	if err != nil {
		fmt.Println("Unable to create label:", err)
	}

	// Add the label to the window.
	win.Add(l) */

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		fmt.Println("Unable to create box:", err)
	}

	listbox, err := gtk.ListBoxNew()
	if err != nil {
		fmt.Println("Unable to create listbox:", err)
	}
	row, _ := gtk.ListBoxRowNew()
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	row.Add(hbox)
	item, _ := gtk.LabelNew("Item")
	hbox.PackStart(item, true, true, 0)
	listbox.Add(row)
	box.PackStart(listbox, true, true, 0)
	win.Add(box)
	//box.Add(listbox)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()

}
