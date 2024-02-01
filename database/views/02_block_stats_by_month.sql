CREATE MATERIALIZED VIEW IF NOT EXISTS block_stats_by_month
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 month', day.ts) AS ts,
		sum(tx_count) as tx_count,
		mean(rollup(block_time_pct)) as block_time,
		rollup(block_time_pct) as block_time_pct,
		sum(gas_wanted) as gas_wanted,
		sum(gas_used) as gas_used,
		(case when sum(gas_wanted) > 0 then sum(fee) / sum(gas_wanted) else 0 end) as gas_price,
		(case when sum(gas_wanted) > 0 then sum(gas_used) / sum(gas_wanted) else 0 end) as gas_efficiency,
		sum(fee) as fee,
		sum(supply_change) as supply_change,
		sum(bytes_in_block) as bytes_in_block,
		sum(data_size) as data_size,
		sum(bytes_in_block)/(count(*) * 86400.0) as bps,
		max(bps_max) as bps_max,
		min(bps_min) as bps_min,
		sum(data_size)/(count(*) * 86400.0) as rbps,
		max(rbps_max) as rbps_max,
		min(rbps_min) as rbps_min,
		sum(tx_count)/(count(*) * 86400.0) as tps,
		max(tps_max) as tps_max,
		min(tps_min) as tps_min	
	from block_stats_by_day as day
	group by 1
	order by 1 desc;

CALL add_view_refresh_job('block_stats_by_month', INTERVAL '1 minute', INTERVAL '1 hour');