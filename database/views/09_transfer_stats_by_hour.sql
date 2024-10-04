CREATE MATERIALIZED VIEW IF NOT EXISTS transfer_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 hour'::interval, time) AS ts,
		transfer.asset as asset,
		count(*) as transfers_count,
		sum(amount) as amount
	from transfer
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('transfer_stats_by_hour', INTERVAL '1 minute', INTERVAL '1 minute');
