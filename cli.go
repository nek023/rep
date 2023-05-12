package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type CLI struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
	rep    replacement
}

type replacement struct {
	Old string
	New string
}

func NewCLI(stdout, stderr io.Writer, stdin io.Reader) *CLI {
	return &CLI{
		stdout: os.Stdout,
		stderr: os.Stderr,
		stdin:  os.Stdin,
	}
}

func (cli *CLI) Run(args []string) error {
	var filename string
	if len(args) > 1 {
		filename = args[1]
	}

	text, err := cli.readInput(filename)
	if err != nil {
		return err
	}

	app := tview.NewApplication()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(text)

	inputField := tview.NewInputField().
		SetLabel("%s/").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetChangedFunc(func(str string) {
			r := parseCommand(str)
			cli.rep = r

			go func(r replacement) {
				w := textView.BatchWriter()
				defer w.Close()

				w.Clear()
				err := renderReplacement(w, text, r)
				if err != nil {
					panic(err)
				}
			}(r)
		}).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()

			if key == tcell.KeyEnter {
				r := strings.NewReplacer(cli.rep.Old, cli.rep.New)
				_, err := r.WriteString(cli.stdout, text)
				if err != nil {
					panic(err)
				}
			}
		})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(inputField, 1, 0, true).
		AddItem(textView, 0, 1, true)

	return app.SetRoot(flex, true).Run()
}

func (cli *CLI) readInput(name string) (string, error) {
	var r io.Reader
	if name == "" || name == "-" {
		r = cli.stdin
	} else {
		f, err := os.Open(name)
		if err != nil {
			return "", errors.WithStack(err)
		}
		defer f.Close()
		r = f
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(b), nil
}

func parseCommand(str string) replacement {
	elems := strings.SplitN(str, "/", 2)
	rep := replacement{Old: elems[0]}
	if len(elems) >= 2 {
		rep.New = elems[1]
	}
	return rep
}

func renderReplacement(w io.Writer, str string, rep replacement) error {
	if rep.Old == "" {
		fmt.Fprint(w, str)
		return nil
	}

	r := strings.NewReplacer(rep.Old, fmt.Sprintf("[black:red]%s[black:green]%s[white:black]", rep.Old, rep.New))
	_, err := r.WriteString(w, str)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
