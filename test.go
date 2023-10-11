package main

import "fmt"

// dummy type to make constants more understandable
type TGate int

// enumeration of constants to assign types to gates
const (
	IN TGate = iota
	OUT
	NOT
	AND
)

// has state and pointer to gate to which it is connected
type Pin struct {
	p_state bool
	gate 	*Gate
}

// has type, state, array of pointers to pins that are connected
// to it, a pointer to the wire to whuich it is connected
type Gate struct {
	g_type	TGate
	g_state bool
	pins 	[]*Pin
	out_pin *Wire
}

// has state and array of pointers to pins to which it is connected
type Wire struct {
	w_state  bool
	transmit []*Pin
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
	default:
		fmt.Println("The gate is of unknown type")
	}
}

// returns gate's state
func (g *Gate) get_state() bool {
	return g.g_state
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

// updates gate's state by NOT logic, uses 1 pin connection, calls update on wire if connected
func (g *Gate) report_state_NOT() {
	g.g_state = !g.pins[0].p_state
	if g.out_pin != nil {
		g.out_pin.wire_report_state(g.g_state)
	}
}

func main() {
	fmt.Print("\nSimple Truth Table Solver\n\n")

	// declare all gates, i/o buffers are also gates
	var gate_and Gate
	var gate_not Gate
	var in_buff1 Gate
	var in_buff2 Gate
	var out_buff Gate

	// declare all pins
	var pin_1	Pin
	var pin_2	Pin
	var pin_n	Pin
	var pin_o	Pin

	// declare all wires
	var wire_1  Wire
	var wire_2  Wire
	var wire_n	Wire
	var wire_o  Wire

	// to use easily
	var state_0 = false
	var state_1 = true

	// set states of pins, connect to dedicated gates
	pin_1 = Pin{p_state: state_0, gate: &gate_and}
	pin_2 = Pin{p_state: state_0, gate: &gate_and}
	pin_n = Pin{p_state: state_0, gate: &gate_not}
	pin_o = Pin{p_state: state_0, gate: &out_buff}

	// set types and states of gates, in gate has only output, out gate has only 1 input
	in_buff1 = Gate{g_type: IN, g_state: state_0, out_pin: &wire_1}
	in_buff2 = Gate{g_type: IN, g_state: state_0, out_pin: &wire_2}
	out_buff = Gate{g_type: OUT, g_state: state_0, pins: []*Pin{&pin_o}}
	gate_and = Gate{g_type: AND, g_state: state_0, pins: []*Pin{&pin_1, &pin_2}, out_pin: &wire_n}
	gate_not = Gate{g_type: NOT, g_state: state_0, pins: []*Pin{&pin_n}, out_pin: &wire_o}

	// set states of wires, connect to dedicated pins
	wire_1 = Wire{w_state: state_0, transmit: []*Pin{&pin_1}}
	wire_2 = Wire{w_state: state_0, transmit: []*Pin{&pin_2}}
	wire_n = Wire{w_state: state_0, transmit: []*Pin{&pin_n}}
	wire_o = Wire{w_state: state_0, transmit: []*Pin{&pin_o}}

	// can be used multiple times
	in_buff1.update_state(state_1)
	in_buff2.update_state(state_1)

	// PRINTING
	fmt.Print("Inputs are: ", in_buff1.g_state, " and ", in_buff2.g_state, "\n\n")
	fmt.Print("State of AND is: ", gate_and.g_state, "\n\n")
	fmt.Print("State of NOT is: ", gate_not.g_state, "\n\n")
	fmt.Println("The output is: ", out_buff.get_state())
}