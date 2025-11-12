package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func updateAPIRoutesFile(modelName string) error {
	routesFile := "../api/routes.go"

	content, err := os.ReadFile(routesFile)
	if err != nil {
		return fmt.Errorf("error reading routes file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	insertIndex := -1
	for i, line := range lines {
		if strings.Contains(line, "app.Controllers.") && i > 0 {
			insertIndex = i + 1
		}
	}

	if insertIndex == -1 {
		return fmt.Errorf("could not find insertion point in routes file")
	}

	newLine := fmt.Sprintf("\tapp.Controllers.%s.SetBasicRoutes(deps)", modelName)

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:insertIndex]...)
	newLines = append(newLines, newLine)
	newLines = append(newLines, lines[insertIndex:]...)

	err = os.WriteFile(routesFile, []byte(strings.Join(newLines, "\n")), 0o644)
	if err != nil {
		return fmt.Errorf("error writing routes file: %w", err)
	}

	log.Printf("updated: %s with %s controller", routesFile, modelName)
	return nil
}

func updateControllersFile(modelName string) error {
	controllersFile := "../controllers/controllers.go"

	content, err := os.ReadFile(controllersFile)
	if err != nil {
		return fmt.Errorf("error reading controllers file: %w", err)
	}

	importLine := fmt.Sprintf(`	"app/controllers/%s_controller"`, strings.ToLower(modelName))
	structLine := fmt.Sprintf("\t%s *%s_controller.Controllers", modelName, strings.ToLower(modelName))

	setupLine := fmt.Sprintf("\t\t%s: %s_controller.Get(d),", modelName, strings.ToLower(modelName))

	lines := strings.Split(string(content), "\n")

	newLines := make([]string, 0, len(lines)+3)
	inImports := false
	importsAdded := false
	structAdded := false
	setupAdded := false

	for _, line := range lines {
		if strings.Contains(line, "import (") {
			inImports = true
			newLines = append(newLines, line)
			continue
		}

		if inImports && !importsAdded {
			if strings.Contains(line, ")") {
				newLines = append(newLines, importLine)
				importsAdded = true
				inImports = false
			}
		}

		if strings.Contains(line, "type Controllers struct {") && !structAdded {
			newLines = append(newLines, line)
			inserted := false
			for i := len(newLines) - 1; i >= 0; i-- {
				if strings.Contains(newLines[i], "type Controllers struct {") {
					// Insert after the opening brace
					newLines = append(newLines[:i+1], append([]string{structLine}, newLines[i+1:]...)...)
					structAdded = true
					inserted = true
					break
				}
			}
			if !inserted {
				newLines = append(newLines, structLine)
				structAdded = true
			}
			continue
		}

		if strings.Contains(line, "return &Controllers{") && !setupAdded {
			newLines = append(newLines, line)
			inserted := false
			for i := len(newLines) - 1; i >= 0; i-- {
				if strings.Contains(newLines[i], "return &Controllers{") {
					// Insert after the opening brace
					newLines = append(newLines[:i+1], append([]string{setupLine}, newLines[i+1:]...)...)
					setupAdded = true
					inserted = true
					break
				}
			}
			if !inserted {
				newLines = append(newLines, setupLine)
				setupAdded = true
			}
			continue
		}

		newLines = append(newLines, line)
	}

	err = os.WriteFile(controllersFile, []byte(strings.Join(newLines, "\n")), 0o644)
	if err != nil {
		return fmt.Errorf("error writing controllers file: %w", err)
	}

	log.Printf("updated: %s with %s controller", controllersFile, modelName)
	return nil
}

func updateModelsFile(modelName string) error {
	modelsFile := "../models/models.go"

	content, err := os.ReadFile(modelsFile)
	if err != nil {
		return fmt.Errorf("error reading models file: %w", err)
	}
	importLine := fmt.Sprintf(`	"app/models/%s"`, strings.ToLower(modelName))

	structLine := fmt.Sprintf("\t%s *%s.Queries", modelName, strings.ToLower(modelName))
	setupLine := fmt.Sprintf("\t\t%s: %s.New(d),", modelName, strings.ToLower(modelName))

	lines := strings.Split(string(content), "\n")

	newLines := make([]string, 0, len(lines)+3)
	inImports := false
	importsAdded := false
	structAdded := false
	setupAdded := false

	for _, line := range lines {
		if strings.Contains(line, "import (") {
			inImports = true
			newLines = append(newLines, line)
			continue
		}

		if inImports && !importsAdded {
			if strings.Contains(line, ")") {
				newLines = append(newLines, importLine)
				importsAdded = true
				inImports = false
			}
		}

		if strings.Contains(line, "type Models struct {") && !structAdded {
			newLines = append(newLines, line)
			inserted := false
			for i := len(newLines) - 1; i >= 0; i-- {
				if strings.Contains(newLines[i], "type Models struct {") {
					// Insert after the opening brace
					newLines = append(newLines[:i+1], append([]string{structLine}, newLines[i+1:]...)...)
					structAdded = true
					inserted = true
					break
				}
			}
			if !inserted {
				newLines = append(newLines, structLine)
				structAdded = true
			}
			continue
		}

		if strings.Contains(line, "return &Models{") && !setupAdded {
			newLines = append(newLines, line)
			inserted := false
			for i := len(newLines) - 1; i >= 0; i-- {
				if strings.Contains(newLines[i], "return &Models{") {
					newLines = append(newLines[:i+1], append([]string{setupLine}, newLines[i+1:]...)...)
					setupAdded = true
					inserted = true
					break
				}
			}
			if !inserted {
				newLines = append(newLines, setupLine)
				setupAdded = true
			}
			continue
		}

		newLines = append(newLines, line)
	}

	err = os.WriteFile(modelsFile, []byte(strings.Join(newLines, "\n")), 0o644)
	if err != nil {
		return fmt.Errorf("error writing models file: %w", err)
	}

	log.Printf("updated: %s with %s model", modelsFile, modelName)
	return nil
}
