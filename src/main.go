package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	git "github.com/wwleak/tidy/local"
)

func main() {

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	var err error
	var local *git.Local
	repository, err := local.Init()
	if err != nil {
		panic(err)
	}

	instance := git.Local{Repository: repository}

	// repository.GetMergedBranches()

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

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		fmt.Println("Unable to create box:", err)
	}

	listbox, err := gtk.ListBoxNew()
	if err != nil {
		fmt.Println("Unable to create listbox:", err)
	}

	branches := []string{}
	branches, _ = instance.GetMergedBranches()
	fmt.Println(len(branches))
	//for _, name := range branches {
	row, _ := gtk.ListBoxRowNew()
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	row.Add(hbox)
	item, _ := gtk.LabelNew("Item")
	hbox.PackStart(item, true, true, 0)
	listbox.Add(row)
	box.PackStart(listbox, true, true, 0)

	btn, err := gtk.ButtonNewWithLabel("Delete")
	if err != nil {
		panic(err)
	}
	btn.Connect("clicked", nil)
	box.Add(btn)

	win.Add(box)
	//}

	//box.Add(listbox)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()

}
