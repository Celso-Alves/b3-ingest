package models

import (
	"time"

	"b3-ingest/internal/domain/models"
)

type Trading struct {
	DataNegocio                time.Time `gorm:"column:data_negocio;type:date"`
	CodigoInstrumento          string    `gorm:"column:codigo_instrumento"`
	PrecoNegocio               float64   `gorm:"column:preco_negocio"`
	QuantidadeNegociada        int64     `gorm:"column:quantidade_negociada"`
	HoraFechamento             int64     `gorm:"column:hora_fechamento"`
	HashArquivo                string    `gorm:"column:hash_arquivo"`
	CodigoIdentificadorNegocio int       `gorm:"column:codigo_identificador_negocio"`
}

// ToTradingORMModel converts a domain Trading model to an ORM Trading model.
func ToTradingORMModel(domain models.Trading) Trading {
	return Trading{
		DataNegocio:                domain.DataNegocio,
		CodigoInstrumento:          domain.CodigoInstrumento,
		PrecoNegocio:               domain.PrecoNegocio,
		QuantidadeNegociada:        domain.QuantidadeNegociada,
		HoraFechamento:             domain.HoraFechamento,
		HashArquivo:                domain.HashArquivo,
		CodigoIdentificadorNegocio: domain.CodigoIdentificadorNegocio,
	}
}

// ToTradingDomainModel converts an ORM Trading model to a domain Trading model.
func ToTradingDomainModel(orm Trading) models.Trading {
	return models.Trading{
		DataNegocio:                orm.DataNegocio,
		CodigoInstrumento:          orm.CodigoInstrumento,
		PrecoNegocio:               orm.PrecoNegocio,
		QuantidadeNegociada:        orm.QuantidadeNegociada,
		HoraFechamento:             orm.HoraFechamento,
		HashArquivo:                orm.HashArquivo,
		CodigoIdentificadorNegocio: orm.CodigoIdentificadorNegocio,
	}
}
