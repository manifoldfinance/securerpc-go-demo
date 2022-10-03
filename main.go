package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
  // stats not supported post merge
	j               = "application/json"
	flashbotURL     = "https://api.securerpc.com/v1"
	stats           = "flashbots_getUserStats"
	sb		          = "manifold_sendBundle"
	flashbotXHeader = "X-Flashbots-Signature"
	p               = "POST"
)

// pubAddr: 0x45Bdc5051bBf3319DA289BA222E1c5F99136Ed43
// https://github.com/flashbots/rpc-endpoint/commit/d42ae0b665833391398628358ca03dae1964868d#diff-311b8f116ec381569dfecd92d10675b9ff9acffcbc86d66ca4bc0aca79b6e991L35
// Flashbots RPC Relay Test Key
var (
	privateKey, _ = crypto.HexToECDSA(
    "7bdeed70a07d5a45546e83a88dd430f71348592e747d2d3eb23f32db003eb0e1",
//		"2e19800fcbbf0abb7cf6d72ee7171f08943bc8e5c3568d1d7420e52136898154",
	)
)

func flashbotHeader(signature []byte, privateKey *ecdsa.PrivateKey) string {
	return crypto.PubkeyToAddress(privateKey.PublicKey).Hex() +
		":" + hexutil.Encode(signature)
}

func main() {
	mevHTTPClient := &http.Client{
		Timeout: time.Second * 3,
	}
	// Fri Sep 30 16:01:18 PDT 2022
	// 15_649_299 
	currentBlock := big.NewInt(15_649_299)
	params := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  sb,
		"params": []interface{}{
			fmt.Sprintf("0x%x", currentBlock.Uint64()),
		},
	}
	payload, _ := json.Marshal(params)
	req, _ := http.NewRequest(p, flashbotURL, bytes.NewBuffer(payload))
	headerReady, _ := crypto.Sign(
		accounts.TextHash([]byte(hexutil.Encode(crypto.Keccak256(payload)))),
		privateKey,
	)
	req.Header.Add("content-type", j)
	req.Header.Add("Accept", j)
	req.Header.Add(flashbotXHeader, flashbotHeader(headerReady, privateKey))
	resp, _ := mevHTTPClient.Do(req)
	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
}
