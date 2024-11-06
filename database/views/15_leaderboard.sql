CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard AS
   with board as (
	select 
            app_id,
            sum(size) as size, 
            sum(actions_count) as actions_count, 
            max(last_time) as last_time, 
            min(first_time) as first_time
        from (
            select
                rollup_id, 
                sender_id,
                sum(size) as size, 
                sum(actions_count) as actions_count, 
                max(last_time) as last_time, 
                min(first_time) as first_time
            from app_stats_by_month
            group by 1, 2
        ) as agg
        inner join app_id on app_id.address_id = agg.sender_id AND (app_id.rollup_id = agg.rollup_id OR app_id.rollup_id = 0)
        group by 1
    ) 
    select 
        board.size, 
        board.actions_count, 
        board.last_time, 
        board.first_time,
        board.size / (select sum(size) from board) as size_pct,
        board.actions_count / (select sum(actions_count) from board)as actions_count_pct,
        app.*
    from board
    inner join app on app.id = board.app_id;

CALL add_job_refresh_materialized_view();