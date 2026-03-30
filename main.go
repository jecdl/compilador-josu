package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Tabla de símbolos para rastrear variables
var symbolTable = make(map[string]bool)

func main() {
	// 1. Configuración de archivos
	inputName := "main.josu"
	outputDir := "build"
	outputName := outputDir + "/main.go"

	// Crear carpeta de salida si no existe
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	inputFile, err := os.Open(inputName)
	if err != nil {
		fmt.Printf("Error: No se encontro el archivo %s\n", inputName)
		return
	}
	defer inputFile.Close()

	outputFile, _ := os.Create(outputName)
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	
	// Cabecera del archivo Go generado
	writer.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n")

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// --- ANALIZADOR SINTÁCTICO ---
		if strings.HasPrefix(line, "print ") {
			// Manejo de print
			content := strings.TrimPrefix(line, "print ")
			writer.WriteString(fmt.Sprintf("\tfmt.Println(%s)\n", content))
		} else if strings.Contains(line, "=") {
			// Manejo de variables
			parts := strings.SplitN(line, "=", 2)
			varName := strings.TrimSpace(parts[0])
			varValue := strings.TrimSpace(parts[1])

			if !symbolTable[varName] {
				// Declaración inicial
				writer.WriteString(fmt.Sprintf("\tvar %s = %s\n", varName, varValue))
				symbolTable[varName] = true
			} else {
				// Re-asignación
				writer.WriteString(fmt.Sprintf("\t%s = %s\n", varName, varValue))
			}
		}
	}

	writer.WriteString("}\n")
	writer.Flush()
	fmt.Printf("Compilacion finalizada. El resultado esta en: %s\n", outputName)
}