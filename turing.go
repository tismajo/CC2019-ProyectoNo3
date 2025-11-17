package main

import (
	"errors"
	"fmt"
	"strings"
)

// Transition interna para búsquedas rápidas
type Transition struct {
	FromState      string
	FromCache      string
	ReadSymbol     string
	ToState        string
	ToCache        string
	WriteSymbol    string
	Displacement   string // "L", "R", "S"
}

// Máquina de Turing
type TuringMachine struct {
	States         []string
	InitialState   string
	FinalState     string
	Alphabet       []string
	TapeAlphabet   []string
	Transitions    []Transition
	blankSymbol    string // símbolo lógico de blank internamente (puede ser "")
}

// Creates a new TM from parsed YAMLMachine
func NewTuringFromYAML(y YAMLMachine) *TuringMachine {
	tm := &TuringMachine{
		States:       y.QStates.QList,
		InitialState: y.QStates.Initial,
		FinalState:   y.QStates.Final,
		Alphabet:     y.Alphabet,
		TapeAlphabet: y.TapeAlphabet,
		blankSymbol:  "", // representamos blank internamente como la cadena vacía
	}
	for _, d := range y.Delta {
		tr := Transition{
			FromState:    d.Params.InitialState,
			FromCache:    d.Params.MemCacheValue, // puede ser ""
			ReadSymbol:   d.Params.TapeInput,
			ToState:      d.Output.FinalState,
			ToCache:      d.Output.MemCacheValue,
			WriteSymbol:  d.Output.TapeOutput,
			Displacement: strings.ToUpper(d.Output.TapeDisplacement),
		}
		tm.Transitions = append(tm.Transitions, tr)
	}
	return tm
}

// Tape representada como mapa de posiciones a símbolos (string).
// Posiciones no presentes se interpretan como blank (cadena vacía).
type tape struct {
	cells map[int]string
}

// createTapeFromString carga la cadena de entrada (cada símbolo es 1 char/1 token separado).
// Asume que la entrada ya es una secuencia de símbolos que queremos mapear uno a uno.
func createTapeFromString(s string) *tape {
	t := &tape{cells: make(map[int]string)}
	for i, ch := range s {
		// cada carácter lo guardamos como string con tamaño 1
		t.cells[i] = string(ch)
	}
	return t
}

func (t *tape) read(pos int) string {
	if v, ok := t.cells[pos]; ok {
		return v
	}
	return ""
}

func (t *tape) write(pos int, sym string) {
	if sym == "" {
		// blank: eliminamos la celda para mantener la representación esparcida
		delete(t.cells, pos)
	} else {
		t.cells[pos] = sym
	}
}

// findTransition busca una transición que coincida exactamente con (state, memCache, tapeSymbol).
// Retorna nil si no hay ninguna.
func (tm *TuringMachine) findTransition(state, memCache, tapeSymbol string) *Transition {
	for _, tr := range tm.Transitions {
		// comparaciones exactas; memCache y tapeSymbol pueden ser "" (blank)
		if tr.FromState == state && tr.FromCache == memCache && tr.ReadSymbol == tapeSymbol {
			// copia para devolución
			tmp := tr
			return &tmp
		}
	}
	return nil
}

// displaySymbol convierte símbolo interno "" (blank) a visible "_" para impresión
func (tm *TuringMachine) displaySymbol(sym string) string {
	if sym == "" {
		return "_"
	}
	return sym
}

// buildInstantaneousDescription genera la ID en formato α q β  (y muestra mem_cache)
func (tm *TuringMachine) buildInstantaneousDescription(t *tape, head int, state string, memCache string) string {
	// buscamos rango relevante: desde el mínimo índice con algo hasta máximo.
	min := head
	max := head
	for pos := range t.cells {
		if pos < min {
			min = pos
		}
		if pos > max {
			max = pos
		}
	}
	// si t.cells está vacío, mostramos al menos la posición 0..0
	if len(t.cells) == 0 {
		min = 0
		max = 0
	}

	var leftBuilder strings.Builder
	for i := min; i < head; i++ {
		leftBuilder.WriteString(tm.displaySymbol(t.read(i)))
	}
	var rightBuilder strings.Builder
	// head position symbol y el resto
	currentSymbol := tm.displaySymbol(t.read(head))
	rightBuilder.WriteString(currentSymbol)
	for i := head + 1; i <= max; i++ {
		rightBuilder.WriteString(tm.displaySymbol(t.read(i)))
	}

	// α q β
	id := fmt.Sprintf("%s (%s) %s    [mem=%s]  head=%d", leftBuilder.String(), state, rightBuilder.String(), tm.displaySymbol(memCache), head)
	return id
}

// Simulate ejecuta la TM sobre una cadena de entrada y retorna la lista de IDs y si fue aceptada.
func (tm *TuringMachine) Simulate(input string) ([]string, bool, error) {
	// inicializaciones
	t := createTapeFromString(input)
	head := 0
	state := tm.InitialState
	memCache := "" // inicialmente vacío (blank)
	ids := []string{}

	// guardamos la primera ID (estado inicial antes de aplicar transiciones)
	ids = append(ids, tm.buildInstantaneousDescription(t, head, state, memCache))

	// limitador de pasos para evitar loops infinitos (por seguridad), por defecto alto.
	const maxSteps = 100000
	step := 0

	for {
		if step > maxSteps {
			return ids, false, errors.New("límite de pasos excedido (posible loop infinito)")
		}

		// Si ya estamos en estado final, aceptamos y terminamos (la ID final ya fue agregada)
		if state == tm.FinalState {
			return ids, true, nil
		}

		// leemos símbolo actual en la cinta (puede ser blank -> "")
		currSym := t.read(head)

		tr := tm.findTransition(state, memCache, currSym)
		if tr == nil {
			// no hay transición aplicable -> rechazado
			return ids, false, nil
		}

		// aplicamos la transición
		// escribir símbolo
		t.write(head, tr.WriteSymbol)
		// actualizar memoria caché
		memCache = tr.ToCache
		// mover cabeza
		switch tr.Displacement {
		case "L":
			head--
		case "R":
			head++
		case "S", "": // stay
			// no hacer nada
		default:
			// si viene algo distinto, lo tratamos como error
			return ids, false, fmt.Errorf("desplazamiento inválido en transición: %s", tr.Displacement)
		}
		// cambiar estado
		state = tr.ToState

		// agregamos la ID posterior a la transición
		ids = append(ids, tm.buildInstantaneousDescription(t, head, state, memCache))

		step++
	}
}
