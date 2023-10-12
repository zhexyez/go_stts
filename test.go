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
	gate 	*Gate
}

// has type, state, array of pointers to pins that are connected
// to it, a pointer to the wire to whuich it is connected
type Gate struct {
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
	gates 		[]*Gate
	pins  		[]*Pin
	wires 		[]*Wire
	exposed_in 	[]*Pin
	exposed_out *Wire
}

func NewComponent(name string, number_of_inputs int) Component {
	component_counter++
	return Component{name: name, number: component_counter, exposed_in: make([]*Pin, number_of_inputs)}
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

	

	// declare all gates, i/o buffers are also gates
	var gate_and Gate
	var gate_not Gate
	var gate_or  Gate
	var in_buff1 Gate
	var in_buff2 Gate
	var out_buff Gate

	// declare all pins
	var pin_1	Pin
	var pin_2	Pin
	var pin_3	Pin
	var pin_4	Pin
	var pin_5	Pin
	var pin_6	Pin

	// declare all wires
	var wire_1  Wire
	var wire_2  Wire
	var wire_3	Wire
	var wire_4	Wire
	var wire_5  Wire

	// to use easily
	var state_0 = false
	var state_1 = true

	// set types and states of gates, in gate has only output, out gate has only 1 input
	in_buff1 = Gate{g_type: IN, g_state: state_0, out_pin: &wire_1}
	in_buff2 = Gate{g_type: IN, g_state: state_0, out_pin: &wire_2}
	out_buff = Gate{g_type: OUT, g_state: state_0, pins: []*Pin{&pin_6}}
	gate_and = Gate{g_type: AND, g_state: state_0, pins: []*Pin{&pin_4, &pin_5}, out_pin: &wire_5}
	gate_not = Gate{g_type: NOT, g_state: state_0, pins: []*Pin{&pin_3}, out_pin: &wire_3}
	gate_or  = Gate{g_type: OR,  g_state: state_0, pins: []*Pin{&pin_1, &pin_2}, out_pin: &wire_4}

	// set states of pins, connect to dedicated gates
	pin_1 = Pin{p_state: state_0, gate: &gate_or}
	pin_2 = Pin{p_state: state_0, gate: &gate_or}
	pin_3 = Pin{p_state: state_0, gate: &gate_not}
	pin_4 = Pin{p_state: state_0, gate: &gate_and}
	pin_5 = Pin{p_state: state_0, gate: &gate_and}
	pin_6 = Pin{p_state: state_0, gate: &out_buff}

	// set states of wires, connect to dedicated pins
	wire_1 = Wire{w_state: state_0, transmit: []*Pin{&pin_1}}
	wire_2 = Wire{w_state: state_0, transmit: []*Pin{&pin_3, &pin_5}}
	wire_3 = Wire{w_state: state_0, transmit: []*Pin{&pin_2}}
	wire_4 = Wire{w_state: state_0, transmit: []*Pin{&pin_4}}
	wire_5 = Wire{w_state: state_0, transmit: []*Pin{&pin_6}}

	// can be used multiple times
	in_buff1.update_state(state_0)
	in_buff2.update_state(state_0)

	// PRINTING
	fmt.Print("\nInputs are: ", in_buff1.g_state, " and ", in_buff2.g_state, "\n")
	fmt.Println("The output is: ", out_buff.get_state())

	// can be used multiple times
	in_buff1.update_state(state_0)
	in_buff2.update_state(state_1)

	// PRINTING
	fmt.Print("\nInputs are: ", in_buff1.g_state, " and ", in_buff2.g_state, "\n")
	fmt.Println("The output is: ", out_buff.get_state())

	// can be used multiple times
	in_buff1.update_state(state_1)
	in_buff2.update_state(state_0)

	// PRINTING
	fmt.Print("\nInputs are: ", in_buff1.g_state, " and ", in_buff2.g_state, "\n")
	fmt.Println("The output is: ", out_buff.get_state())

	// can be used multiple times
	in_buff1.update_state(state_1)
	in_buff2.update_state(state_1)

	// PRINTING
	fmt.Print("\nInputs are: ", in_buff1.g_state, " and ", in_buff2.g_state, "\n")
	fmt.Println("The output is: ", out_buff.get_state())
}