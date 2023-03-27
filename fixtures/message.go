package fixtures

import (
	"errors"
	"fmt"
)

// Command is an implementation of dogma.Command used for testing.
type Command[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (c Command[T]) MessageDescription() string {
	return fmt.Sprintf("%v", c)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (c Command[T]) Validate() error {
	if c.Invalid != "" {
		return errors.New(c.Invalid)
	}
	return nil
}

// Event is an implementation of dogma.Event used for testing.
type Event[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (e Event[T]) MessageDescription() string {
	return fmt.Sprintf("%v", e)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (e Event[T]) Validate() error {
	if e.Invalid != "" {
		return errors.New(e.Invalid)
	}
	return nil
}

// Timeout is an implementation of dogma.Timeout used for testing.
type Timeout[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (t Timeout[T]) MessageDescription() string {
	return fmt.Sprintf("%v", t)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (t Timeout[T]) Validate() error {
	if t.Invalid != "" {
		return errors.New(t.Invalid)
	}
	return nil
}

// MessageA is type used as a dogma.Message in tests.
type MessageA struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageA) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageA1 is an instance of MessageA with a distinct value.
	MessageA1 = MessageA{"A1"}
	// MessageA2 is an instance of MessageA with a distinct value.
	MessageA2 = MessageA{"A2"}
	// MessageA3 is an instance of MessageA with a distinct value.
	MessageA3 = MessageA{"A3"}
)

// MessageB is type used as a dogma.Message in tests.
type MessageB struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageB) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageB1 is an instance of MessageB with a distinct value.
	MessageB1 = MessageB{"B1"}
	// MessageB2 is an instance of MessageB with a distinct value.
	MessageB2 = MessageB{"B2"}
	// MessageB3 is an instance of MessageB with a distinct value.
	MessageB3 = MessageB{"B3"}
)

// MessageC is type used as a dogma.Message in tests.
type MessageC struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageC) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageC1 is an instance of MessageC with a distinct value.
	MessageC1 = MessageC{"C1"}
	// MessageC2 is an instance of MessageC with a distinct value.
	MessageC2 = MessageC{"C2"}
	// MessageC3 is an instance of MessageC with a distinct value.
	MessageC3 = MessageC{"C3"}
)

// MessageD is type used as a dogma.Message in tests.
type MessageD struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageD) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageD1 is an instance of MessageD with a distinct value.
	MessageD1 = MessageD{"D1"}
	// MessageD2 is an instance of MessageD with a distinct value.
	MessageD2 = MessageD{"D2"}
	// MessageD3 is an instance of MessageD with a distinct value.
	MessageD3 = MessageD{"D3"}
)

// MessageE is type used as a dogma.Message in tests.
type MessageE struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageE) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageE1 is an instance of MessageE with a distinct value.
	MessageE1 = MessageE{"E1"}
	// MessageE2 is an instance of MessageE with a distinct value.
	MessageE2 = MessageE{"E2"}
	// MessageE3 is an instance of MessageE with a distinct value.
	MessageE3 = MessageE{"E3"}
)

// MessageF is type used as a dogma.Message in tests.
type MessageF struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageF) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageF1 is an instance of MessageF with a distinct value.
	MessageF1 = MessageF{"F1"}
	// MessageF2 is an instance of MessageF with a distinct value.
	MessageF2 = MessageF{"F2"}
	// MessageF3 is an instance of MessageF with a distinct value.
	MessageF3 = MessageF{"F3"}
)

// MessageG is type used as a dogma.Message in tests.
type MessageG struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageG) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageG1 is an instance of MessageG with a distinct value.
	MessageG1 = MessageG{"G1"}
	// MessageG2 is an instance of MessageG with a distinct value.
	MessageG2 = MessageG{"G2"}
	// MessageG3 is an instance of MessageG with a distinct value.
	MessageG3 = MessageG{"G3"}
)

// MessageH is type used as a dogma.Message in tests.
type MessageH struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageH) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageH1 is an instance of MessageH with a distinct value.
	MessageH1 = MessageH{"H1"}
	// MessageH2 is an instance of MessageH with a distinct value.
	MessageH2 = MessageH{"H2"}
	// MessageH3 is an instance of MessageH with a distinct value.
	MessageH3 = MessageH{"H3"}
)

// MessageI is type used as a dogma.Message in tests.
type MessageI struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageI) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageI1 is an instance of MessageI with a distinct value.
	MessageI1 = MessageI{"I1"}
	// MessageI2 is an instance of MessageI with a distinct value.
	MessageI2 = MessageI{"I2"}
	// MessageI3 is an instance of MessageI with a distinct value.
	MessageI3 = MessageI{"I3"}
)

// MessageJ is type used as a dogma.Message in tests.
type MessageJ struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageJ) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageJ1 is an instance of MessageJ with a distinct value.
	MessageJ1 = MessageJ{"J1"}
	// MessageJ2 is an instance of MessageJ with a distinct value.
	MessageJ2 = MessageJ{"J2"}
	// MessageJ3 is an instance of MessageJ with a distinct value.
	MessageJ3 = MessageJ{"J3"}
)

// MessageK is type used as a dogma.Message in tests.
type MessageK struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageK) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageK1 is an instance of MessageK with a distinct value.
	MessageK1 = MessageK{"K1"}
	// MessageK2 is an instance of MessageK with a distinct value.
	MessageK2 = MessageK{"K2"}
	// MessageK3 is an instance of MessageK with a distinct value.
	MessageK3 = MessageK{"K3"}
)

// MessageL is type used as a dogma.Message in tests.
type MessageL struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageL) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageL1 is an instance of MessageL with a distinct value.
	MessageL1 = MessageL{"L1"}
	// MessageL2 is an instance of MessageL with a distinct value.
	MessageL2 = MessageL{"L2"}
	// MessageL3 is an instance of MessageL with a distinct value.
	MessageL3 = MessageL{"L3"}
)

// MessageM is type used as a dogma.Message in tests.
type MessageM struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageM) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageM1 is an instance of MessageM with a distinct value.
	MessageM1 = MessageM{"M1"}
	// MessageM2 is an instance of MessageM with a distinct value.
	MessageM2 = MessageM{"M2"}
	// MessageM3 is an instance of MessageM with a distinct value.
	MessageM3 = MessageM{"M3"}
)

// MessageN is type used as a dogma.Message in tests.
type MessageN struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageN) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageN1 is an instance of MessageN with a distinct value.
	MessageN1 = MessageN{"N1"}
	// MessageN2 is an instance of MessageN with a distinct value.
	MessageN2 = MessageN{"N2"}
	// MessageN3 is an instance of MessageN with a distinct value.
	MessageN3 = MessageN{"N3"}
)

// MessageO is type used as a dogma.Message in tests.
type MessageO struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageO) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageO1 is an instance of MessageO with a distinct value.
	MessageO1 = MessageO{"O1"}
	// MessageO2 is an instance of MessageO with a distinct value.
	MessageO2 = MessageO{"O2"}
	// MessageO3 is an instance of MessageO with a distinct value.
	MessageO3 = MessageO{"O3"}
)

// MessageP is type used as a dogma.Message in tests.
type MessageP struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageP) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageP1 is an instance of MessageP with a distinct value.
	MessageP1 = MessageP{"P1"}
	// MessageP2 is an instance of MessageP with a distinct value.
	MessageP2 = MessageP{"P2"}
	// MessageP3 is an instance of MessageP with a distinct value.
	MessageP3 = MessageP{"P3"}
)

// MessageQ is type used as a dogma.Message in tests.
type MessageQ struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageQ) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageQ1 is an instance of MessageQ with a distinct value.
	MessageQ1 = MessageQ{"Q1"}
	// MessageQ2 is an instance of MessageQ with a distinct value.
	MessageQ2 = MessageQ{"Q2"}
	// MessageQ3 is an instance of MessageQ with a distinct value.
	MessageQ3 = MessageQ{"Q3"}
)

// MessageR is type used as a dogma.Message in tests.
type MessageR struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageR) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageR1 is an instance of MessageR with a distinct value.
	MessageR1 = MessageR{"R1"}
	// MessageR2 is an instance of MessageR with a distinct value.
	MessageR2 = MessageR{"R2"}
	// MessageR3 is an instance of MessageR with a distinct value.
	MessageR3 = MessageR{"R3"}
)

// MessageS is type used as a dogma.Message in tests.
type MessageS struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageS) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageS1 is an instance of MessageS with a distinct value.
	MessageS1 = MessageS{"S1"}
	// MessageS2 is an instance of MessageS with a distinct value.
	MessageS2 = MessageS{"S2"}
	// MessageS3 is an instance of MessageS with a distinct value.
	MessageS3 = MessageS{"S3"}
)

// MessageT is type used as a dogma.Message in tests.
type MessageT struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageT) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageT1 is an instance of MessageT with a distinct value.
	MessageT1 = MessageT{"T1"}
	// MessageT2 is an instance of MessageT with a distinct value.
	MessageT2 = MessageT{"T2"}
	// MessageT3 is an instance of MessageT with a distinct value.
	MessageT3 = MessageT{"T3"}
)

// MessageU is type used as a dogma.Message in tests.
type MessageU struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageU) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageU1 is an instance of MessageU with a distinct value.
	MessageU1 = MessageU{"U1"}
	// MessageU2 is an instance of MessageU with a distinct value.
	MessageU2 = MessageU{"U2"}
	// MessageU3 is an instance of MessageU with a distinct value.
	MessageU3 = MessageU{"U3"}
)

// MessageV is type used as a dogma.Message in tests.
type MessageV struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageV) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageV1 is an instance of MessageV with a distinct value.
	MessageV1 = MessageV{"V1"}
	// MessageV2 is an instance of MessageV with a distinct value.
	MessageV2 = MessageV{"V2"}
	// MessageV3 is an instance of MessageV with a distinct value.
	MessageV3 = MessageV{"V3"}
)

// MessageW is type used as a dogma.Message in tests.
type MessageW struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageW) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageW1 is an instance of MessageW with a distinct value.
	MessageW1 = MessageW{"W1"}
	// MessageW2 is an instance of MessageW with a distinct value.
	MessageW2 = MessageW{"W2"}
	// MessageW3 is an instance of MessageW with a distinct value.
	MessageW3 = MessageW{"W3"}
)

// MessageX is type used as a dogma.Message in tests.
type MessageX struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageX) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageX1 is an instance of MessageX with a distinct value.
	MessageX1 = MessageX{"X1"}
	// MessageX2 is an instance of MessageX with a distinct value.
	MessageX2 = MessageX{"X2"}
	// MessageX3 is an instance of MessageX with a distinct value.
	MessageX3 = MessageX{"X3"}
)

// MessageY is type used as a dogma.Message in tests.
type MessageY struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageY) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageY1 is an instance of MessageY with a distinct value.
	MessageY1 = MessageY{"Y1"}
	// MessageY2 is an instance of MessageY with a distinct value.
	MessageY2 = MessageY{"Y2"}
	// MessageY3 is an instance of MessageY with a distinct value.
	MessageY3 = MessageY{"Y3"}
)

// MessageZ is type used as a dogma.Message in tests.
type MessageZ struct {
	Value interface{}
}

// Validate returns m.Value if it is an error.
func (m MessageZ) Validate() error {
	err, _ := m.Value.(error)
	return err
}

var (
	// MessageZ1 is an instance of MessageZ with a distinct value.
	MessageZ1 = MessageZ{"Z1"}
	// MessageZ2 is an instance of MessageZ with a distinct value.
	MessageZ2 = MessageZ{"Z2"}
	// MessageZ3 is an instance of MessageZ with a distinct value.
	MessageZ3 = MessageZ{"Z3"}
)
