ALTER TABLE IF EXISTS public.rollup_action ADD COLUMN IF NOT EXISTS sender_id int8 NOT NULL DEFAULT 0;

--bun:split

COMMENT ON COLUMN public.rollup_action.sender_id IS 'Internal id of sender address';

--bun:split

with actions as (
	select signer_id, rollup_action.tx_id from rollup_action
	left join tx on tx_id = tx.id
)
update rollup_action as ra
set sender_id = actions.signer_id 
from actions
where ra.tx_id = actions.tx_id;
