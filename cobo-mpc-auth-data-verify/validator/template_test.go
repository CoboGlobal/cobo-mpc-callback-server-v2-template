package validator

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplatesParse(t *testing.T) {
	templateDir := "template_datas/json_templates"

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Skipf("Template directory %s does not exist, skipping test", templateDir)
		return
	}

	var successCount, failureCount int
	var failures []string

	err := filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".j2") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: failed to read file - %v", path, err))
			failureCount++
			return nil
		}

		templateContent := string(content)

		_, err = getGonjaTemplate(templateContent)
		if err != nil {
			t.Errorf("%s: %v", path, err)
			failures = append(failures, fmt.Sprintf("%s: %v", path, err))
			failureCount++
		} else {
			successCount++
		}

		return nil
	})

	if err != nil {
		t.Errorf("Error walking template directory: %v", err)
		return
	}

	// output statistics
	t.Logf("Template parsing results:")
	t.Logf("  Total templates processed: %d", successCount+failureCount)
	t.Logf("  Successfully parsed: %d", successCount)
	t.Logf("  Failed to parse: %d", failureCount)

	if len(failures) > 0 {
		t.Logf("Failed templates:")
		for _, failure := range failures {
			t.Logf("  - %s", failure)
		}
	}

	if failureCount > 0 {
		t.Errorf("Failed to parse %d templates", failureCount)
	}
}

func TestTemplateBuildMessage(t *testing.T) {
	exampleDataDir := "template_datas/example_datas"
	templateDir := "template_datas/json_templates"

	if _, err := os.Stat(exampleDataDir); os.IsNotExist(err) {
		t.Skipf("Example data directory %s does not exist, skipping test", exampleDataDir)
		return
	}

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Skipf("Template directory %s does not exist, skipping test", templateDir)
		return
	}

	var successCount, failureCount int
	var failures []string

	err := filepath.WalkDir(exampleDataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}

		// parse filename
		fileName := d.Name()
		fileNameWithoutExt := strings.TrimSuffix(fileName, ".json")

		// find version pattern: _X.Y.Z_ or _X.Y.Z
		// use regex or string split to extract
		parts := strings.Split(fileNameWithoutExt, "_")
		if len(parts) < 3 {
			failures = append(failures, fmt.Sprintf("%s: cannot parse filename format", path))
			failureCount++
			return nil
		}

		// find version part (format: X.Y.Z)
		var templateName string
		var version string
		var foundVersion bool

		for i, part := range parts {
			// check if it is version format (X.Y.Z)
			if strings.Contains(part, ".") && len(strings.Split(part, ".")) == 3 {
				// check if all are digits and dots
				isVersion := true
				for _, char := range part {
					if char != '.' && (char < '0' || char > '9') {
						isVersion = false
						break
					}
				}

				if isVersion {
					version = part
					// template name is all parts before version
					templateName = strings.Join(parts[:i], "_")
					foundVersion = true
					break
				}
			}
		}

		if !foundVersion {
			failures = append(failures, fmt.Sprintf("%s: cannot find version number in filename", path))
			failureCount++
			return nil
		}

		// build expected template file name
		expectedTemplateName := fmt.Sprintf("%s_%s.json.j2", templateName, version)
		templatePath := filepath.Join(templateDir, expectedTemplateName)

		// check if template file exists
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			failures = append(failures, fmt.Sprintf("%s: template file %s does not exist", path, expectedTemplateName))
			failureCount++
			return nil
		}

		// read template file content
		templateContent, err := os.ReadFile(templatePath)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: failed to read template file %s - %v", path, expectedTemplateName, err))
			failureCount++
			return nil
		}

		// read example data file content
		exampleDataContent, err := os.ReadFile(path)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: failed to read example data file - %v", path, err))
			failureCount++
			return nil
		}

		// build message from template
		s := NewStatementBuilder(string(templateContent))
		message, err := s.Build(string(exampleDataContent))
		if err != nil {
			t.Errorf("%s: failed to build message from template %s - %v", path, templatePath, err)
			failures = append(failures, fmt.Sprintf("%s: failed to build message from template %s - %v", path, templatePath, err))
			failureCount++
			return nil
		}

		// convert message to interface
		var data interface{}
		err = json.Unmarshal([]byte(message), &data)
		if err != nil {
			t.Errorf("%s: template %s failed to unmarshal message to interface - %v", path, templatePath, err)
			failures = append(failures, fmt.Sprintf("%s: template %s failed to unmarshal message to interface - %v", path, templatePath, err))
			failureCount++
			return nil
		}
		//fmt.Printf("path: %v\n", fileName)
		//fmt.Printf("temp: %v\n", expectedTemplateName)
		successCount++

		return nil
	})

	if err != nil {
		t.Errorf("Error walking example data directory: %v", err)
		return
	}

	// output statistics
	t.Logf("Template mapping and message building results:")
	t.Logf("  Total example files processed: %d", successCount+failureCount)
	t.Logf("  Successfully processed: %d", successCount)
	t.Logf("  Failed to process: %d", failureCount)

	// if there are failures, output detailed information
	if len(failures) > 0 {
		t.Logf("Failed processing:")
		for _, failure := range failures {
			t.Logf("  - %s", failure)
		}
	}

	// if there are failures, let the test fail
	if failureCount > 0 {
		t.Errorf("Failed to process %d example files", failureCount)
	}
}
