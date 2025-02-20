// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package genesis

import (
	"strconv"
	"strings"

	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	nodeTypes "github.com/celenium-io/astria-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
)

func (module *Module) parseConstants(appState nodeTypes.AppState, consensus pkgTypes.ConsensusParams, data *parsedData) {
	// block
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameBlock,
		Name:   "block_max_bytes",
		Value:  strconv.FormatInt(consensus.Block.MaxBytes, 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameBlock,
		Name:   "block_max_gas",
		Value:  strconv.FormatInt(consensus.Block.MaxGas, 10),
	})

	// evidence
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameEvidence,
		Name:   "max_age_num_blocks",
		Value:  strconv.FormatInt(consensus.Evidence.MaxAgeNumBlocks, 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameEvidence,
		Name:   "max_age_duration",
		Value:  consensus.Evidence.MaxAgeDuration.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameEvidence,
		Name:   "max_bytes",
		Value:  strconv.FormatInt(consensus.Evidence.MaxBytes, 10),
	})

	// validator
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameValidator,
		Name:   "pub_key_types",
		Value:  strings.Join(consensus.Validator.PubKeyTypes, ","),
	})

	// version
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameVersion,
		Name:   "app",
		Value:  strconv.FormatUint(consensus.Version.AppVersion, 10),
	})

	// generic
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "authority_sudo_address",
		Value:  appState.AuthoritySudoAddress.Value,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "native_asset_base_denomination",
		Value:  appState.NativeAssetBaseDenomination,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ibc_sudo_address",
		Value:  appState.IbcSudoAddress.Value,
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "bridge_lock_base",
		Value:  appState.Fees.BridgeLock.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "bridge_lock_multiplier",
		Value:  appState.Fees.BridgeLock.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "bridge_unlock_base",
		Value:  appState.Fees.BridgeUnlock.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "bridge_unlock_multiplier",
		Value:  appState.Fees.BridgeUnlock.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "transfer_base",
		Value:  appState.Fees.Transfer.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "transfer_multiplier",
		Value:  appState.Fees.Transfer.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "rollup_data_submission_base",
		Value:  appState.Fees.RollupDataSubmission.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "rollup_data_submission_multiplier",
		Value:  appState.Fees.RollupDataSubmission.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ics20_withdrawal_base",
		Value:  appState.Fees.Ics20Withdrawal.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ics20_withdrawal_multiplier",
		Value:  appState.Fees.Ics20Withdrawal.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "init_bridge_account_base",
		Value:  appState.Fees.InitBridgeAccount.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "init_bridge_account_multiplier",
		Value:  appState.Fees.InitBridgeAccount.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "bridge_sudo_change_base",
		Value:  appState.Fees.BridgeSudoChange.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "bridge_sudo_change_multiplier",
		Value:  appState.Fees.BridgeSudoChange.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ibc_relay_base",
		Value:  appState.Fees.IbcRelay.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ibc_relay_multiplier",
		Value:  appState.Fees.IbcRelay.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "validator_update_base",
		Value:  appState.Fees.ValidatorUpdate.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "validator_update_multiplier",
		Value:  appState.Fees.ValidatorUpdate.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "fee_asset_change_base",
		Value:  appState.Fees.FeeAssetChange.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "fee_asset_change_multiplier",
		Value:  appState.Fees.FeeAssetChange.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "fee_change_base",
		Value:  appState.Fees.FeeChange.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "fee_change_multiplier",
		Value:  appState.Fees.FeeChange.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ibc_relayer_change_base",
		Value:  appState.Fees.IbcRelayerChange.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ibc_relayer_change_multiplier",
		Value:  appState.Fees.IbcRelayerChange.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "sudo_address_change_base",
		Value:  appState.Fees.SudoAddressChange.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "sudo_address_change_multiplier",
		Value:  appState.Fees.SudoAddressChange.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ibc_sudo_change_base",
		Value:  appState.Fees.IbcSudoChange.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ibc_sudo_change_multiplier",
		Value:  appState.Fees.IbcSudoChange.Multiplier.String(),
	})

	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "bridge_transfer_base",
		Value:  appState.Fees.BridgeTransfer.Base.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "bridge_transfer_multiplier",
		Value:  appState.Fees.BridgeTransfer.Multiplier.String(),
	})
}
