ALTER TABLE public."bridge" ADD IF NOT EXISTS disable_deposits bool DEFAULT false NOT NULL;

--bun:split

ALTER TABLE public."bridge" ALTER COLUMN disable_deposits SET STORAGE PLAIN;

--bun:split

COMMENT ON COLUMN public."bridge".disable_deposits IS 'Disable deposits to the bridge account';
