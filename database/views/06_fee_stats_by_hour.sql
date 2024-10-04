CREATE MATERIALIZED VIEW IF NOT EXISTS fee_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 hour'::interval, time) AS ts,
		fee.asset as asset,
		count(*) as fee_count,
		sum(amount) as amount,
		min(amount) as min_amount,
		max(amount) as max_amount,
		avg(amount) as avg_amount,
		percentile_agg(amount) as amount_pct
	from fee
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('fee_stats_by_hour', INTERVAL '1 minute', INTERVAL '1 minute');
