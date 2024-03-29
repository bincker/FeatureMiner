package main

import (
	"fmt"
	"github.com/AndroSOM/FeatureMiner/helper"
	"github.com/AndroSOM/FeatureMiner/miner"
	"github.com/AndroSOM/FeatureMiner/setup"
	"github.com/mattn/go-gtk/gtk"
	"runtime"
	"strconv"
)

var inputFolder, outputFolder string
var apks []string

func CreateWindow() *gtk.Window {
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("AndroSOM Feature Miner")
	window.SetDefaultSize(600, 600)
	vbox := gtk.NewVBox(false, 1)
	CreateMenu(window, vbox)

	notebook := gtk.NewNotebook()
	main_page := gtk.NewFrame("Miner")
	main_page.Add(MinerPage())
	notebook.AppendPage(main_page, gtk.NewLabel("Feature Miner"))
	setup_page := gtk.NewFrame("Setup")
	setup_page.Add(SetupPage())
	notebook.AppendPage(setup_page, gtk.NewLabel("Dependencies"))

	vbox.Add(notebook)
	window.Add(vbox)
	return window
}

func CreateMenu(w *gtk.Window, vbox *gtk.VBox) {
	action_group := gtk.NewActionGroup("my_group")
	ui_manager := CreateUIManager()
	accel_group := ui_manager.GetAccelGroup()
	w.AddAccelGroup(accel_group)

	action_group.AddAction(gtk.NewAction("FileMenu", "File", "", ""))

	action_filequit := gtk.NewAction("FileQuit", "", "", gtk.STOCK_QUIT)
	action_filequit.Connect("activate", gtk.MainQuit)
	action_group.AddActionWithAccel(action_filequit, "")

	action_group.AddAction(gtk.NewAction("HelpMenu", "Help", "", ""))

	action_help_about := gtk.NewAction("HelpAbout", "About", "", "")
	action_help_about.Connect("activate", func() {
		dialog := gtk.NewAboutDialog()
		dialog.SetProgramName("FeatureMiner")
		dialog.SetComments("FeatureMiner is part of the AndroSOM project which was built at the university of applied sciences Salzburg as part of a master's thesis.")
		dialog.SetAuthors(helper.Authors())
		dialog.SetLicense("Copyright (c) 2014 AndroSOM\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of button software and associated documentation files (the \"Software\"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and button permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.")
		dialog.SetWrapLicense(true)
		dialog.Run()
		dialog.Destroy()
	})
	action_group.AddActionWithAccel(action_help_about, "")

	ui_manager.InsertActionGroup(action_group, 0)
	menubar := ui_manager.GetWidget("/MenuBar")
	vbox.PackStart(menubar, false, false, 0)
	eventbox := gtk.NewEventBox()
	vbox.PackStart(eventbox, false, false, 0)
}

func MinerPage() *gtk.VBox {
	vbox := gtk.NewVBox(false, 1)

	framebox := gtk.NewVBox(false, 1)

	framebox.Add(helper.SetFolder("Input", &inputFolder))
	framebox.Add(helper.SetFolder("Output", &outputFolder))

	load_button_hbox := gtk.NewHBox(false, 0)
	load_button_hbox.SetBorderWidth(5)
	load_button_hbox.SetSizeRequest(-1, 60)
	load_button := gtk.NewButtonWithLabel("Load Android APKs")
	load_button_hbox.PackStart(load_button, true, true, 0)
	framebox.Add(load_button_hbox)

	apk_count_label_hbox := gtk.NewHBox(false, 0)
	apk_count_label_hbox.SetBorderWidth(5)
	apk_count_label_hbox.SetSizeRequest(-1, 30)
	apk_count_label := gtk.NewLabel("Loaded APKs: None.. Please load APKs to continue.")
	apk_count_label_hbox.PackStart(apk_count_label, true, true, 0)
	framebox.Add(apk_count_label_hbox)

	vbox.PackStart(framebox, false, true, 0)

	/* Static Analysis */

	static_analysis_frame := gtk.NewFrame("Static Analysis")
	static_analysis_frame.SetBorderWidth(5)
	static_analysis_frame.SetSensitive(false)

	static_analysis_hbox := gtk.NewHBox(false, 0)
	static_analysis_hbox.SetBorderWidth(10)
	static_analysis_hbox.SetSizeRequest(-1, 60)

	static_analysis_progress := gtk.NewProgressBar()
	static_analysis_start_button := gtk.NewButtonWithLabel("Start Analysis")

	static_analysis_hbox.PackStart(static_analysis_start_button, false, true, 5)
	static_analysis_hbox.PackStart(static_analysis_progress, true, true, 5)
	static_analysis_frame.Add(static_analysis_hbox)
	vbox.PackStart(static_analysis_frame, false, true, 0)

	/* Dynamic Analyisi */

	dynamic_analysis_frame := gtk.NewFrame("Dynamic Analysis")
	dynamic_analysis_frame.SetBorderWidth(5)
	dynamic_analysis_frame.SetSensitive(false)

	dynamic_analysis_hbox := gtk.NewHBox(false, 0)
	dynamic_analysis_hbox.SetBorderWidth(10)
	dynamic_analysis_hbox.SetSizeRequest(-1, 60)

	dynamic_analysis_progress := gtk.NewProgressBar()
	dynamic_analysis_start_button := gtk.NewButtonWithLabel("Start Analysis")

	dynamic_analysis_hbox.PackStart(dynamic_analysis_start_button, false, true, 5)
	dynamic_analysis_hbox.PackStart(dynamic_analysis_progress, true, true, 5)
	dynamic_analysis_frame.Add(dynamic_analysis_hbox)
	vbox.PackStart(dynamic_analysis_frame, false, true, 0)

	/* Events */

	load_button.Connect("clicked", func() {
		if helper.FolderExists(inputFolder, "Please provide a valid input directory!") && helper.FolderExists(outputFolder, "Please provide a valid output directory!") {
			miner.LoadAPKs(inputFolder, &apks, load_button)
			if len(apks) > 0 {
				apk_count_label.SetLabel("Loaded APKs: " + strconv.Itoa(len(apks)))
				static_analysis_frame.SetSensitive(true)
				dynamic_analysis_frame.SetSensitive(true)
			} else {
				apk_count_label.SetLabel("No APKs found in this input folder..")
				static_analysis_frame.SetSensitive(false)
				dynamic_analysis_frame.SetSensitive(false)
			}
		}
	})

	static_analysis_start_button.Connect("clicked", func() {
		fmt.Println("starting static analysis..")
		static_analysis_start_button.SetSensitive(false)
		miner.Analysis(&apks, outputFolder, static_analysis_progress, "static_analysis.py", runtime.NumCPU())
		static_analysis_start_button.SetSensitive(true)
	})

	dynamic_analysis_start_button.Connect("clicked", func() {
		fmt.Println("starting dynamic analysis..")
		dynamic_analysis_start_button.SetSensitive(false)
		miner.Analysis(&apks, outputFolder, dynamic_analysis_progress, "dynamic_analysis.py", 1)
		dynamic_analysis_start_button.SetSensitive(true)
	})

	return vbox
}

func SetupPage() *gtk.VBox {
	vbox := gtk.NewVBox(false, 1)

	static_analysis_frame := gtk.NewFrame("Static Analysis")
	static_analysis_frame.SetBorderWidth(5)
	setup.SetupAndroguard(static_analysis_frame)

	dynamic_analysis_frame := gtk.NewFrame("Dynamic Analysis")
	dynamic_analysis_frame.SetBorderWidth(5)
	//setup.SetupAndroguard(dynamic_analysis_frame)

	vbox.Add(static_analysis_frame)
	vbox.Add(dynamic_analysis_frame)
	return vbox
}

func CreateUIManager() *gtk.UIManager {
	UI_INFO := `
<ui>
  <menubar name='MenuBar'>
    <menu action='FileMenu'>
      <menuitem action='FileQuit' />
    </menu>
    <menu action='HelpMenu'>
      <menuitem action='HelpAbout'/>
    </menu>
  </menubar>
</ui>
`
	ui_manager := gtk.NewUIManager()
	ui_manager.AddUIFromString(UI_INFO)
	return ui_manager
}

func main() {
	gtk.Init(nil)
	window := CreateWindow()
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.Connect("destroy", gtk.MainQuit)
	window.ShowAll()
	gtk.Main()
}
