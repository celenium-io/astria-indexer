CREATE MATERIALIZED VIEW IF NOT EXISTS fee_stats_by_month
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 month'::interval, fee_stats_by_day.ts) AS ts,
		fee_stats_by_day.asset as asset,
		sum(fee_count) as fee_count,
		sum(amount) as amount,
		min(min_amount) as min_amount,
		max(max_amount) as max_amount,
		mean(rollup(amount_pct)) as avg_amount,
		rollup(amount_pct) as amount_pct
	from fee_stats_by_day
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('fee_stats_by_month', INTERVAL '1 minute', INTERVAL '1 hour');
