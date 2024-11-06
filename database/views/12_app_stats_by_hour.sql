CREATE MATERIALIZED VIEW IF NOT EXISTS app_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 hour'::interval, rollup_action.time) AS time, 
        rollup_action.rollup_id, 
        rollup_action.sender_id, 
        sum(rollup_action.size) as size, 
        count(*) as actions_count, 
        max(rollup_action.time) as last_time,
        min(rollup_action.time) as first_time
    from rollup_action
    group by 1, 2, 3
	with no data;
        
CALL add_view_refresh_job('app_stats_by_hour', NULL, INTERVAL '1 minute');