CREATE MATERIALIZED VIEW IF NOT EXISTS transfer_stats_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 day'::interval, transfer_stats_by_hour.ts) AS ts,
		transfer_stats_by_hour.asset as asset,
		sum(transfers_count) as transfers_count,
		sum(amount) as amount
	from transfer_stats_by_hour
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('transfer_stats_by_day', INTERVAL '1 minute', INTERVAL '1 minute');
