package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

type OptionsType struct {
	CaseInsensitive bool `json:"caseInsensitive,omitempty"`
	WholeWord       bool `json:"wholeWord,omitempty"`
	WholeLine       bool `json:"wholeLine,omitempty"`
	FilenameOnly    bool `json:"filenameOnly,omitempty"`
	FilesWoMatches  bool `json:"filesWoMatches,omitempty"`
}

func (a *App) SelectFolder() string {
	selection, err := rt.OpenDirectoryDialog(a.ctx, rt.OpenDialogOptions{
		Title: "Select Folder",
	})
	if err != nil {
		log.Println("Error selecting a folder")
	}
	return selection
}

func (a *App) SelectFile() string {
	selection, err := rt.OpenFileDialog(a.ctx, rt.OpenDialogOptions{
		Title: "Select File",
	})
	if err != nil {
		log.Println("Error selecting a file")
	}
	return selection
}

func (a *App) Search(path string, pattern string, options OptionsType) string {
	if path == "" {
		_, err := rt.MessageDialog(a.ctx, rt.MessageDialogOptions{
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
		_, err := rt.MessageDialog(a.ctx, rt.MessageDialogOptions{
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
