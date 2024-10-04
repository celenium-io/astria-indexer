CREATE MATERIALIZED VIEW IF NOT EXISTS fee_stats_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 day'::interval, fee_stats_by_hour.ts) AS ts,
		fee_stats_by_hour.asset as asset,
		sum(fee_count) as fee_count,
		sum(amount) as amount,
		min(min_amount) as min_amount,
		max(max_amount) as max_amount,
		mean(rollup(amount_pct)) as avg_amount,
		rollup(amount_pct) as amount_pct
	from fee_stats_by_hour
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('fee_stats_by_day', INTERVAL '1 minute', INTERVAL '1 minute');
