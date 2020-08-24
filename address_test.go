package goslp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil"
)

func TestDecodeSlpAddress(t *testing.T) {
	addrStr := "simpleledger:qrkjty23a5yl7vcvcnyh4dpnxxzuzs4lzqvesp65yq"
	prefix, data, _ := DecodeCashAddress(addrStr)
	if prefix != "simpleledger" {
		t.Fatal("decode failed")
	}
	if len(data) != 34 {
		t.Fatal("data wrong length")
	}
}

func TestDecodeAddressMainnet(t *testing.T) {
	addrStr := "qrkjty23a5yl7vcvcnyh4dpnxxzuzs4lzqvesp65yq"
	addr, err := DecodeAddress(addrStr, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addr.String())
	if addr.String() != addrStr {
		t.Fatal("decode failed")
	}
}

func TestCashAddrFailedDecodeAddressMainnet(t *testing.T) {
	addrStr := "qqqmy7340gd5esk26zvgxmh8fpkq36e32vv6gd69dv"
	_, err := DecodeAddress(addrStr, &chaincfg.MainNetParams)
	if err == nil {
		t.Fatal(err)
	}
}

func TestCashAddr2FailedDecodeAddressMainnet(t *testing.T) {
	addrStr := "bitcoincash:qqqmy7340gd5esk26zvgxmh8fpkq36e32vv6gd69dv"
	_, err := DecodeAddress(addrStr, &chaincfg.MainNetParams)
	if err == nil {
		t.Fatal(err)
	}
}

func TestConvertP2pkhBchToSlpAddress(t *testing.T) {
	addrStr := "qprqzzhhve7sgysgf8h29tumywnaeyqm7y6e869uc6"
	params := &chaincfg.Params{
		SlpAddressPrefix:  "simpleledger",
		CashAddressPrefix: "bitcoincash",
	}
	addr, _ := bchutil.DecodeAddress(addrStr, params)
	hash := addr.(*bchutil.AddressPubKeyHash).Hash160()

	slpAddr, _ := NewAddressPubKeyHash(hash[:], params)
	if slpAddr.String() != "qprqzzhhve7sgysgf8h29tumywnaeyqm7ykzvpsuxy" {
		t.Fatal("incorrect conversion from cashAddr to bchAddr")
	}
}

func TestConvertP2shBchToSlpAddress(t *testing.T) {
	addrStr := "pzmj0ueqasnsw80a26th5t2gsz5evcxsps2tavljvp"
	params := &chaincfg.Params{
		SlpAddressPrefix:  "simpleledger",
		CashAddressPrefix: "bitcoincash",
	}
	addr, _ := bchutil.DecodeAddress(addrStr, params)
	hash := addr.(*bchutil.AddressScriptHash).Hash160()

	expHash, _ := hex.DecodeString("0e13f9ce9a6be9653d812a030e936767b503debe")
	if bytes.Equal(expHash, hash[:]) {
		t.Fatal("incorrect hash decoded from bchutil")
	}

	slpAddr, _ := NewAddressScriptHashFromHash(hash[:], params)
	if slpAddr.String() != "pzmj0ueqasnsw80a26th5t2gsz5evcxspsxskh2jjl" {
		t.Fatal("incorrect conversion from cashAddr to bchAddr")
	}
}
