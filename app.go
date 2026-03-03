package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

var GameFolderPath string

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OpenGitHub opens the GitHub repo page in the default browser
func (a *App) OpenGitHub() {
	runtime.BrowserOpenURL(a.ctx, "https://github.com/SDxBacon/kcd2-mod-dualdialog-tool")
}

// OpenGitHub opens the GitHub repo page in the default browser
func (a *App) OpenNexusMod() {
	runtime.BrowserOpenURL(a.ctx, "https://www.nexusmods.com/kingdomcomedeliverance2/mods/656")
}

// SelectFolder selects KCM2 folder and returns the path
func (a *App) SelectGameFolder() (string, error) {
	options := runtime.OpenDialogOptions{
		Title:                "Please select the Kingdom Come: Deliverance II folder",
		CanCreateDirectories: false,
	}
	folderPath, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		GameFolderPath = ""
		return "", err
	}

	// check if the folder is a valid KCM2 folder
	err = validateKCM2Folder(folderPath)
	if err != nil {
		GameFolderPath = ""
		return "", err
	}

	// Save the folder path
	GameFolderPath = folderPath
	// Return the folder path
	return folderPath, nil
}

func (a *App) CreateModZip(
	mainLanguage string,
	subLanguage string,
	categories []CategoryConfig,
	mainPatchPaths []string,
	subPatchPaths []string,
) (string, error) {
	// Ask user to select output folder
	options := runtime.OpenDialogOptions{
		Title:                "Please select the output folder",
		CanCreateDirectories: false,
	}
	outputFolder, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil || outputFolder == "" {
		return "", err
	}

	// Process and export the mod zip
	if err := ProcessAndExportModZip(mainLanguage, subLanguage, outputFolder, categories, mainPatchPaths, subPatchPaths); err != nil {
		return "", err
	}

	// return zip file location and nil as error
	return filepath.Join(outputFolder, "Dual Dialog.zip"), nil
}

// GetDefaultCategories returns the built-in category configurations to the frontend.
func (a *App) GetDefaultCategories() []CategoryConfig {
	return DefaultCategories()
}

// SelectPatchFile opens a file-picker dialog filtered to XML files and returns the chosen path.
func (a *App) SelectPatchFile() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Select a community patch XML file",
		Filters: []runtime.FileFilter{
			{DisplayName: "XML Files (*.xml)", Pattern: "*.xml"},
		},
	}
	path, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return "", err
	}
	return path, nil
}

// PreviewPatchFile returns the number of entries found in a patch XML file.
func (a *App) PreviewPatchFile(xmlPath string) (int, error) {
	return CountPatchEntries(xmlPath)
}

// IsValidKCM2Folder checks if the given path is a valid KCM2 folder.
// It only requires a "Localization" subfolder to be present, so users
// can test with just the language pak files without a full game install.
func validateKCM2Folder(path string) error {
	localizationPath := filepath.Join(path, "Localization")
	info, err := os.Stat(localizationPath)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("Localization folder not found in the selected directory")
	}
	return nil
}


