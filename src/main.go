package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	"encoding/json"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	git "github.com/wwleak/tidy/local"

	preferences "github.com/wwleak/tidy/settings"
	"github.com/wwleak/tidy/window"
)

const (
	COLUMN_NAME = iota
)

var (
	s *preferences.Settings
)

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

	win := window.New("Tidy")

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	css, _ := gtk.CssProviderNew()
	css.LoadFromPath("styles.css")

	screen, _ := gdk.ScreenGetDefault()

	gtk.AddProviderForScreen(screen, css, gtk.STYLE_PROVIDER_PRIORITY_USER)

	menubar := window.SetMenuBar()
	menuitem := window.SetMenuItem("File")
	filemenu := window.SetMenuNew()
	settings := window.SetMenuItem("Settings")
	close := window.SetMenuItem("Close")

	close.Connect("activate", func() {
		win.Close()
	})

	settings.Connect("activate", func() {
		settings := window.New("Settings")

		settings.Connect("destroy", func() {
			settings.Close()
		})

		dirlabel := window.SetLabel("Repository folder")
		dirlabel.SetName("folder")
		dirlabel.SetHAlign(gtk.ALIGN_START)

		directory := window.SetFolderChooserButton("Repository")

		reflabel := window.SetLabel("Branch's name")
		reflabel.SetName("branch")
		reflabel.SetHAlign(gtk.ALIGN_START)

		box := window.SetBox(gtk.ORIENTATION_VERTICAL, 5)
		box.SetName("config")

		branch := window.SetEntry()

		box.PackStart(dirlabel, false, false, 0)
		box.PackStart(directory, false, false, 0)
		box.PackStart(reflabel, false, false, 0)
		box.PackStart(branch, false, false, 0)

		btn := window.SetButton("Save")
		btn.Connect("clicked", func() {
			branch, _ := branch.GetText()

			config := preferences.Settings{
				Repository: preferences.Repository{
					Branch: branch,
					Folder: directory.GetFilename(),
				},
			}
			config.Save("/tmp/tidy.yaml")
		})

		box.PackEnd(btn, false, true, 5)

		settings.Add(box)

		settings.ShowAll()
	})

	filemenu.Append(settings)
	filemenu.Append(close)
	menuitem.SetSubmenu(filemenu)
	menubar.Append(menuitem)

	box := window.SetBox(gtk.ORIENTATION_VERTICAL, 0)

	frame := window.SetFrame()
	frame.SetName("frame")

	tree, store := window.SetTreeView("Merged Branches")
	// tree.SetName("tree")
	frame.Add(tree)

	branches := []string{}
	branches, _ = instance.GetMergedBranches()

	for _, name := range branches {
		window.AddTreeViewRow(store, name)
	}

	box.PackStart(menubar, false, true, 0)
	box.PackStart(frame, true, true, 0)

	hbox := window.SetBox(gtk.ORIENTATION_HORIZONTAL, 5)
	hbox.SetName("tools")

	search := window.SetButton("Search")

	hbox.PackStart(search, false, true, 0)

	delete := window.SetButton("Delete")

	delete.Connect("clicked", func() {
		instance.DeleteBranches(branches)
		store.Clear()
	})

	hbox.PackEnd(delete, true, true, 0)

	box.PackEnd(hbox, false, false, 0)

	win.Add(box)

	win.ShowAll()

	gtk.Main()

}
