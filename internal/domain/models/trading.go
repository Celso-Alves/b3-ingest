package models

import "time"

type Trading struct {
	DataNegocio                time.Time
	CodigoInstrumento          string
	PrecoNegocio               float64
	QuantidadeNegociada        int64
	HoraFechamento             int64
	HashArquivo                string
	CodigoIdentificadorNegocio int
}
