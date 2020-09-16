package goslp

import (
	"bytes"
	"encoding/hex"
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

func TestDecodeAddressMainnetSlp(t *testing.T) {
	slpAddrStr := "qrkjty23a5yl7vcvcnyh4dpnxxzuzs4lzqvesp65yq"
	addr, err := DecodeAddress(slpAddrStr, &chaincfg.MainNetParams, false)
	if err != nil {
		t.Fatal(err)
	}
	if addr.String() != slpAddrStr {
		t.Fatal("decode failed")
	}
}

func TestDecodeAddressMainnetCash(t *testing.T) {
	slpAddrStr := "qrkjty23a5yl7vcvcnyh4dpnxxzuzs4lzqvesp65yq"
	cashAddrStr := "qrkjty23a5yl7vcvcnyh4dpnxxzuzs4lzqqzm60567"
	addr, err := DecodeAddress(slpAddrStr, &chaincfg.MainNetParams, true)
	if err != nil {
		t.Fatal(err)
	}
	if addr.String() != cashAddrStr {
		t.Fatal("decode failed")
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
