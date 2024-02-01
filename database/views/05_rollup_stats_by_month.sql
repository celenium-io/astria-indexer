CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_month
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 month'::interval, rollup_stats_by_day.ts) AS ts,
		rollup_stats_by_day.rollup_id as rollup_id,
		sum(actions_count) as actions_count,
		sum(size) as size,
		min(min_size) as min_size,
		max(max_size) as max_size,
		mean(rollup(size_pct)) as avg_size,
		rollup(size_pct) as size_pct
	from rollup_stats_by_day
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('rollup_stats_by_month', INTERVAL '1 minute', INTERVAL '1 hour');
