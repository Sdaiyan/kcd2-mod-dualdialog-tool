package main

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Languages represents the supported game languages
type Languages int

// Language enum values
const (
	English Languages = iota
	ChineseSimplified
	ChineseTraditional
	Japanese
	Korean
	Spanish
	French
	German
	Italian
	Polish
	Czech
	Portuguese
	Turkish
	Ukrainian
	Russian
)

// String returns the string representation of the language
func (l Languages) String() string {
	return [...]string{
		"English",
		"Chinese - Simplified",
		"Chinese - Traditional",
		"Japanese",
		"Korean",
		"Spanish",
		"French",
		"German",
		"Italian",
		"Polish",
		"Czech",
		"Portuguese",
		"Turkish",
		"Ukrainian",
		"Russian",
	}[l]
}

// GetLanguagePakName returns the PAK file name for the language
func (l Languages) GetLanguagePakName() string {
	switch l {
	case English:
		return "English_xml.pak"
	case ChineseSimplified:
		return "Chineses_xml.pak"
	case ChineseTraditional:
		return "Chineset_xml.pak"
	case Japanese:
		return "Japanese_xml.pak"
	case Korean:
		return "Korean_xml.pak"
	case Spanish:
		return "Spanish_xml.pak"
	case French:
		return "French_xml.pak"
	case German:
		return "German_xml.pak"
	case Italian:
		return "Italian_xml.pak"
	case Polish:
		return "Polish_xml.pak"
	case Czech:
		return "Czech_xml.pak"
	case Portuguese:
		return "Portuguese_xml.pak"
	case Turkish:
		return "Turkish_xml.pak"
	case Ukrainian:
		return "Ukrainian_xml.pak"
	case Russian:
		return "Russian_xml.pak"
	default:
		return "English_xml.pak"
	}
}

// LanguageFromString returns the Languages enum from a string representation
func LanguageFromString(s string) (Languages, bool) {
	for i, name := range [...]string{
		"English",
		"Chinese - Simplified",
		"Chinese - Traditional",
		"Japanese",
		"Korean",
		"Spanish",
		"French",
		"German",
		"Italian",
		"Polish",
		"Czech",
		"Portuguese",
		"Turkish",
		"Ukrainian",
		"Russian",
	} {
		if name == s {
			return Languages(i), true
		}
	}
	return English, false
}

// rowRegex matches <Row> entries including those with multi-line Cell content.
// (?s) enables dotall mode so '.' matches newline characters.
var rowRegex = regexp.MustCompile(`(?s)<Row><Cell>(.*?)</Cell><Cell>(.*?)</Cell><Cell>(.*?)</Cell></Row>`)

// CategoryConfig defines configuration for a single text category to bilingualize.
type CategoryConfig struct {
	ID         string `json:"id"`
	SourceFile string `json:"sourceFile"`
	OutputFile string `json:"outputFile"`
	// Separator is the literal string inserted between main and sub language text.
	// Use `\n` (backslash-n, 2 chars) for in-game line break, or e.g. " / " for inline slash.
	Separator string `json:"separator"`
	Enabled   bool   `json:"enabled"`
}

// DefaultCategories returns the built-in category configurations.
func DefaultCategories() []CategoryConfig {
	return []CategoryConfig{
		{ID: "dialog", SourceFile: "text_ui_dialog.xml", OutputFile: "text_dualdialog_dialog.xml", Separator: `\n`, Enabled: true},
		{ID: "quest", SourceFile: "text_ui_quest.xml", OutputFile: "text_dualdialog_quest.xml", Separator: `\n`, Enabled: true},
		{ID: "items", SourceFile: "text_ui_items.xml", OutputFile: "text_dualdialog_items.xml", Separator: ` / `, Enabled: true},
		{ID: "soul", SourceFile: "text_ui_soul.xml", OutputFile: "text_dualdialog_soul.xml", Separator: ` / `, Enabled: true},
		{ID: "misc", SourceFile: "text_ui_misc.xml", OutputFile: "text_dualdialog_misc.xml", Separator: ` / `, Enabled: true},
		{ID: "minigames", SourceFile: "text_ui_minigames.xml", OutputFile: "text_dualdialog_minigames.xml", Separator: ` / `, Enabled: true},
		{ID: "menus", SourceFile: "text_ui_menus.xml", OutputFile: "text_dualdialog_menus.xml", Separator: `\n`, Enabled: false},
		{ID: "tutorials", SourceFile: "text_ui_tutorials.xml", OutputFile: "text_dualdialog_tutorials.xml", Separator: `\n`, Enabled: false},
	}
}

// readFileFromPakReader extracts a named file's bytes from an already-open zip reader.
func readFileFromPakReader(reader *zip.ReadCloser, fileName string) ([]byte, error) {
	for _, f := range reader.File {
		if f.Name == fileName {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("%s not found in pak", fileName)
}

// isNearlyContained returns true when one string wholly contains the other.
func isNearlyContained(str1, str2 string) bool {
	return strings.Contains(str1, str2) || strings.Contains(str2, str1)
}

// buildLangMap parses XML data and returns a key → col3 map.
// Uses dotall regex so multi-line <Cell> entries are handled correctly.
func buildLangMap(data []byte) map[string]string {
	matches := rowRegex.FindAllSubmatch(data, -1)
	result := make(map[string]string, len(matches))
	for _, m := range matches {
		result[string(m[1])] = string(m[3])
	}
	return result
}

// loadPatchMap reads a local XML patch file and returns a key → col3 map.
func loadPatchMap(xmlPath string) (map[string]string, error) {
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read patch file %s: %w", xmlPath, err)
	}
	return buildLangMap(data), nil
}

// CountPatchEntries returns the number of entries in a patch XML file.
func CountPatchEntries(xmlPath string) (int, error) {
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read patch file %s: %w", xmlPath, err)
	}
	return len(rowRegex.FindAllSubmatch(data, -1)), nil
}

// applyPatches overlays one or more patch maps onto a base map (later patches win).
func applyPatches(base map[string]string, patches ...map[string]string) map[string]string {
	for _, patch := range patches {
		for k, v := range patch {
			base[k] = v
		}
	}
	return base
}

// generateMergedXML produces the merged dual-language XML string for one category.
// mainData is the source XML; mainPatch overrides main-language col3 values;
// subMap provides the secondary-language lookup; separator is inserted between them.
func generateMergedXML(mainData []byte, mainPatch, subMap map[string]string, separator string) string {
	matches := rowRegex.FindAllSubmatch(mainData, -1)

	var sb strings.Builder
	sb.Grow(len(mainData) + len(mainData)/2)
	sb.WriteString("<Table>\n")

	for _, m := range matches {
		key := string(m[1])
		col2 := string(m[2])
		col3 := string(m[3])

		// Apply main-language patch if a correction exists for this key.
		if patched, ok := mainPatch[key]; ok {
			col3 = patched
		}

		textInSub, exists := subMap[key]

		if !exists || isNearlyContained(col3, textInSub) || col3 == textInSub {
			fmt.Fprintf(&sb, "<Row><Cell>%s</Cell><Cell>%s</Cell><Cell>%s</Cell></Row>\n", key, col2, col3)
		} else {
			fmt.Fprintf(&sb, "<Row><Cell>%s</Cell><Cell>%s</Cell><Cell>%s%s%s</Cell></Row>\n", key, col2, col3, separator, textInSub)
		}
	}

	sb.WriteString("</Table>")
	return sb.String()
}

// ProcessAndExportModZip generates "Dual Dialog.zip" in outputFolder.
// categories controls which XML files are processed and with which separator.
// mainPatchPaths / subPatchPaths are optional paths to community-fix XML files
// that are merged over the respective language pak data before bilingualization.
func ProcessAndExportModZip(
	mainLanguage, subLanguage, outputFolder string,
	categories []CategoryConfig,
	mainPatchPaths, subPatchPaths []string,
) error {
	mainLang, _ := LanguageFromString(mainLanguage)
	subLang, _ := LanguageFromString(subLanguage)

	mainPakPath := filepath.Join(GameFolderPath, "Localization", mainLang.GetLanguagePakName())
	subPakPath := filepath.Join(GameFolderPath, "Localization", subLang.GetLanguagePakName())

	fmt.Println("Main pak:", mainPakPath)
	fmt.Println("Sub  pak:", subPakPath)

	// Load community patch maps (main language).
	mainPatchMaps := make([]map[string]string, 0, len(mainPatchPaths))
	for _, p := range mainPatchPaths {
		pm, err := loadPatchMap(p)
		if err != nil {
			return err
		}
		mainPatchMaps = append(mainPatchMaps, pm)
	}

	// Load community patch maps (sub language).
	subPatchMaps := make([]map[string]string, 0, len(subPatchPaths))
	for _, p := range subPatchPaths {
		pm, err := loadPatchMap(p)
		if err != nil {
			return err
		}
		subPatchMaps = append(subPatchMaps, pm)
	}

	// Open PAK files once; we will read multiple XML files from each.
	mainPakReader, err := zip.OpenReader(mainPakPath)
	if err != nil {
		return fmt.Errorf("failed to open main pak: %w", err)
	}
	defer mainPakReader.Close()

	subPakReader, err := zip.OpenReader(subPakPath)
	if err != nil {
		return fmt.Errorf("failed to open sub pak: %w", err)
	}
	defer subPakReader.Close()

	// Build the inner PAK (a ZIP in memory) containing one merged XML per category.
	var innerBuf bytes.Buffer
	innerZip := zip.NewWriter(&innerBuf)
	innerZip.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.DefaultCompression)
	})

	for _, cat := range categories {
		if !cat.Enabled {
			continue
		}

		mainData, err := readFileFromPakReader(mainPakReader, cat.SourceFile)
		if err != nil {
			fmt.Printf("Warning: %s not in main pak, skipping: %v\n", cat.SourceFile, err)
			continue
		}

		subData, err := readFileFromPakReader(subPakReader, cat.SourceFile)
		if err != nil {
			fmt.Printf("Warning: %s not in sub pak, skipping: %v\n", cat.SourceFile, err)
			continue
		}

		// Build the sub-language lookup map and apply sub patches.
		subMap := buildLangMap(subData)
		applyPatches(subMap, subPatchMaps...)

		// Build the main-language patch overlay.
		mainPatch := make(map[string]string)
		applyPatches(mainPatch, mainPatchMaps...)

		mergedXML := generateMergedXML(mainData, mainPatch, subMap, cat.Separator)

		fh := zip.FileHeader{
			Name:     cat.OutputFile,
			Method:   zip.Deflate,
			Modified: time.Now().Truncate(time.Second),
		}
		w, err := innerZip.CreateHeader(&fh)
		if err != nil {
			return fmt.Errorf("failed to create entry %s in inner pak: %w", cat.OutputFile, err)
		}
		if _, err = io.WriteString(w, mergedXML); err != nil {
			return fmt.Errorf("failed to write %s to inner pak: %w", cat.OutputFile, err)
		}
		fmt.Printf("Processed category: %s (%s)\n", cat.ID, cat.OutputFile)
	}

	if err := innerZip.Close(); err != nil {
		return fmt.Errorf("failed to finalise inner pak: %w", err)
	}

	// Write the outer "Dual Dialog.zip".
	outputZipPath := filepath.Join(outputFolder, "Dual Dialog.zip")
	outputZipFile, err := os.Create(outputZipPath)
	if err != nil {
		return fmt.Errorf("failed to create output zip: %w", err)
	}
	defer outputZipFile.Close()

	outputZipWriter := zip.NewWriter(outputZipFile)
	defer outputZipWriter.Close()

	// Inner pak goes in as Localization/{mainLang}.pak
	pakEntry, err := outputZipWriter.Create(fmt.Sprintf("Localization/%s", mainLang.GetLanguagePakName()))
	if err != nil {
		return fmt.Errorf("failed to create pak entry in output zip: %w", err)
	}
	if _, err = pakEntry.Write(innerBuf.Bytes()); err != nil {
		return fmt.Errorf("failed to write pak to output zip: %w", err)
	}

	// mod.manifest
	manifestEntry, err := outputZipWriter.Create("mod.manifest")
	if err != nil {
		return fmt.Errorf("failed to create manifest entry in zip: %w", err)
	}
	if _, err = io.WriteString(manifestEntry, manifest); err != nil {
		return fmt.Errorf("failed to write manifest to zip: %w", err)
	}

	return nil
}
