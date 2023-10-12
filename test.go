package main

import (
	"fmt"
)

// dummy type to make constants more understandable
type TGate int

// enumeration of constants to assign types to gates
const (
	COMP TGate = iota
	IN
	OUT
	NOT
	AND
	OR
)

var component_counter int = 0
var gate_counter      int = 0
var pin_counter       int = 0
var wire_counter      int = 0

// has state and pointer to gate to which it is connected
type Pin struct {
	number	int
	p_state bool
	w_connected bool
	gate 	*Gate
}

// has type, state, array of pointers to pins that are connected
// to it, a pointer to the wire to whuich it is connected
type Gate struct {
	name	string
	number 	int
	g_type	TGate
	g_state bool
	pins 	[]*Pin
	out_pin *Wire
}

// has state and array of pointers to pins to which it is connected
type Wire struct {
	number	 int
	w_state  bool
	transmit []*Pin
}

type Component struct {
	name		string
	number      int
	c_type		TGate
	gates 		[]*Gate
	pins  		[]*Pin
	wires 		[]*Wire
	exposed_in 	[]*Pin
	exposed_out *Wire
}

func NewComponent(name string, number_of_inputs int) *Component {
	component_counter++
	return &Component{name: name, number: component_counter, c_type: COMP, exposed_in: make([]*Pin, number_of_inputs)}
}

func AddC2C(where_to_add *Component, what_to_add *Component) {
	// implement
}

func NewPin(gate *Gate) *Pin {
	pin_counter++
	return &Pin{number: pin_counter, p_state: false, w_connected: false, gate: gate}
}

func NewGate(g_type TGate) *Gate {
	gate_counter++
	gate := Gate{g_type: g_type, number: gate_counter, g_state: false}
	switch g_type {
	case IN:
		gate.name = "IN"
		return &gate
	case OUT:
		gate.name = "OUT"
		pin_in := NewPin(&gate)
		gate.pins = append(gate.pins, pin_in)
		return &gate
	case NOT:
		gate.name = "NOT"
		pin_in := NewPin(&gate)
		gate.pins = append(gate.pins, pin_in)
		return &gate
	case AND:
		gate.name = "AND"
		pin_in_1 := NewPin(&gate)
		pin_in_2 := NewPin(&gate)
		gate.pins = append(gate.pins, pin_in_1, pin_in_2)
		return &gate
	case OR:
		gate.name = "OR"
		pin_in_1 := NewPin(&gate)
		pin_in_2 := NewPin(&gate)
		gate.pins = append(gate.pins, pin_in_1, pin_in_2)
		return &gate
	default:
		fmt.Println("Given type is unknown")
		return nil
	}
}

func NewWire(from_gate *Gate, to_pin *Pin) {
	wire_counter++
	wire := Wire{number: wire_counter, w_state: false, transmit: []*Pin{to_pin}}
	from_gate.out_pin = &wire
}

func UpdateWire(wire *Wire, to_pin *Pin) {
	wire.transmit = append(wire.transmit, to_pin)
}

func G2G(gate_from *Gate, gate_to *Gate, one_out_to_n_inputs int) {
	
	if (gate_from != nil) && (gate_to != nil) && (0 < one_out_to_n_inputs) && (one_out_to_n_inputs < 3){
		
		if gate_from.out_pin != nil {
			
			tmp_size := len(gate_from.out_pin.transmit)
			for _, pin := range gate_to.pins {
				if !pin.w_connected {
					UpdateWire(gate_from.out_pin, pin)
					pin.w_connected = true
					break
				}
			}
			if len(gate_from.out_pin.transmit) == tmp_size {
				fmt.Println("Connection unsuccessful because there are no avialable pins on the Gate 2")
			}

		} else {
			for _, pin := range gate_to.pins {
				if !pin.w_connected {
					NewWire(gate_from, pin)
					pin.w_connected = true
					break
				}
			}
			if gate_from.out_pin == nil {
				fmt.Println("Connection unsuccessful because there are no avialable pins on the Gate 2")
			}
		}

	} else {
		fmt.Println("Connection unsuccessful! RULES: Gate 1 cannot be nil | Gate 2 cannot be nil | Number of inputs should be either 1 or 2")
	}
}

// called when wire should chage its state, the caller is gate
func (w *Wire) wire_report_state(rec_state bool) {
	w.w_state = rec_state
	for _, pin := range w.transmit {
		pin.pin_report_state(w.w_state)
	}
}

// called when pin should chage its state, the caller is wire
func (p *Pin) pin_report_state(rec_state bool) {
	p.p_state = rec_state
	switch p.gate.g_type {
	case IN:
		p.gate.update_state(p.p_state)
	case OUT:
		p.gate.update_state(p.p_state)
	case NOT:
		p.gate.report_state_NOT()
	case AND:
		p.gate.report_state_AND()
	case OR:
		p.gate.report_state_OR()
	default:
		fmt.Println("The gate is of unknown type")
	}
}

// returns gate's state
func (g *Gate) get_state() bool {
	return g.g_state
}

func (g *Gate) get_detailed_state() {
	fmt.Print("\n\n--- Selected gate detailed state ---\n\n")
	fmt.Println("Number of gate element: ", g.number)
	fmt.Println("Type name: ", g.name)
	fmt.Println("State: ", g.g_state)
	fmt.Println("Number of pins: ", len(g.pins))
	for i, pin := range g.pins {
		if pin == nil {
			fmt.Println("_g_ The pin ", i, " is nil!")
		} else {
			fmt.Println("_g_ Number of pin element: ", pin.number)
			fmt.Println("_g_ State: ", pin.p_state)
			fmt.Println("_g_ Is connected: ", pin.w_connected)
		}
	}
	if g.out_pin == nil {
		fmt.Println("There is no output wire!")
	} else {
		fmt.Println("Number of wire element: ", g.out_pin.number)
		fmt.Println("State: ", g.out_pin.w_state)
		fmt.Println("Number of connections", len(g.out_pin.transmit))
		for i, pin := range g.out_pin.transmit {
			if pin == nil {
				fmt.Println("_w_ The pin ", i, " is nil!")
			} else {
				fmt.Println("_w_ Number of pin element: ", pin.number)
				fmt.Println("_w_ State: ", pin.p_state)
				fmt.Println("_w_ Is connected: ", pin.w_connected)
			}
		}
	}
	fmt.Print("\n--- END ---\n\n")
}

// updates gate's state, calles update on wire if connected
func (g *Gate) update_state(rec_state bool) {
	g.g_state = rec_state
	if g.out_pin != nil {
		g.out_pin.wire_report_state(g.g_state)
	}
}

// updates gate's state by AND logic, uses 2 pin connection, calls update on wire if connected
func (g *Gate) report_state_AND() {
	if g.out_pin != nil {
		g.g_state = g.pins[0].p_state && g.pins[1].p_state
		g.out_pin.wire_report_state(g.g_state)
	} else {
		g.g_state = g.pins[0].p_state
	}
}

// updates gate's state by OR logic, uses 2 pin connection, calls update on wire if connected
func (g *Gate) report_state_OR() {
	if g.out_pin != nil {
		g.g_state = g.pins[0].p_state || g.pins[1].p_state
		g.out_pin.wire_report_state(g.g_state)
	} else {
		g.g_state = g.pins[0].p_state
	}
}

// updates gate's state by NOT logic, uses 1 pin connection, calls update on wire if connected
func (g *Gate) report_state_NOT() {
	g.g_state = !g.pins[0].p_state
	if g.out_pin != nil {
		g.out_pin.wire_report_state(g.g_state)
	}
}

func main() {
	fmt.Print("\nSimple Truth Table Solver\n\n")

	//var component = NewComponent("NAND Gate", 2)
	var gate_and = NewGate(AND)
	var gate_not = NewGate(NOT)
	var gate_inbuff_1  = NewGate(IN)
	var gate_inbuff_2  = NewGate(IN)
	var gate_outbuff   = NewGate(OUT)

	G2G(gate_inbuff_1, gate_and, 1)
	G2G(gate_inbuff_2, gate_and, 1)
	G2G(gate_and, gate_not, 1)
	G2G(gate_not, gate_outbuff, 1)

	gate_inbuff_1.update_state(true)
	gate_inbuff_2.update_state(false)

	fmt.Print("\nInputs are: ", gate_inbuff_1.g_state, " and ", gate_inbuff_2.g_state, "\n")
	//gate_and.get_detailed_state()
	fmt.Println("The output is: ", gate_outbuff.get_state())
}