package config

import "strconv"

// IntOrUndefined holds the value of a flag of type `int`,
// that's by default `undefined`.
// We use this instead of just `int` to differentiate `undefined`
// and `zero` values.
type IntOrUndefined struct {
	value *int
}

func (s *IntOrUndefined) Type() string {
	return "int"
}

func (s *IntOrUndefined) Value() *int {
	return s.value
}

func (s *IntOrUndefined) Set(v string) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	s.value = &i
	return nil
}

func (s *IntOrUndefined) SetNil() error {
	s.value = nil
	return nil
}

func (s *IntOrUndefined) String() string {
	if s.value == nil {
		return ""
	}
	return strconv.Itoa(*s.value)
}

func NewIntOrUndefined(v *int) IntOrUndefined {
	return IntOrUndefined{value: v}
}
