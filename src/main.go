package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	git "github.com/wwleak/tidy/local"
)

const (
	COLUMN_NAME = iota
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
	win.SetResizable(false)
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

	tree, err := gtk.TreeViewNew()
	if err != nil {
		fmt.Println("Unable to create treeview:", err)
	}

	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		fmt.Println("Unable to create cell:", err)
	}

	text := fmt.Sprintf("Merged Branches")
	column, err := gtk.TreeViewColumnNewWithAttribute(text, cell, "text", COLUMN_NAME)
	// column.AddAttribute()

	if err != nil {
		fmt.Println("Unable to create column:", err)
	}

	tree.AppendColumn(column)

	store, err := gtk.ListStoreNew(glib.TYPE_STRING)

	tree.SetModel(store)

	branches := []string{}
	branches, _ = instance.GetMergedBranches()

	for _, name := range branches {
		addRow(store, name)
	}

	box.PackStart(tree, true, true, 0)

	/* for _, name := range branches {
		row, _ := gtk.ListBoxRowNew()
		hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		row.Add(hbox)
		item, _ := gtk.LabelNew(name)
		hbox.PackStart(item, true, true, 0)
		list.Add(row)
	}

	box.PackStart(list, true, true, 0) */

	btn, err := gtk.ButtonNewWithLabel("Delete")
	if err != nil {
		panic(err)
	}

	btn.Connect("clicked", remove(store, git.Local{Repository: repository}, branches))

	box.PackEnd(btn, false, false, 0)

	win.Add(box)

	win.SetDefaultSize(300, 400)

	win.ShowAll()

	gtk.Main()

}

func addRow(store *gtk.ListStore, text string) {
	iter := store.Append()

	err := store.Set(iter,
		[]int{COLUMN_NAME},
		[]interface{}{text})

	if err != nil {
		fmt.Println("Unable to add element to:", err)
	}
}

func remove(store *gtk.ListStore, instance *git.Local, branches []string) error {
	instance.DeleteBranches(branches)
	store.Clear()
	return nil
}
