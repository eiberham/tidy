package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	git "github.com/wwleak/tidy/local"
	// "github.com/wwleak/tidy/settings"
)

const (
	COLUMN_NAME = iota
)

/*var (
	s *settings.Settings
)*/

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

	/*config, err := s.Open("/tmp/settings.yaml")
	if err != nil {
		fmt.Println("Something bad happened")
	}

	fmt.Println("Done, %T", config)*/

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
		sett.SetDefaultSize(600, 300)
		sett.SetResizable(false)

		trgt, err := gtk.LabelNew("Choose your repository")
		if err != nil {
			panic(err)
		}

		repo, err := gtk.FileChooserButtonNew("Repository", gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)
		if err != nil {
			panic(err)
		}

		grid, err := gtk.GridNew()
		if err != nil {
			panic(err)
		}

		grid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)

		grid.Add(trgt)
		grid.Add(repo)

		trgt.SetHExpand(true)
		repo.SetHExpand(true)

		/* rama, err := gtk.LabelNew("What's your branch's name ?")
		if err != nil {
			panic(err)
		}

		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		if err != nil {
			panic(err)
		} */

		/* box.Add(rama)
		grid.Add(box) */

		sett.Add(grid)

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
