package window

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	COLUMN_NAME = iota
)

// ...
func New(title string) *gtk.Window {
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

// ...
func SetButton(label string) *gtk.Button {
	button, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		panic(err)
	}

	return button
}

// ...
func SetLabel(text string) *gtk.Label {
	label, err := gtk.LabelNew(text)
	if err != nil {
		panic(err)
	}
	return label
}

// ...
func SetBox(orientation gtk.Orientation) *gtk.Box {
	box, err := gtk.BoxNew(orientation, 0)
	if err != nil {
		panic(err)
	}
	return box
}

// ...
func SetTreeView(caption string) (*gtk.TreeView, *gtk.ListStore) {
	tree, err := gtk.TreeViewNew()
	if err != nil {
		panic(err)
	}

	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		panic(err)
	}

	title := fmt.Sprintf(caption)
	var column *gtk.TreeViewColumn
	column, err = gtk.TreeViewColumnNewWithAttribute(title, cell, "text", COLUMN_NAME)

	if err != nil {
		panic(err)
	}

	tree.AppendColumn(column)
	store, err := gtk.ListStoreNew(glib.TYPE_STRING)
	tree.SetModel(store)

	return tree, store
}

// ...
func AddTreeViewRow(store *gtk.ListStore, text string) {
	iter := store.Append()

	err := store.Set(iter,
		[]int{COLUMN_NAME},
		[]interface{}{text})

	if err != nil {
		panic(err)
	}

	return
}
