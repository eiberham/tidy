package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	git "github.com/wwleak/tidy/local"
)

func main() {
	gtk.Init(nil)

	var err error
	var local *git.Local
	repository, err := local.Init()
	if err != nil {
		panic(err)
	}

	instance := git.Local{Repository: repository}

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

	for _, name := range branches {
		row, _ := gtk.ListBoxRowNew()
		hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		row.Add(hbox)
		item, _ := gtk.LabelNew(name)
		hbox.PackStart(item, true, true, 0)
		listbox.Add(row)
	}

	box.PackStart(listbox, true, true, 0)

	btn, err := gtk.ButtonNewWithLabel("Delete")
	if err != nil {
		panic(err)
	}
	btn.Connect("clicked", func() {
		instance.DeleteBranches(branches)
	})
	box.Add(btn)

	win.Add(box)

	win.SetDefaultSize(500, 400)

	win.ShowAll()

	gtk.Main()

}
