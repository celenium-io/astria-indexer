update price 
set height = blocks.height
from (
	select height, time from block
	where block.time >= '2025-04-01T00:00:00Z'
) as blocks
where blocks.time = price.time;