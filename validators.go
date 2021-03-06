package atmin

import (
	"bytes"
)

type ExactValidator struct{}

func (v ExactValidator) Validate(initial, current []byte) bool {
	return bytes.Equal(initial, current)
}

func (m Minimizer) ValidateExact() Minimizer {
	m.val = ExactValidator{}
	return m
}

type StringValidator struct {
	needle []byte
}

func (v StringValidator) Validate(initial, current []byte) bool {
	return bytes.Contains(current, v.needle)
}

func (m Minimizer) ValidateString(needle string) Minimizer {
	m.val = StringValidator{needle: []byte(needle)}
	return m
}
