CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 day'::interval, rollup_stats_by_hour.ts) AS ts,
		rollup_stats_by_hour.rollup_id as rollup_id,
		sum(actions_count) as actions_count,
		sum(size) as size,
		min(min_size) as min_size,
		max(max_size) as max_size,
		mean(rollup(size_pct)) as avg_size,
		rollup(size_pct) as size_pct,
		min(first_time) as first_time,
		max(last_time) as last_time
	from rollup_stats_by_hour
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('rollup_stats_by_day', INTERVAL '1 minute', INTERVAL '1 minute');
