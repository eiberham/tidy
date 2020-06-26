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

func setupWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		panic(err)
	}

	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetDefaultSize(300, 400)
	win.SetResizable(false)

	return win
}

func add(store *gtk.ListStore, text string) {
	iter := store.Append()

	err := store.Set(iter,
		[]int{COLUMN_NAME},
		[]interface{}{text})

	if err != nil {
		panic(err)
	}

	return
}

func main() {
	gtk.Init(nil)

	var err error
	var local *git.Local
	repository, err := local.Init()
	if err != nil {
		panic(err)
	}

	instance := git.Local{Repository: repository}

	win := setupWindow("Tidy")

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		panic(err)
	}

	tree, err := gtk.TreeViewNew()
	if err != nil {
		panic(err)
	}

	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		panic(err)
	}

	text := fmt.Sprintf("Merged Branches")
	column, err := gtk.TreeViewColumnNewWithAttribute(text, cell, "text", COLUMN_NAME)
	// column.AddAttribute()

	if err != nil {
		panic(err)
	}

	tree.AppendColumn(column)

	store, err := gtk.ListStoreNew(glib.TYPE_STRING)

	tree.SetModel(store)

	branches := []string{}
	branches, _ = instance.GetMergedBranches()

	for _, name := range branches {
		add(store, name)
	}

	box.PackStart(tree, true, true, 0)

	btn, err := gtk.ButtonNewWithLabel("Delete")
	if err != nil {
		panic(err)
	}

	btn.Connect("clicked", func() {
		instance.DeleteBranches(branches)
		store.Clear()
	})

	box.PackEnd(btn, false, false, 0)

	win.Add(box)

	win.ShowAll()

	gtk.Main()

}
