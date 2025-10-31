CREATE MATERIALIZED VIEW IF NOT EXISTS block_stats_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 day'::interval, hour.ts) AS ts,
		sum(tx_count) as tx_count,
		mean(rollup(block_time_pct)) as block_time,
		rollup(block_time_pct) as block_time_pct,
		sum(supply_change) as supply_change,
		sum(bytes_in_block) as bytes_in_block,
		sum(data_size) as data_size,
		sum(bytes_in_block)/86400.0 as bps,
		max(bps_max) as bps_max,
		min(bps_min) as bps_min,
		sum(data_size)/86400.0 as rbps,
		max(rbps_max) as rbps_max,
		min(rbps_min) as rbps_min,
		sum(tx_count)/86400.0 as tps,
		max(tps_max) as tps_max,
		min(tps_min) as tps_min
	from block_stats_by_hour as hour
	group by 1
	order by 1 desc;

CALL add_view_refresh_job('block_stats_by_day', INTERVAL '1 minute', INTERVAL '5 minute');
