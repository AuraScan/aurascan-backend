package chain

import (
	"aurascan-backend/model"
	"aurascan-backend/util"
	"ch-common-package/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func GetLatestHeight() int64 {
	var height int64

	res, err := http.Get(util.AleoNodeApi + "/latest/height")
	if err != nil {
		fmt.Printf("GetLatestHeight Get failed | %v", err)
		return 0
	}

	jsonData, err := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &height)
	return height
}

func GetLatestHash() {
	var blockHash string

	res, err := http.Get(util.AleoNodeApi + "/latest/hash")
	if err != nil {
		fmt.Printf("GetLatestHeight Get failed | %v", err)
		return
	}

	jsonData, err := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &blockHash)
	fmt.Printf("GetLatestHash hash:%v\n", blockHash)
}

func GetLatestBlock() (model.Block, error) {
	var block model.Block

	res, err := http.Get(util.AleoNodeApi + "/latest/block")
	if err != nil {
		return model.Block{}, fmt.Errorf("GetLatestBlock Get failed | %v", err)
	}

	jsonData, err := io.ReadAll(res.Body)
	if err = json.Unmarshal(jsonData, &block); err != nil {
		return model.Block{}, fmt.Errorf("GetLatestBlock Unmarshal failed | %v", err)
	}
	return block, nil
}

func GetBlockByHeight(height int64) model.Block {
	var block model.Block

	res, err := http.Get(util.AleoNodeApi + "/block/" + strconv.FormatInt(height, 10))
	if err != nil {
		logger.Errorf("GetBlockByHeight Get failed | %v", err)
		return model.Block{}
	}

	jsonData, err := io.ReadAll(res.Body)

	if err = json.Unmarshal(jsonData, &block); err != nil {
		logger.Errorf("GetBlockByHeight Unmarshal failed | %v", err)
		return model.Block{}
	}
	//fmt.Printf("GetBlockByHeight block:%v\n", block)
	return block
}

func GetGenesisBlockByHeight() model.BlockGenesis {
	var block model.BlockGenesis

	res, err := http.Get(util.AleoNodeApi + "/block/" + strconv.FormatInt(0, 10))
	if err != nil {
		logger.Errorf("GetBlockByHeight Get failed | %v", err)
		return model.BlockGenesis{}
	}

	jsonData, err := io.ReadAll(res.Body)

	if err = json.Unmarshal(jsonData, &block); err != nil {
		logger.Errorf("GetBlockByHeight Unmarshal failed | %v", err)
		return model.BlockGenesis{}
	}
	return block
}

func GetProgramById(programId string) string {
	res, err := http.Get(util.AleoNodeApi + "/program/" + programId)
	if err != nil {
		fmt.Printf("GetProgramById Get failed | %v", err)
		return ""
	}
	var program string
	jsonData, err := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &program)

	return program
}

func GetProgramNameById(programId string) []string {
	res, err := http.Get(util.AleoNodeApi + "/program/" + programId + "/mappings")
	if err != nil {
		fmt.Printf("GetProgramNameById Get failed | %v", err)
		return []string{}
	}
	var names = make([]string, 0)
	jsonData, err := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &names)

	return names
}

// 通过key获取值
func GetProgramMapValue(programId, mappingName, MappingKey string) string {
	res, err := http.Get(util.AleoNodeApi + "/program/" + programId + "/mapping/" + mappingName + "/" + MappingKey)
	if err != nil {
		fmt.Printf("GetProgramNameById Get failed | %v", err)
		return ""
	}
	var str string
	jsonData, err := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &str)

	return str
}

func GetTransactionById(id string) model.Transaction {
	var transaction model.Transaction

	res, err := http.Get(util.AleoNodeApi + "/transaction/" + id)
	if err != nil {
		fmt.Printf("GetTransactionById Get failed | %v", err)
		return model.Transaction{}
	}

	jsonData, err := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &transaction)
	fmt.Printf("GetTransactionById transaction:%v\n", transaction)
	return transaction
}

func GetTransactionsByHeight(height int64) []*model.Transaction {
	var transactions []*model.Transaction

	res, err := http.Get(util.AleoNodeApi + "/transactions/" + strconv.FormatInt(height, 10))
	if err != nil {
		fmt.Printf("GetTransactionsByHeight Get failed | %v", err)
		return nil
	}

	jsonData, err := io.ReadAll(res.Body)
	json.Unmarshal(jsonData, &transactions)
	fmt.Printf("GetTransactionsByHeight block:%v\n", transactions)
	return transactions
}
