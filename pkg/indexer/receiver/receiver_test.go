// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	ic "github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/celenium-io/astria-indexer/pkg/node/mock"
	nodeTypes "github.com/celenium-io/astria-indexer/pkg/node/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/suite"
)

const testIndexerName = "test_indexer"

// ModuleTestSuite -
type ModuleTestSuite struct {
	suite.Suite
	api *mock.MockApi
}

func (s *ModuleTestSuite) InitApi(configureApi func()) {
	ctrl := gomock.NewController(s.T())
	s.api = mock.NewMockApi(ctrl)

	if configureApi != nil {
		configureApi()
	}
}

var cfgDefault = ic.Indexer{
	Name:         testIndexerName,
	ThreadsCount: 1,
	StartLevel:   0,
	BlockPeriod:  10,
}

func (s *ModuleTestSuite) createModule() Module {
	return NewModule(cfgDefault, s.api)
}

func (s *ModuleTestSuite) createModuleEmptyState(cfgOptional *ic.Indexer) Module {
	cfg := cfgDefault
	if cfgOptional != nil {
		cfg = *cfgOptional
	}

	receiverModule := NewModule(cfg, s.api)
	return receiverModule
}

func (s *ModuleTestSuite) TestModule_SuccessOnStop() {
	s.InitApi(func() {
		s.api.EXPECT().Status(gomock.Any()).Return(nodeTypes.Status{}, nil).MinTimes(0)
	})

	receiverModule := s.createModule()

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	stopperModule := stopper.NewModule(cancelCtx)
	err := stopperModule.AttachTo(&receiverModule, StopOutput, stopper.InputName)
	s.Require().NoError(err)

	receiverModule.Init(&storage.State{
		Id:         1,
		Name:       testIndexerName,
		LastHeight: 1000,
		LastHash:   hashOf1000Block,
		LastTime:   time.Time{},
		ChainId:    "explorer-test",
	})

	stopperCtx, stopperCtxCancel := context.WithCancel(context.Background())
	defer stopperCtxCancel()

	stopperModule.Start(stopperCtx)
	receiverModule.Start(ctx)

	defer func() {
		s.Require().NoError(receiverModule.Close())
	}()

	receiverModule.MustOutput(StopOutput).Push(struct{}{})

	<-ctx.Done()

	s.Require().ErrorIs(context.Canceled, ctx.Err())
}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
