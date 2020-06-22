package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
)

var (
	local Local
)

func main() {

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
