Video: https://www.canva.com/design/DAG45D5Mmu8/rlxdp_hbaS8qQjRLgL68lQ/watch?utm_content=DAG45D5Mmu8&utm_campaign=designshare&utm_medium=link2&utm_source=uniquelinks&utlId=h0356c0d65d 


# Máquina de Turing – Proyecto CC2019

Este proyecto implementa un **simulador de Máquina de Turing** en Go.  
La máquina se define mediante un archivo **YAML**, compatible con estados, memoria, alfabeto, cinta y transiciones formales.  
El simulador ejecuta las descripciones instantáneas paso a paso y muestra si la cadena es **aceptada** o **rechazada**.

---

## Estructura del Proyecto

.
├── main.go
├── turing.go          # Lógica de la máquina
├── info.yaml       # Definición formal de la MT
└── README.txt

---

## Formato del YAML

El YAML define:

- **Estados (`q_states`)**  
- **Alfabeto de entrada (`alphabet`)**  
- **Alfabeto de cinta (`tape_alphabet`)**  
- **Transiciones (`delta`)**  
- **Cadenas a simular (`simulation_strings`)**

Ejemplo mínimo:

q_states:
  q_list: ["q0", "q1", "qf"]
  initial: "q0"
  final: "qf"

alphabet: ["a", "b"]
tape_alphabet: ["a", "b", "_"]

delta:
  - params:
      initial_state: "q0"
      mem_cache_value: ""
      tape_input: "a"
    output:
      final_state: "q1"
      mem_cache_value: ""
      tape_output: "a"
      tape_displacement: "R"

---

## Cómo ejecutar

1. Tener Go instalado (1.20+).
2. Ejecutar:

go run .

El programa cargará automáticamente el archivo `info.yaml` y simulará todas las cadenas listadas.

---

## 🧪 Ejemplo incluido

La máquina de ejemplo **acepta todas las cadenas que terminan en `a`** y rechaza las que terminan en `b`.

Ejemplos aceptados:
- `a`
- `ba`
- `abba`

Ejemplos rechazados:
- `b`
- `bb`
- `bab`
