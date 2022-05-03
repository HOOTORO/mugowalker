package mapman

import (
	"worker/esperia"
)

type MapExporter interface {
	Export() error
}

type MapImporter interface {
	Import(interface{}) esperia.Esperia
}

func He() {
	m := 1 + 1
	_ = m
}
