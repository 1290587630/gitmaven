// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/stevenroose/go-bitcoin-core-rpc"
	"log"
)

func main() {

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host: "172.27.3.24:16332",
		User: "rpcusername",
		Pass: "rpcpassword",
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Shutdown()

	client.WalletPassphrase("btcc.com",90000)

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)

	address ,err := btcutil.DecodeAddress("mm1mayxAKbtwjSs594x66YTrjjzvRiCcVj",&chaincfg.TestNet3Params)
	if err != nil{
		return
	}
	fmt.Println(address)
	addrs:=  []btcutil.Address{address}
	//addrs[0] =

	listUnspentResult,err := client.ListUnspentMinMaxAddresses(0,99999,addrs)
	fmt.Println(listUnspentResult)

	//var el btcjson.ListUnspentResult
	for _,el := range  listUnspentResult{

		fmt.Println("++++++++++++")

		fmt.Println("txid",el.Address)
		fmt.Println("txid",el.TxID)
		fmt.Println("vout:",el.Vout)
		fmt.Println("amount:",el.Amount)
		fmt.Println("address",el.Address)
		fmt.Println("=============")
	}

}
