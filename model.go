package main

// Estructuras para parseo YAML
type YAMLMachine struct {
	QStates          QStates         `yaml:"q_states"`
	Alphabet         []string        `yaml:"alphabet"`
	TapeAlphabet     []string        `yaml:"tape_alphabet"`
	Delta            []YAMLDeltaItem `yaml:"delta"`
	SimulationInputs []string        `yaml:"simulation_strings"`
}

type QStates struct {
	QList   []string `yaml:"q_list"`
	Initial string   `yaml:"initial"`
	Final   string   `yaml:"final"`
}

type YAMLDeltaItem struct {
	Params YAMLDeltaParams `yaml:"params"`
	Output YAMLDeltaOutput `yaml:"output"`
}

type YAMLDeltaParams struct {
	InitialState   string `yaml:"initial_state"`
	MemCacheValue  string `yaml:"mem_cache_value"` // puede ser vacío (blank)
	TapeInput      string `yaml:"tape_input"`
}

type YAMLDeltaOutput struct {
	FinalState      string `yaml:"final_state"`
	MemCacheValue   string `yaml:"mem_cache_value"` // valor a dejar en caché (puede ser vacío)
	TapeOutput      string `yaml:"tape_output"`
	TapeDisplacement string `yaml:"tape_displacement"` // "L", "R", "S"
}
