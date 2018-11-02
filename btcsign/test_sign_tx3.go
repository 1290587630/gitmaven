package main

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"log"
)

const (
	TXID        = "a0dd46a7da775ba9ccfba6247b1343076e07105adcbae733a9e167f26e931a75" // UTXO TXID
	AMOUNT      = 6942902                                                            // UTXO 余额
	FROMADDRESS = "mm1mayxAKbtwjSs594x66YTrjjzvRiCcVj"                               // 转出地址
	TOADDRESS1  = "2N4U7fBDD8MKtagdXgpMnmqW6gf2CcnKttw"                              // 目标地址
	TOADDRESS2  = "mm1mayxAKbtwjSs594x66YTrjjzvRiCcVj"                               // 找零地址
	FROMWIF     = "cTrTtAVHm6eiEXtHgqze9s4NjnmGUJSwEpph3PW5krkHoapamr7U"             // 转出地址的私钥
	TXFEE       = 10000                                                              // 手续费
	MINFEE      = 546                                                                // 最小消费
	nindex		= 2
)

func main() {
	params := &chaincfg.TestNet3Params
	params.Name = "testnet3"
	params.Net = wire.BitcoinNet(0x0B110907) //mainnet 0xf9beb4d9
	params.PubKeyHashAddrID = byte(0x6F)
	params.PrivateKeyID = byte(0x80)

	tx := wire.NewMsgTx(wire.TxVersion)

	// input
	fromAddress, err := btcutil.DecodeAddress(FROMADDRESS, params)
	if err != nil {
		log.Fatalf("invalid address: %v", err)
	}
	hash, err := chainhash.NewHashFromStr(TXID)
	if err != nil {
		log.Fatalf("could not get hash from transaction ID: %v", err)
	}
	outPoint := wire.NewOutPoint(hash, nindex)
	txIn := wire.NewTxIn(outPoint, nil, nil)
	tx.AddTxIn(txIn)

	pkscript, err := txscript.PayToAddrScript(fromAddress)
	if err != nil {
		log.Fatalf("could not get pay from address script: %v", err)
	}

	// output
	// Pay the minimum network fee so that nodes will broadcast the tx.
	toAddress1, err := btcutil.DecodeAddress(TOADDRESS1, params)
	if err != nil {
		log.Fatalf("invalid address: %v", err)
	}
	script1, err := txscript.PayToAddrScript(toAddress1)
	if err != nil {
		log.Fatalf("could not get pay to address script: %v", err)
	}
	outCoin1 := int64(MINFEE)
	txOut1 := wire.NewTxOut(outCoin1, script1)
	tx.AddTxOut(txOut1)

	// OP_RETURN
	// Put a message here with your name or MIT ID number so I can find your
	// submission on the blockchain.
	//opReturnData := []byte{'o', 'm', 'n', 'i',
	//	0x00, 0x00,
	//	0x00, 0x00,
	//	0x00, 0x00, 0x00, 0x01,
	//	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a}
	//// build the op_return output script
	//// this is the OP_RETURN opcode, followed by a data push opcode, then the data.
	//opReturnScript, err :=
	//	txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData(opReturnData).Script()
	//if err != nil {
	//	panic(err)
	//}
	//opreturn := wire.NewTxOut(0, opReturnScript)
	//tx.AddTxOut(opreturn)

	toAddress2, err := btcutil.DecodeAddress(TOADDRESS2, params)
	if err != nil {
		log.Fatalf("invalid address: %v", err)
	}
	script2, err := txscript.PayToAddrScript(toAddress2)
	if err != nil {
		log.Fatalf("could not get pay to address script: %v", err)
	}
	outCoin2 := int64(AMOUNT-TXFEE) - outCoin1
	txOut2 := wire.NewTxOut(outCoin2, script2)
	tx.AddTxOut(txOut2)

	// private key
	wif, err := btcutil.DecodeWIF(FROMWIF)
	if err != nil {
		log.Fatalf("could not decode wif: %v", err)
	}

	sig, err := txscript.SignatureScript(
		tx,                  // The tx to be signed.
		0,                   // The index of the txin the signature is for.
		pkscript,            // The other half of the script from the PubKeyHash.
		txscript.SigHashAll, // The signature flags that indicate what the sig covers.
		wif.PrivKey,         // The key to generate the signature with.
		true,                // The compress sig flag. This saves space on the blockchain.
	)

	if err != nil {
		log.Fatalf("could not generate signature: %v", err)
	}

	tx.TxIn[0].SignatureScript = sig

	log.Printf("signed raw transaction: %s", txToHex(tx))

	// Validate signature
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(pkscript, tx, 0, flags, nil, nil, outCoin1)
	if err != nil {
		log.Printf("err != nil: %v\n", err)
	}
	if err := vm.Execute(); err != nil {
		log.Printf("vm.Execute > err != nil: %v\n", err)
	}
	log.Printf("verify ok\n")
}

func txToHex(tx *wire.MsgTx) string {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)
	return hex.EncodeToString(buf.Bytes())
}
