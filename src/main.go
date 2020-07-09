package main

// git for-each-ref --format='%(committerdate) %09 %(authorname) %09 %(refname)' --sort=committerdate

import (
	// "encoding/json"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	git "github.com/wwleak/tidy/repository"

	preferences "github.com/wwleak/tidy/settings"
	"github.com/wwleak/tidy/window"
)

var (
	config     *preferences.Settings
	repository *git.Repository
)

func configured() bool {
	settings, err := config.Open("/tmp/tidy.yaml")
	config = settings

	if err != nil {
		return false
	}

	return true
}

func load(folder *gtk.FileChooserButton, branch *gtk.ComboBoxText) {
	folder.SetCurrentFolder(config.Repository.Folder)
	branch.SetSensitive(true)
	self, err := repository.Init(config.Repository.Folder)
	if err != nil {
		panic(err)
	}

	repository = &git.Repository{Self: self}
	branches := repository.GetBranches()

	var selItem int
	for key, item := range branches[:] {
		if item == config.Repository.Branch {
			selItem = key
		}
		branch.AppendText(item)
	}

	branch.SetActive(selItem)

	return
}

func main() {

	gtk.Init(nil)

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

		branch := window.SetComboBox()
		branch.SetSensitive(false)

		directory.Connect("selection-changed", func() {
			folder := directory.GetFilename()
			branch.SetSensitive(true)

			self, err := repository.Init(folder)
			if err != nil {
				panic(err)
			}

			repository = &git.Repository{Self: self}
			branches := repository.GetBranches()
			fmt.Println(branches)

			for _, item := range branches {
				branch.AppendText(item)
			}
		})

		/* grid := window.SetGrid()

		rmtlabel := window.SetLabel("Enable Remote ?")

		toggle := window.SetSwitch()

		grid.Attach(rmtlabel, 0, 100, 250, 100)
		grid.Attach(toggle, 250, 100, 250, 100) */

		box.PackStart(dirlabel, false, false, 0)
		box.PackStart(directory, false, false, 0)
		box.PackStart(reflabel, false, false, 0)
		box.PackStart(branch, false, false, 0)
		// box.PackStart(toggle, false, false, 0)
		// box.PackStart(grid, false, false, 0)

		save := window.SetButton("Save")
		save.SetName("save")
		save.Connect("clicked", func() {
			// branch, _ := branch.GetText()
			option, _ := branch.GetEntry()
			branch, _ := option.GetText()

			config := preferences.Settings{
				Repository: preferences.Repository{
					Branch: branch,
					Folder: directory.GetFilename(),
				},
			}
			if success, _ := config.Save("/tmp/tidy.yaml"); success {
				window.ShowDialogMessage(win, gtk.MESSAGE_INFO, "Done", "Configuration set successfully")
			}

		})

		box.PackEnd(save, false, true, 5)

		settings.Add(box)

		settings.ShowAll()

		if configured() {
			load(directory, branch)
		}
	})

	filemenu.Append(settings)
	filemenu.Append(close)
	menuitem.SetSubmenu(filemenu)
	menubar.Append(menuitem)

	box := window.SetBox(gtk.ORIENTATION_VERTICAL, 0)

	scrolled := window.SetScrolledWindow()

	frame := window.SetFrame()
	frame.SetName("frame")

	tree, store := window.SetTreeView("Merged Branches")

	scrolled.Add(tree)

	frame.Add(scrolled)

	box.PackStart(menubar, false, true, 0)
	box.PackStart(frame, true, true, 0)

	hbox := window.SetBox(gtk.ORIENTATION_HORIZONTAL, 5)
	hbox.SetName("tools")

	search := window.SetButton("Search")

	search.Connect("clicked", func() {
		if configured() {
			self, err := repository.Init(config.Repository.Folder)
			if err != nil {
				panic(err)
			}

			repository = &git.Repository{Self: self}
			branches := []string{}
			branches, _ = repository.GetMergedBranches(config.Repository.Branch)
			store.Clear()
			for _, name := range branches {
				window.AddTreeViewRow(store, name)
			}

			return
		}
		window.ShowDialogMessage(win, gtk.MESSAGE_WARNING, "Oops", "You haven't set any configuration!")
	})

	hbox.PackStart(search, false, true, 0)

	delete := window.SetButton("Delete")
	delete.SetName("delete")

	delete.Connect("clicked", func() {
		branches := []string{}
		branches, _ = repository.GetMergedBranches(config.Repository.Branch)
		repository.DeleteBranches(branches)
		store.Clear()
	})

	hbox.PackEnd(delete, true, true, 0)

	box.PackEnd(hbox, false, false, 0)

	win.Add(box)

	win.ShowAll()

	gtk.Main()

}
