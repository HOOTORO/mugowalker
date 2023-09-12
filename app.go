package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type OptionsType struct {
	CaseInsensitive bool `json:"caseInsensitive,omitempty"`
	WholeWord       bool `json:"wholeWord,omitempty"`
	WholeLine       bool `json:"wholeLine,omitempty"`
	FilenameOnly    bool `json:"filenameOnly,omitempty"`
	FilesWoMatches  bool `json:"filesWoMatches,omitempty"`
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Connected to  %s!", name)
}

func (a *App) SelectFolder() string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Folder",
	})
	if err != nil {
		log.Println("Error selecting a folder")
	}
	return selection
}

func (a *App) SelectFile() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
	})
	if err != nil {
		log.Println("Error selecting a file")
	}
	return selection
}

func (a *App) Search(path string, pattern string, options OptionsType) string {
	if path == "" {
		_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "ERROR",
			Message:       "No path was entered",
			Buttons:       []string{"OK"},
			DefaultButton: "OK",
		})
		if err != nil {
			log.Println(err)
			return ""
		}
		return ""
	}

	if pattern == "" {
		_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "ERROR",
			Message:       "No pattern was entered",
			Buttons:       []string{"OK"},
			DefaultButton: "OK",
		})
		if err != nil {
			log.Println(err)
			return ""
		}
		return ""
	}

	if options.CaseInsensitive {
		pattern = "(?i)" + pattern
	}
	if options.WholeWord {
		pattern = `\b` + pattern + `\b`
	}
	if options.WholeLine {
		pattern = "^" + pattern + "$"
	}

	results, err := walkDir(path, pattern, options)
	if err != nil {
		log.Fatalln("ERROR walking the directory!")
	}

	if len(results) > 0 {
		return results
	} else {
		return "Not results were found"
	}
}

// func (a *App) ParseImage(path string) string {
// 	ip := ocr.ExtractText(path)
// }

func walkDir(dirToWalk string, pattern string, options OptionsType) (string, error) {
	var matches []string
	err := filepath.Walk(dirToWalk, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			r, err := regexp.Compile(pattern)
			if err != nil {
				return err
			}
			if r.MatchString(info.Name()) {
				matches = append(matches, path)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	results := strings.Join(matches, "\n")
	return strings.ReplaceAll(results, dirToWalk, ""), nil
}
