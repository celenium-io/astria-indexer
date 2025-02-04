// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"fmt"

	internalAstria "github.com/celenium-io/astria-indexer/internal/astria"
	"github.com/celenium-io/astria-indexer/internal/storage"
	clientTypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	connectionTypes "github.com/cosmos/ibc-go/v9/modules/core/03-connection/types"
	channelTypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
	lightTypes "github.com/cosmos/ibc-go/v9/modules/light-clients/07-tendermint"
	"github.com/fatih/structs"
	"github.com/goccy/go-json"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type IbcTransfer struct {
	Amount   decimal.Decimal `json:"amount"`
	Denom    string          `json:"denom"`
	Receiver string          `json:"receiver"`
	Sender   string          `json:"sender"`
}

func parseIbcMessages(typ string, data []byte, action *storage.Action, ctx *Context) error {
	switch typ {

	// channel messages
	case "/ibc.core.channel.v1.MsgRecvPacket":
		var msg channelTypes.MsgRecvPacket
		if err := handleType(&msg, data, action); err != nil {
			return err
		}

		var transfer IbcTransfer
		if err := json.Unmarshal(msg.Packet.Data, &transfer); err != nil {
			return nil
		}
		asset := fmt.Sprintf("%s/%s/%s", msg.Packet.GetDestPort(), msg.Packet.GetDestChannel(), transfer.Denom)
		if err := handleTransfer(transfer, action, ctx, asset, false); err != nil {
			return errors.Wrap(err, "transfer handling")
		}
	case "/ibc.core.channel.v1.MsgTimeout":
		var msg channelTypes.MsgTimeout
		if err := handleType(&msg, data, action); err != nil {
			return err
		}

		var transfer IbcTransfer
		if err := json.Unmarshal(msg.Packet.Data, &transfer); err != nil {
			return nil
		}
		if err := handleTransfer(transfer, action, ctx, transfer.Denom, true); err != nil {
			return errors.Wrap(err, "transfer handling")
		}
	case "/ibc.core.channel.v1.MsgChannelOpenInit":
		var msg channelTypes.MsgChannelOpenInit
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.channel.v1.MsgChannelOpenTry":
		var msg channelTypes.MsgChannelOpenTry
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.channel.v1.MsgChannelOpenAck":
		var msg channelTypes.MsgChannelOpenAck
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.channel.v1.MsgChannelOpenConfirm":
		var msg channelTypes.MsgChannelOpenConfirm
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.channel.v1.MsgChannelCloseInit":
		var msg channelTypes.MsgChannelCloseInit
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.channel.v1.MsgChannelCloseConfirm":
		var msg channelTypes.MsgChannelCloseConfirm
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.channel.v1.MsgAcknowledgement":
		var msg channelTypes.MsgAcknowledgement
		if err := handleType(&msg, data, action); err != nil {
			return err
		}

	// connection messages
	case "/ibc.core.connection.v1.MsgConnectionOpenInit":
		var msg connectionTypes.MsgConnectionOpenInit
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.connection.v1.MsgConnectionOpenTry":
		var msg connectionTypes.MsgConnectionOpenTry
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.connection.v1.MsgConnectionOpenAck":
		var msg connectionTypes.MsgConnectionOpenAck
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.connection.v1.MsgConnectionOpenConfirm":
		var msg connectionTypes.MsgConnectionOpenConfirm
		if err := handleType(&msg, data, action); err != nil {
			return err
		}

	// clients messages
	case "/ibc.core.client.v1.MsgCreateClient":
		var msg clientTypes.MsgCreateClient
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.client.v1.MsgUpdateClient":
		var msg clientTypes.MsgUpdateClient
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	case "/ibc.core.client.v1.MsgUpgradeClient":
		var msg clientTypes.MsgUpgradeClient
		if err := handleType(&msg, data, action); err != nil {
			return err
		}
	// case "/ibc.core.client.v1.MsgSubmitMisbehaviour":
	// 	var msg clientTypes.MsgSubmitMisbehaviour
	// 	if err := handleType(&msg, data, action); err != nil {
	// 		return err
	// 	}
	default:
	}
	return nil
}

func handleType[T proto.Message](msg T, data []byte, action *storage.Action) error {
	if err := proto.Unmarshal(data, msg); err != nil {
		return errors.Wrap(err, "unmarshal IBC message")
	}
	m := structs.Map(msg)

	for key, value := range m {
		if valMap, ok := value.(map[string]any); ok {
			if typeUrl, ok := valMap["TypeUrl"]; ok {
				switch typeUrl {
				case "/ibc.lightclients.tendermint.v1.ClientState":
					if data, ok := valMap["Value"]; ok {
						var csMsg lightTypes.ClientState
						if err := proto.Unmarshal(data.([]byte), &csMsg); err != nil {
							return errors.Wrap(err, "unmarshal client state")
						}
						m[key] = structs.Map(csMsg)
					}
				case "/ibc.lightclients.tendermint.v1.ConsensusState":
					if data, ok := valMap["Value"]; ok {
						var csMsg lightTypes.ConsensusState
						if err := proto.Unmarshal(data.([]byte), &csMsg); err != nil {
							return errors.Wrap(err, "unmarshal consensus state")
						}
						m[key] = structs.Map(csMsg)
					}
				case "/ibc.lightclients.tendermint.v1.Header":
					if data, ok := valMap["Value"]; ok {
						var csMsg lightTypes.Header
						if err := proto.Unmarshal(data.([]byte), &csMsg); err != nil {
							return errors.Wrap(err, "unmarshal header")
						}
						m[key] = structs.Map(csMsg)
					}
				}
			}
		}

	}

	action.Data["msg"] = m
	delete(action.Data, "raw")
	return nil
}

func handleTransfer(transfer IbcTransfer, action *storage.Action, ctx *Context, asset string, isRefund bool) error {
	var addr string
	var amount = transfer.Amount.Copy()

	switch {
	case internalAstria.IsAddress(transfer.Receiver):
		addr = transfer.Receiver
		if isRefund {
			amount = amount.Neg()
		}
	case internalAstria.IsCompatAddress(transfer.Receiver):
		a, err := internalAstria.CompatToAstria(transfer.Receiver)
		if err != nil {
			return err
		}
		addr = a
		if isRefund {
			amount = amount.Neg()
		}
	case internalAstria.IsAddress(transfer.Sender):
		addr = transfer.Sender
		if !isRefund {
			amount = amount.Neg()
		}
	case internalAstria.IsCompatAddress(transfer.Sender):
		a, err := internalAstria.CompatToAstria(transfer.Sender)
		if err != nil {
			return err
		}
		addr = a
		if !isRefund {
			amount = amount.Neg()
		}
	}

	if addr != "" {
		address := ctx.Addresses.Set(addr, action.Height, amount, asset, 0, 0)
		action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
			Address:  address,
			Height:   action.Height,
			Currency: asset,
			Update:   amount,
		})
	}
	return nil
}
