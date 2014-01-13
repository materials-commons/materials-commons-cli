package object

import ()

type Validater interface {
	IsValid() error
}

type Ider interface {
	IdExists(string) bool
}

type Item interface {
	Validater
	Ider
	Add() error
	Delete() error
	Update() error
}
