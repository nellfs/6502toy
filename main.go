package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/nellfs/6502toy/internal/cpu"
	"github.com/nellfs/6502toy/internal/terminal"
)

const inputCodeFile = "input/code.asm"
const binaryOutputFile = "build/program.bin"

func writeCode(buff []byte) {
	err := os.WriteFile(inputCodeFile, buff, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	// window setup
	app := adw.NewApplication("com.nellfs.6502toy", gio.ApplicationFlagsNone)

	app.ConnectActivate(func() { activate(app) })
	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}

}

func activate(app *adw.Application) {
	win := adw.NewApplicationWindow(&app.Application)
	win.SetTitle("6502 Toy")
	win.SetSizeRequest(320, 320)
	win.SetDefaultSize(800, 600)

	// sidebar
	sidebarHeaderBar := adw.NewHeaderBar()
	sidebarView := adw.NewToolbarView()
	sidebarView.AddTopBar(sidebarHeaderBar)

	// CPU model
	emu := &cpu.CPU{}
	emu.PC = 0x8000

	// sidebar content: CPU registers
	regsPage := adw.NewPreferencesPage()
	regsPage.SetTitle("CPU")

	regsGroup := adw.NewPreferencesGroup()
	regsGroup.SetTitle("Registers")

	regARow := adw.NewActionRow()
	regARow.SetTitle("A")

	regPCRow := adw.NewActionRow()
	regPCRow.SetTitle("PC")

	regsGroup.Add(regARow)
	regsGroup.Add(regPCRow)
	regsPage.Add(regsGroup)
	sidebarView.SetContent(regsPage)

	updateRegisters := func() {
		regARow.SetSubtitle(fmt.Sprintf("0x%02X", emu.A))
		regPCRow.SetSubtitle(fmt.Sprintf("0x%04X", emu.PC))
	}

	updateRegisters()

	// main headerbar
	mainHeaderBar := adw.NewHeaderBar()
	mainHeaderBar.SetShowTitle(false)

	// main content
	textView := gtk.NewTextView()
	textView.SetMarginStart(8)
	textView.SetMarginBottom(16)
	textView.SetMarginEnd(8)
	textView.AddCSSClass("code-view")

	editorScrolled := gtk.NewScrolledWindow()
	editorScrolled.SetChild(textView)
	editorScrolled.SetVExpand(true)
	editorScrolled.SetHExpand(true)

	// bottom "terminal" for logs/errors
	logView := gtk.NewTextView()

	logWriter := slog.New(terminal.SimpleHandler(terminal.NewTextViewWriter(logView), slog.LevelDebug))
	slog.SetDefault(logWriter)

	logScrolled := gtk.NewScrolledWindow()
	logScrolled.SetChild(logView)
	logScrolled.SetSizeRequest(-1, 200) // minimum height for log/terminal area

	// load existing code, if any
	if data, err := os.ReadFile(inputCodeFile); err == nil {
		textView.Buffer().SetText(string(data))
	}

	// main headerbar button
	buildBtn := gtk.NewButtonFromIconName("weather-tornado-symbolic")
	buildBtn.SetTooltipText("Build")

	buildBtn.ConnectClicked(func() {
		// compile and save code snippet
		if err := compile(textView.Buffer()); err != nil {
			slog.Error("build failed", "error", err)
		}
	})

	runBtn := gtk.NewButtonFromIconName("media-playback-start-symbolic")
	runBtn.SetTooltipText("Run")

	runBtn.ConnectClicked(func() {
		if err := compile(textView.Buffer()); err != nil {
			slog.Error("build before run failed", "error", err)
			return
		}

		// cpu runtime
		data, err := os.ReadFile(binaryOutputFile)
		if err != nil {
			slog.Error("failed to read program binary", "error", err, "path", binaryOutputFile)
			return
		}
		copy(emu.Mem[0x8000:], data)

		// reset CPU state and run until BRK or a safety limit
		emu.PC = 0x8000
		emu.A = 0
		emu.Run(100000)

		updateRegisters()
	})

	mainHeaderBar.PackStart(buildBtn)
	mainHeaderBar.PackStart(runBtn)

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromString(`
textview, textview text {
	background: transparent;
}

.code-view, .code-view text {
	font-size: 16px;
	font-family: monospace;
}

.log-view, .log-view text {
	font-size: 14px;
	font-family: monospace;
}
`)
	gtk.StyleContextAddProviderForDisplay(
		textView.Display(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	mainView := adw.NewToolbarView()
	mainView.AddTopBar(mainHeaderBar)

	// paned main area: editor on top, log "terminal" at the bottom
	mainPaned := gtk.NewPaned(gtk.OrientationVertical)
	mainPaned.SetStartChild(editorScrolled)
	mainPaned.SetEndChild(logScrolled)
	mainPaned.SetResizeStartChild(true)
	mainPaned.SetResizeEndChild(false)
	mainPaned.SetShrinkStartChild(false)
	mainPaned.SetShrinkEndChild(false)
	mainPaned.SetPosition(350)

	mainView.SetContent(mainPaned)

	// split view
	split := adw.NewOverlaySplitView()
	split.SetSidebar(sidebarView)
	split.SetContent(mainView)
	split.SetSidebarPosition(gtk.PackStart)

	win.SetContent(split)
	win.Present()
}

func compile(textViewBuffer *gtk.TextBuffer) error {
	buff := textViewBuffer
	start := buff.StartIter()
	end := buff.EndIter()

	text := buff.Text(start, end, true)
	writeCode([]byte(text))

	err := run("ca65", inputCodeFile, "-o", "build/object.o")
	if err != nil {
		slog.Error("could not assemble code", "error", err)
		return err
	}

	err = run("ld65", "build/object.o", "-t", "none", "-o", binaryOutputFile)
	if err != nil {
		slog.Error("could not link object", "error", err)
		return err
	}
	return nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
