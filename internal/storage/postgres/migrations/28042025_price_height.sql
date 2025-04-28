update price 
set height = blocks.height
from (
	select block.height, price.time from price
	left join block on block.time = price.time
) as blocks;
