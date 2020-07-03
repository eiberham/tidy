package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	"encoding/json"
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	git "github.com/wwleak/tidy/local"
	"github.com/wwleak/tidy/settings"
)

const (
	COLUMN_NAME = iota
)

var (
	s *settings.Settings
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
	var err error

	config, err := s.Open("/tmp/tidy.yaml")
	if err != nil {
		fmt.Println("Something bad happened")
	}

	fmt.Printf("Result: %v\n", config)

	data, err := json.Marshal(config)
	fmt.Printf("%s\n", data)

	gtk.Init(nil)

	var local *git.Local
	repository, err := local.Init()
	if err != nil {
		panic(err)
	}

	instance := git.Local{Repository: repository}

	win := setupWindow("Tidy")

	// Menu bar

	menubar, err := gtk.MenuBarNew()

	menuitem, err := gtk.MenuItemNewWithLabel("File")

	filemenu, err := gtk.MenuNew()

	fileitem, err := gtk.MenuItemNewWithLabel("Settings")

	closeitem, err := gtk.MenuItemNewWithLabel("Close")

	closeitem.Connect("activate", func() {
		win.Close()
	})

	fileitem.Connect("activate", func() {
		sett, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
		if err != nil {
			panic(err)
		}

		sett.SetTitle("Settings")
		sett.Connect("destroy", func() {
			sett.Close()
		})
		sett.SetPosition(gtk.WIN_POS_CENTER)
		sett.SetDefaultSize(300, 300)
		sett.SetResizable(false)

		trgt, err := gtk.LabelNew("Repository Folder")
		if err != nil {
			panic(err)
		}
		trgt.SetMarginTop(5)
		trgt.SetMarginStart(5)
		trgt.SetHAlign(gtk.ALIGN_START)

		repo, err := gtk.FileChooserButtonNew("Repository", gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)
		if err != nil {
			panic(err)
		}
		repo.SetMarginTop(5)
		repo.SetMarginStart(5)
		repo.SetMarginEnd(5)

		repo.Connect("selection-changed", func() {
			folder := repo.GetFilename()
			fmt.Printf("folder: %s ", folder)
		})

		brch, err := gtk.LabelNew("Branch's Name")
		if err != nil {
			panic(err)
		}
		brch.SetMarginTop(5)
		brch.SetMarginStart(5)
		brch.SetHAlign(gtk.ALIGN_START)

		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		if err != nil {
			panic(err)
		}

		entry, err := gtk.EntryNew()
		if err != nil {
			panic(err)
		}
		entry.SetMarginTop(5)
		entry.SetMarginStart(5)
		entry.SetMarginEnd(5)

		box.Add(trgt)
		box.Add(repo)
		box.Add(brch)
		box.Add(entry)

		/* grid, err := gtk.GridNew()
		if err != nil {
			panic(err)
		}


		grid.SetMarginStart(20)
		grid.SetMarginEnd(20)
		grid.SetMarginTop(20)
		grid.SetMarginBottom(20)
		grid.SetRowSpacing(20)
		grid.SetColumnSpacing(20)
		grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

		grid.Attach(trgt, 0, 1, 1, 1)
		grid.Attach(repo, 1, 20, 1, 1)

		trgt.SetHExpand(true)
		repo.SetHExpand(true) */

		btn, err := gtk.ButtonNewWithLabel("Save")
		if err != nil {
			panic(err)
		}
		btn.SetMarginStart(5)
		btn.SetMarginEnd(5)
		btn.Connect("clicked", func() {
			fmt.Println("clicked save btn")
		})

		/* grid.Attach(btn, 0, 40, 1, 1)
		sett.Add(grid) */

		box.PackEnd(btn, false, true, 5)

		sett.Add(box)

		sett.ShowAll()
	})

	filemenu.Append(fileitem)
	filemenu.Append(closeitem)

	menuitem.SetSubmenu(filemenu)

	menubar.Append(menuitem)

	// End

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

	box.PackStart(menubar, false, true, 0)
	box.PackStart(tree, true, true, 0)

	srh, err := gtk.ButtonNewWithLabel("Search")
	if err != nil {
		panic(err)
	}
	srh.SetMarginStart(5)
	srh.SetMarginEnd(5)
	srh.SetMarginTop(5)
	srh.SetMarginBottom(5)

	box.PackEnd(srh, false, false, 0)

	btn, err := gtk.ButtonNewWithLabel("Delete")
	if err != nil {
		panic(err)
	}
	btn.SetMarginStart(5)
	btn.SetMarginEnd(5)
	btn.SetMarginTop(5)
	btn.SetMarginBottom(5)

	btn.Connect("clicked", func() {
		instance.DeleteBranches(branches)
		store.Clear()
	})

	box.PackEnd(btn, false, false, 0)

	win.Add(box)

	win.ShowAll()

	gtk.Main()

}
