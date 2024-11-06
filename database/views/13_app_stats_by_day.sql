CREATE MATERIALIZED VIEW IF NOT EXISTS app_stats_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 day'::interval, actions.time) AS time, 
        actions.rollup_id, 
        actions.sender_id, 
        sum(actions.size) as size, 
        sum(actions.actions_count) as actions_count, 
        max(actions.last_time) as last_time,
        min(actions.first_time) as first_time
    from app_stats_by_hour as actions
    group by 1, 2, 3
	with no data;
        
CALL add_view_refresh_job('app_stats_by_day', NULL, INTERVAL '5 minute');