CREATE MATERIALIZED VIEW IF NOT EXISTS block_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 hour'::interval, time) AS ts,
		sum(tx_count) as tx_count,
		avg(block_time) as block_time,
		percentile_agg(block_time) as block_time_pct,
		sum(supply_change) as supply_change,
		sum(fee) as fee,
		sum(bytes_in_block) as bytes_in_block,
		sum(data_size) as data_size,
		(sum(bytes_in_block)/3600.0) as bps,
		max(case when block_time > 0 then bytes_in_block::float/(block_time/1000.0) else 0 end) as bps_max,
		min(case when block_time > 0 then bytes_in_block::float/(block_time/1000.0) else 0 end) as bps_min,
		(sum(data_size)/3600.0) as rbps,
		max(case when block_time > 0 then data_size::float/(block_time/1000.0) else 0 end) as rbps_max,
		min(case when block_time > 0 then data_size::float/(block_time/1000.0) else 0 end) as rbps_min,
		(sum(tx_count)/3600.0) as tps,
		max(case when block_time > 0 then tx_count::float/(block_time/1000.0) else 0 end) as tps_max,
		min(case when block_time > 0 then tx_count::float/(block_time/1000.0) else 0 end) as tps_min
	from block_stats
	group by 1
	order by 1 desc;

CALL add_view_refresh_job('block_stats_by_hour', INTERVAL '1 minute', INTERVAL '1 minute');
