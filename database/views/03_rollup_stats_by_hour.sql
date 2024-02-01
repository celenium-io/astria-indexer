CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 hour'::interval, time) AS ts,
		rollup_action.rollup_id as rollup_id,
		count(*) as actions_count,
		sum(size) as size,
		min(size) as min_size,
		max(size) as max_size,
		avg(size) as avg_size,
		percentile_agg(size) as size_pct
	from rollup_action
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('rollup_stats_by_hour', INTERVAL '1 minute', INTERVAL '1 minute');
