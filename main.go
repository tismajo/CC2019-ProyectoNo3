package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func main() {
	yamlPath := flag.String("config", "info.yaml", "ruta al archivo YAML con la descripción de la MT")
	flag.Parse()

	data, err := ioutil.ReadFile(*yamlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error leyendo archivo YAML: %v\n", err)
		os.Exit(1)
	}

	var parsed YAMLMachine
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		fmt.Fprintf(os.Stderr, "Error parseando YAML: %v\n", err)
		os.Exit(1)
	}

	tm := NewTuringFromYAML(parsed)

	fmt.Println("== Máquina cargada ==")
	fmt.Printf("Estados: %v\n", tm.States)
	fmt.Printf("Estado inicial: %s\n", tm.InitialState)
	fmt.Printf("Estado final: %s\n", tm.FinalState)
	fmt.Printf("Alfabeto: %v\n", tm.Alphabet)
	fmt.Printf("Alfabeto de cinta: %v\n", tm.TapeAlphabet)
	fmt.Printf("Transiciones: %d\n\n", len(tm.Transitions))

	// Simular cada cadena
	for _, input := range parsed.SimulationInputs {
		fmt.Println("======================================")
		fmt.Printf("Simulando entrada: \"%s\"\n", input)
		ids, accepted, err := tm.Simulate(input)
		if err != nil {
			fmt.Printf("Error durante simulación: %v\n", err)
			continue
		}
		fmt.Println("Descripciones instantáneas (ID) por paso:")
		for i, id := range ids {
			fmt.Printf("%3d: %s\n", i, id)
		}
		if accepted {
			fmt.Printf(">> Resultado: ACEPTADA (llegó al estado final %s)\n", tm.FinalState)
		} else {
			fmt.Printf(">> Resultado: RECHAZADA (no hay transición aplicable hacia estado final)\n")
		}
		fmt.Println("======================================\n")
	}
}
