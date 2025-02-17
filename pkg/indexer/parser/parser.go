// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"

	"github.com/celenium-io/astria-indexer/pkg/node"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

type Module struct {
	modules.BaseModule

	api          node.Api
	bridgeAssets map[string]string
}

var _ modules.Module = (*Module)(nil)

const (
	InputName  = "blocks"
	OutputName = "data"
	StopOutput = "stop"
)

func NewModule(api node.Api) Module {
	m := Module{
		BaseModule: modules.New("parser"),
		api:        api,
	}
	m.CreateInput(InputName)
	m.CreateOutput(OutputName)
	m.CreateOutput(StopOutput)

	return m
}

func (p *Module) Init(ctx context.Context, bridgeAssets map[string]string) {
	p.bridgeAssets = bridgeAssets
}

func (p *Module) Start(ctx context.Context) {
	p.Log.Info().Msg("starting parser module...")
	p.G.GoCtx(ctx, p.listen)
}

func (p *Module) Close() error {
	p.Log.Info().Msg("closing...")
	p.G.Wait()
	return nil
}
