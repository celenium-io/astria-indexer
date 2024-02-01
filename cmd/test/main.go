// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/base64"
	"log"

	astria "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/sequencer/v1alpha1"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"google.golang.org/protobuf/proto"
)

func main() {
	raw := "CkAJBnO1LvHIhHCN/20WlsFtDD4Rj6E6LqQ4NgcFwY1TmeYdcWXvTss3xZOhC9HgqDeZd+a1mk37v8Aj8XxyU6IBEiDsYdXD5V4Xza1Y2UcI61HZVIjxmIXFBSEIu9lkT8BLNxq8ARKXARKUAQogGbqKuz5LVqMJ32dWxHuX4pjjpy2IRJ02oPrbHKc2ZTkScPhugIQ7msoHglIIlGjdIoO7Dxy1mWnvIFva7Sm+BRCkiIrHIwSJ6AAAgIMb2YKg5WhSydq2/TlLUEbDP5noh0qFtp0pCBiYlun4VLYvr9egJ4tsIEA314WJIGcpTMC4Ql7WsCIZ6APHlUp57IP+4xkaIHBAMcho/T08hKHPqMtF3rpOp0a0Rpf39KbtG49sI5uC"

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		panic(err)
	}

	var tx astria.SignedTransaction
	if err := proto.Unmarshal(decoded, &tx); err != nil {
		panic(err)
	}

	address := ed25519.PubKey(tx.PublicKey).Address()

	log.Println("Address: ", address.String())
	log.Println("Public Key: ", tx.PublicKey)
	log.Println("Signature:", tx.Signature)
	log.Println("Nonce: ", tx.Transaction.Nonce)
	log.Println("Fee asset id: ", tx.Transaction.FeeAssetId)
	log.Println("Actions: ", len(tx.Transaction.Actions))

	for _, action := range tx.Transaction.Actions {
		if val := action.GetSequenceAction(); val != nil {
			log.Println("Sequence action rollup id: ", val.RollupId)
			log.Println("Sequence action data: ", val.Data)
		}
		if val := action.GetTransferAction(); val != nil {
			log.Println("Transfer", val.Amount)
		}
	}
}
