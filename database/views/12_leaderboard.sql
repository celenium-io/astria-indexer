CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard AS
	select 
        app.*,
        agg.size,
        agg.min_size,
        agg.max_size,
        agg.avg_size,
        agg.actions_count,
        agg.last_time,
        agg.first_time
    from (
        select
            rollup_id,
            sum(size) as size, 
            min(min_size) as min_size,
            max(max_size) as max_size,
            mean(rollup(size_pct)) as avg_size,
            sum(actions_count) as actions_count, 
            max(last_time) as last_time, 
            min(first_time) as first_time
        from rollup_stats_by_month
        group by rollup_id
    ) as agg
    inner join app on app.rollup_id = agg.rollup_id;

CALL add_job_refresh_materialized_view();