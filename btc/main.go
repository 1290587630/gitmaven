package main

import (
	"fmt"
	"github.com/Swipecoin/go-bip44"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func GenerateBTC() (string, string, error) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", "", err
	}

	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
	if err != nil {
		return "", "", err
	}
	pubKeySerial := privKey.PubKey().SerializeUncompressed()

	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", err
	}

	return privKeyWif.String(), pubKeyAddress.EncodeAddress(), nil
}

func GenerateBTCTest() (string, string, error) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", "", err
	}

	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, false)
	if err != nil {
		return "", "", err
	}
	pubKeySerial := privKey.PubKey().SerializeUncompressed()

	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.TestNet3Params)
	if err != nil {
		return "", "", err
	}

	return privKeyWif.String(), pubKeyAddress.EncodeAddress(), nil
}

func test(){
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "Secret Passphrase")

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	// Display mnemonic and keys
	fmt.Println("Mnemonic: ", mnemonic)
	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)
}

func main(){

	testpriv()
}

func testpriv(){
	bitSize := 256
	mnemonic, _ := bip44.NewMnemonic(bitSize)
	fmt.Println(mnemonic)

	seedBytes ,err:= mnemonic.NewSeed("my password")
	if(err !=nil){
		return
	}
	testaddr(seedBytes)

	xKey, err := hdkeychain.NewMaster(seedBytes, &chaincfg.TestNet3Params)
	fmt.Println("xkey:",xKey)
	privKey, err := xKey.ECPrivKey()
	fmt.Println("privKey:",privKey.ToECDSA())

	addressPubKeyHash, err :=xKey.Address(&chaincfg.TestNet3Params)
	fmt.Println("addressPubKeyHash:",addressPubKeyHash)


	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, false)
	fmt.Println("privKeyWif:",privKeyWif)

	pubKeySerial := privKey.PubKey().SerializeUncompressed()

	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.TestNet3Params)
	fmt.Println("pubKeyAddress:",pubKeyAddress)

	pubKeyAddress.EncodeAddress()
	fmt.Println("pubKeyAddress.EncodeAddress():",pubKeyAddress.EncodeAddress())
	return
}


func testaddr(seedBytes []byte){
	xKey, _ := bip44.NewKeyFromSeedBytes(seedBytes, bip44.TESTNET3)
	fmt.Println(xKey)

	accountKey, _ := xKey.BIP44AccountKey(bip44.BitcoinCoinType, 0, true)
	fmt.Println(accountKey)

	externalAddress, _ := accountKey.DeriveP2PKAddress(bip44.ExternalChangeType, 0, bip44.TESTNET3)

	fmt.Println(accountKey)
	fmt.Println(externalAddress)
}
