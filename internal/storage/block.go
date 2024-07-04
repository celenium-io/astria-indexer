// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlock interface {
	storage.Table[*Block]

	Last(ctx context.Context) (Block, error)
	ByHeight(ctx context.Context, height pkgTypes.Level, withStats bool) (Block, error)
	ByHash(ctx context.Context, hash []byte) (Block, error)
	ByProposer(ctx context.Context, proposerId uint64, limit, offset int, order storage.SortOrder) ([]Block, error)
	ListWithStats(ctx context.Context, limit, offset uint64, order storage.SortOrder) ([]*Block, error)
	ByIdWithRelations(ctx context.Context, id uint64) (Block, error)
}

// Block -
type Block struct {
	bun.BaseModel `bun:"table:block" comment:"Table with blocks"`

	Id           uint64         `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height       pkgTypes.Level `bun:"height"                    comment:"The number (height) of this block"`
	Time         time.Time      `bun:"time,pk,notnull"           comment:"The time of block"`
	VersionBlock uint64         `bun:"version_block"             comment:"Block version"`
	VersionApp   uint64         `bun:"version_app"               comment:"App version"`

	Hash               pkgTypes.Hex `bun:"hash"                 comment:"Block hash"`
	ParentHash         pkgTypes.Hex `bun:"parent_hash"          comment:"Hash of parent block"`
	LastCommitHash     pkgTypes.Hex `bun:"last_commit_hash"     comment:"Last commit hash"`
	DataHash           pkgTypes.Hex `bun:"data_hash"            comment:"Data hash"`
	ValidatorsHash     pkgTypes.Hex `bun:"validators_hash"      comment:"Validators hash"`
	NextValidatorsHash pkgTypes.Hex `bun:"next_validators_hash" comment:"Next validators hash"`
	ConsensusHash      pkgTypes.Hex `bun:"consensus_hash"       comment:"Consensus hash"`
	AppHash            pkgTypes.Hex `bun:"app_hash"             comment:"App hash"`
	LastResultsHash    pkgTypes.Hex `bun:"last_results_hash"    comment:"Last results hash"`
	EvidenceHash       pkgTypes.Hex `bun:"evidence_hash"        comment:"Evidence hash"`
	ProposerId         uint64       `bun:"proposer_id,nullzero" comment:"Proposer internal id"`
	ActionTypes        types.Bits   `bun:"action_types"         comment:"Bit mask for action types contained in block"`

	ChainId         string                    `bun:"-"` // internal field for filling state
	ProposerAddress string                    `bun:"-"` // internal field for proposer
	Addresses       map[string]*Address       `bun:"-"` // internal field for saving address
	Rollups         map[string]*Rollup        `bun:"-"` // internal field for saving rollups
	RollupAddress   map[string]*RollupAddress `bun:"-"` // internal field for saving rollup address
	Validators      map[string]*Validator     `bun:"-"` // internal field for updating validators
	BlockSignatures []BlockSignature          `bun:"-"` // internal field for saving block signatures
	Constants       []*Constant               `bun:"-"` // internal field for updating constants
	Bridges         []*Bridge                 `bun:"-"` // internal field for saving bridges

	Txs      []*Tx       `bun:"rel:has-many"`
	Stats    *BlockStats `bun:"rel:has-one,join:height=height"`
	Proposer *Validator  `bun:"rel:belongs-to"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
