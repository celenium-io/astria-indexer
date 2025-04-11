CREATE MATERIALIZED VIEW IF NOT EXISTS price_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 hour'::interval, time) AS time,
		currency_pair,
		last(price, time) AS close,
		first(price, time) AS open,
		max(price) AS high,
		min(price) AS low
	from price
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('price_by_hour', INTERVAL '1 minute', INTERVAL '1 minute');
