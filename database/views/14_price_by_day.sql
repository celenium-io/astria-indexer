CREATE MATERIALIZED VIEW IF NOT EXISTS price_by_day	
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 day'::interval, time) AS time,
		currency_pair,
		last(close, time) AS close,
		first(open, time) AS open,
		max(high) AS high,
		min(low) AS low
	from price_by_hour
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('price_by_day', INTERVAL '1 minute', INTERVAL '1 minute');
