package ingestion

import (
	"bufio"
	"context"
	"encoding/csv"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"b3-ingest/internal/infra/settings"
	"b3-ingest/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

type Service struct {
	DB  *gorm.DB
	DSN string
	Log *logger.Logger
}

func NewService(db *gorm.DB, dsn string, log *logger.Logger) *Service {
	return &Service{DB: db, DSN: dsn, Log: log}
}

func (s *Service) IngestFromCSV(dir string) error {
	s.Log.Info("Iniciando ingestão de CSVs...")
	pool, err := pgxpool.New(context.Background(), s.DSN)
	if err != nil {
		s.Log.Error("Erro ao conectar com o banco de dados: %v", err)
		return err
	}
	defer pool.Close()

	if err := s.prepareDatabase(pool); err != nil {
		return err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		s.Log.Error("Erro ao listar arquivos no diretório %s: %v", dir, err)
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, settings.GetEnvs().IngestionCores)
	var firstErr error
	var mu sync.Mutex

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(f os.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := s.processFile(f.Name(), dir, pool); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				s.Log.Error("erro processando arquivo %s: %v", f.Name(), err)
				mu.Unlock()
			}
		}(f)
	}
	wg.Wait()

	return s.finalizeIngestion(pool, firstErr)
}

func (s *Service) prepareDatabase(pool *pgxpool.Pool) error {

	s.Log.Info("Criando tabela UNLOGGED se necessário...")
	sql := `
		CREATE UNLOGGED TABLE IF NOT EXISTS tradings_unlogged (
			data_negocio date,
			codigo_instrumento text,
			preco_negocio numeric,
			quantidade_negociada bigint,
			hora_fechamento bigint,
			codigo_identificador_negocio bigint
		);
	`
	_, err := pool.Exec(context.Background(), sql)

	return err
}

func (s *Service) processFile(fileName, dir string, pool *pgxpool.Pool) error {
	path := dir + "/" + fileName
	s.Log.Info("Processando: %s", path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	r := csv.NewReader(bufio.NewReaderSize(file, 1<<20))
	r.Comma = ';'
	if _, err := r.Read(); err != nil {
		return err
	}

	copySrc := pgx.CopyFromFunc(func() ([]any, error) {
		record, err := r.Read()
		if err != nil {
			return nil, err
		}
		preco, _ := strconv.ParseFloat(strings.ReplaceAll(record[3], ",", "."), 64)
		qtd, _ := strconv.Atoi(record[4])
		dataNegocio, _ := time.Parse("2006-01-02", record[8])
		horaFechamentoInt, _ := strconv.ParseInt(record[5], 10, 64)
		return []any{dataNegocio, record[1], preco, qtd, horaFechamentoInt, record[6]}, nil
	})

	_, err = pool.CopyFrom(context.Background(), pgx.Identifier{"tradings_unlogged"},
		[]string{"data_negocio", "codigo_instrumento", "preco_negocio", "quantidade_negociada", "hora_fechamento", "codigo_identificador_negocio"}, copySrc)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	s.Log.Debug("Memória após %s: %.2f MB", fileName, float64(m.Alloc)/1024/1024)
	return err
}

func (s *Service) finalizeIngestion(pool *pgxpool.Pool, firstErr error) error {
	s.Log.Info("Finalizando ingestão e executando SQLs finais...")
	finalSQL :=
		`INSERT INTO tradings (data_negocio, codigo_instrumento, preco_negocio, quantidade_negociada, hora_fechamento, codigo_identificador_negocio)
		 SELECT data_negocio, codigo_instrumento, preco_negocio, quantidade_negociada, hora_fechamento, codigo_identificador_negocio
		 FROM tradings_unlogged ON CONFLICT DO NOTHING;

		drop table IF EXISTS tradings_unlogged;

		DO $$ BEGIN
		 IF NOT EXISTS (
		 SELECT 1 FROM pg_constraint WHERE conname = 'unique_trade_constraint'
		 AND conrelid = 'tradings'::regclass) THEN
		 ALTER TABLE tradings ADD CONSTRAINT unique_trade_constraint UNIQUE (
		 data_negocio, codigo_instrumento, hora_fechamento, codigo_identificador_negocio);
		 END IF;
		END$$;
		CREATE INDEX IF NOT EXISTS idx_tradings_ticker_data ON tradings (codigo_instrumento, data_negocio);`

	if _, err := pool.Exec(context.Background(), finalSQL); err != nil {
		s.Log.Error("ao executar SQL final: %v", err)
		return err
	}

	s.Log.Info("Ingestão finalizada.")
	return firstErr
}
