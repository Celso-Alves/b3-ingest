--truncate table public.tradings
TRUNCATE table public.tradings_unlogged;
[b3db].public.tradings.[tradings_pkey]

--papeis que mais 
--WINQ25	28793638
--WDOQ25	4688104
--BOVA11	526864
--BITN25	461871


SELECT setting
FROM pg_settings
WHERE name = 'max_wal_size';


SELECT MAX(soma) FROM (
				SELECT data_negocio, SUM(quantidade_negociada) soma
				FROM tradings
				WHERE codigo_instrumento = 'INDQ25' AND data_negocio >= '2025-07-25'
				GROUP BY data_negocio) AS subquery;

				SELECT codigo_instrumento, count(*) 
				FROM tradings
				
				GROUP BY codigo_instrumento



ALTER TABLE tradings
  ADD COLUMN hora_fechamento TIME to_timestamp(hora_fechamento, 'HH24MISSMS')::time;
ALTER TABLE tradings DROP COLUMN IF EXISTS hora_fechamento;


ALTER TABLE tradings
  
  ALTER COLUMN hora_fechamento TYPE TIME USING to_timestamp(hora_fechamento, 'HH24MISSMS')::time;

-- (Opcional) Para garantir idempotência:
ALTER TABLE tradings
  ADD CONSTRAINT unique_trade_with_id_negocio UNIQUE (
    data_negocio,
    codigo_instrumento,
    preco_negocio,
    quantidade_negociada,
    hora_fechamento,
    codigo_identificador_negocio
  );

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'unique_trade_with_id_negocio'
            AND conrelid = 'tradings'::regclass
    ) THEN
        ALTER TABLE tradings
        ADD CONSTRAINT unique_trade_with_id_negocio UNIQUE (
            data_negocio,
            codigo_instrumento,
            preco_negocio,
            quantidade_negociada,
            hora_fechamento,
            codigo_identificador_negocio
        );
    END IF;
END$$;


ALTER TABLE tradings
ADD CONSTRAINT unique_trade_constraint
UNIQUE (
  data_negocio,
  codigo_instrumento,
  --preco_negocio,
  --quantidade_negociada,
  hora_fechamento,
  codigo_identificador_negocio
);
ALTER TABLE tradings
DROP CONSTRAINT IF EXISTS unique_trade_constraint;

CREATE INDEX IF NOT EXISTS idx_tradings_ticker_data
ON tradings (codigo_instrumento, data_negocio);


DROP INDEX IF EXISTS idx_tradings_ticker_data;

ALTER TABLE tradings ADD COLUMN id TEXT PRIMARY KEY;


ALTER TABLE tradings DROP COLUMN hora_fechamento;

ALTER TABLE tradings ADD COLUMN hora_fechamento BIGINT;
