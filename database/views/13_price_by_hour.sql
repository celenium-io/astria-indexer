CREATE MATERIALIZED VIEW IF NOT EXISTS price_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 hour'::interval, time) AS time,
		currency_pair,
		last(price, time) / pow(10, last(coalesce(pair.decimals, 0), time)) AS close,
		first(price, time) / pow(10, last(coalesce(pair.decimals, 0), time)) AS open,
		max(price) / pow(10, last(coalesce(pair.decimals, 0), time)) AS high,
		min(price) / pow(10, last(coalesce(pair.decimals, 0), time)) AS low
	from price
	left join lateral (select * from market where market.pair = price.currency_pair and market.updated_at <= price.time order by updated_at desc limit 1) pair on true
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('price_by_hour', INTERVAL '1 minute', INTERVAL '1 minute');