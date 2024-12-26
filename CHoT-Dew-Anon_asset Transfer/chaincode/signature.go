package asset

import (
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(hash, signature []byte, addr string) bool {
	pubkey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return false
	}

	if addr == crypto.PubkeyToAddress(*pubkey).Hex() {
		return true
	}
	return false
}

// func VerifyKnownSignatures(hash []byte, sigs [][]byte, addrs []common.Address) bool {
// 	set := make(map[string]bool, len(addrs))
// 	for _, addr := range addrs {
// 		set[addr.Hex()] = true
// 	}

// 	for _, sig := range sigs {
// 		addr, err := VerifySignature(hash, sig)
// 		if err != nil {
// 			return false
// 		}
// 		if !set[addr.Hex()] {
// 			return false
// 		}
// 	}
// 	return true
// }
