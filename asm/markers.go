package asm

import "fmt"

type marker interface {
	placeholder()
}

type LABEL struct {
	Name string
}

func (l LABEL) placeholder() {
}

func (l LABEL) String() string {
	return l.Name
}

type SYMBOL struct {
	Name string
}

func (s SYMBOL) String() string {
	return fmt.Sprintf("%%%s", s.Name)
}

func (s SYMBOL) placeholder() {
}

type NUMBER struct {
	Value uint16
}

func (n NUMBER) placeholder() {
}

func (n NUMBER) String() string {
	return fmt.Sprintf("0x%X", n.Value)
}
