package sync

import (
	"aurascan-backend/chain"
	"aurascan-backend/model"
	"aurascan-backend/util"
	"ch-common-package/logger"
	"ch-common-package/ssdb"
)

// 插入创世区块
func InsertGenesis() {
	var err error
	var blockInDb *model.BlockInDb
	var transactions []*model.TransactionInDb
	var programMap map[string]*model.ProgramInDb
	var transitions []*model.TransitionInDb
	var authorityInDb *model.AuthorityInDb

	block := chain.GetGenesisBlockByHeight()
	blockHash := block.BlockHash
	height := block.Header.MetaData.Height

	if blockHash == "" {
		logger.Errorf("InsertGenesis block info invalid (%d)", height)
	}

	timestamp := block.Header.MetaData.Timestamp
	transaction := block.Transactions

	totalFeeInBlock := 0.0
	totalBaseFeeInBlock := 0.0
	totalPriorityFeeInBlock := 0.0

	authorityInDb = &model.AuthorityInDb{
		Type:      block.Authority.Type,
		Signature: block.Authority.Signature,
	}

	if len(transaction) > 0 {
		for _, v := range transaction {
			baseFee := 0.0
			priorityFee := 0.0

			// execute transitions
			if len(v.Transaction.Execution.Transitions) > 0 {
				for _, trans := range v.Transaction.Execution.Transitions {
					var transition = &model.TransitionInDb{
						Id:            trans.Id,
						TransactionId: v.Transaction.Id,
						State:         v.Status,
						Program:       trans.Program,
						Function:      trans.Function,
						Inputs:        trans.Inputs,
						Outputs:       trans.Outputs,
						Tpk:           trans.Tpk,
						Tcm:           trans.Tcm,
						Height:        height,
						Timestamp:     timestamp,
					}
					model.AddTimesByProgramID(trans.Program)
					transitions = append(transitions, transition)
				}
			}

			//fee transition
			if v.Transaction.Fee.Transition.Id != "" {
				var feeTransition = &model.TransitionInDb{
					Id:            v.Transaction.Fee.Transition.Id,
					TransactionId: v.Transaction.Id,
					State:         v.Status,
					Program:       v.Transaction.Fee.Transition.Program,
					Function:      v.Transaction.Fee.Transition.Function,
					Inputs:        v.Transaction.Fee.Transition.Inputs,
					Outputs:       v.Transaction.Fee.Transition.Outputs,
					Tpk:           v.Transaction.Fee.Transition.Tpk,
					Tcm:           v.Transaction.Fee.Transition.Tcm,
					Height:        height,
					Timestamp:     timestamp,
				}
				model.AddTimesByProgramID(v.Transaction.Fee.Transition.Program)

				for idx, fee := range v.Transaction.Fee.Transition.Inputs {
					//若为public，则依次为storage_fee、priority_fee、
					//若为private，则依次为record、storage_fee、priority_fee
					if fee.Type == "public" {
						if idx == 0 {
							baseFee += util.GetFloatInAleoNum(fee.Value, "u64")
						} else if idx == 1 {
							priorityFee += util.GetFloatInAleoNum(fee.Value, "u64")
						}
					} else if fee.Type == "private" {
						if idx == 1 {
							baseFee += util.GetFloatInAleoNum(fee.Value, "u64")
						} else if idx == 2 {
							priorityFee += util.GetFloatInAleoNum(fee.Value, "u64")
						}
					}
				}
				transitions = append(transitions, feeTransition)
			}

			programId := ""
			if v.Transaction.Deployment.Program != "" {
				programId = util.GetProgramId(v.Transaction.Deployment.Program)
			}
			//if len(v.Transaction.Fee.Transition.Outputs) > 0 {
			//	feeStr := GetFeeFromJson(v.Transaction.Fee.Transition.Outputs[0].Value, v.Transaction.Fee.Transition.Id)
			//	fee = util.GetFloatInAleoNum(feeStr, "u64")
			//}

			fee := baseFee + priorityFee
			totalBaseFeeInBlock += baseFee
			totalPriorityFeeInBlock += priorityFee
			totalFeeInBlock += fee

			var trans = &model.TransactionInDb{
				Type:          v.Transaction.Type,
				TransactionId: v.Transaction.Id,
				BlockHash:     blockHash,
				Height:        height,
				OuterStatus:   v.Status,
				OuterIndex:    v.Index,
				OuterType:     v.Type,
				Timestamp:     timestamp,
				Fee:           fee,
				BaseFee:       baseFee,
				PriorityFee:   priorityFee,
				DeploymentId:  programId,
				Finalize:      v.Finalize,
			}
			transactions = append(transactions, trans)

			if programId != "" {
				var program = &model.ProgramInDb{
					ProgramID:      programId,
					Height:         height,
					Owner:          v.Transaction.Owner.Address,
					OwnerSignature: v.Transaction.Owner.Signature,
					TransactionID:  v.Transaction.Id,
					TimesCalled:    0,
					DeployTime:     timestamp,
					UpdateAt:       timestamp,
				}
				programMap[programId] = program

				programSpec := chain.GetProgramById(programId)
				if err := ssdb.Client.Set(programId, programSpec); err != nil {
					logger.Errorf("SaveLatestBlocksByHeightRange set ssdb for programID=%v | %v", programId, err)
				}
			}
		}
	}

	blockInDb = &model.BlockInDb{
		BlockHash:         block.BlockHash,
		PreviousHash:      block.PreviousHash,
		Epoch:             block.Header.MetaData.Height / 360,
		EpochIndex:        block.Header.MetaData.Height - block.Header.MetaData.Height/360*360,
		PreviousStateRoot: block.Header.PreviousStateRoot,
		TransactionsRoot:  block.Header.TransactionsRoot,
		FinalizeRoot:      block.Header.FinalizeRoot,
		RatificationsRoot: block.Header.RatificationsRoot,
		SolutionsRoot:     block.Header.SolutionsRoot,
		SubdagRoot:        block.Header.SubdagRoot,

		Network:               block.Header.MetaData.Network,
		Round:                 block.Header.MetaData.Round,
		Height:                block.Header.MetaData.Height,
		CumulativeWeight:      block.Header.MetaData.CumulativeWeight,
		CumulativeProofTarget: block.Header.MetaData.CumulativeProofTarget,
		CoinbaseTarget:        block.Header.MetaData.CoinbaseTarget,
		ProofTarget:           block.Header.MetaData.ProofTarget,
		LastCoinbaseTarget:    block.Header.MetaData.LastCoinbaseTarget,
		LastCoinbaseTimestamp: block.Header.MetaData.LastCoinbaseTimestamp,
		Timestamp:             block.Header.MetaData.Timestamp,
		AuthorityType:         block.Authority.Type,
		AbortedTransactionIds: block.AbortedTransactionIds,

		SolutionNum:    0,
		TransactionNum: len(block.Transactions),
		TotalFee:       totalFeeInBlock,
		BaseFee:        totalBaseFeeInBlock,
		PriorityFee:    totalPriorityFeeInBlock,
	}

	{
		if err = authorityInDb.Save(); err != nil {
			logger.Errorf("InsertGenesis authorityInDb.Save | %v", err)
			return
		}

		if err = model.SaveProgramList(programMap); err != nil {
			logger.Errorf("InsertGenesis SaveProgramList | %v", err)
			return
		}

		if err = model.SaveTransitionList(transitions); err != nil {
			logger.Errorf("InsertGenesis SaveTransitionList | %v", err)
			return
		}

		if err = model.SaveTransactionList(transactions); err != nil {
			logger.Errorf("InsertGenesis SaveTransactionList | %v", err)
			return
		}

		if err = blockInDb.Save(); err != nil {
			logger.Errorf("InsertGenesis blockInDb.Save | %v", err)
			return
		}
	}
}
