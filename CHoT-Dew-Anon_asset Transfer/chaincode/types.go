package asset

import (
	"strconv"

	"golang.org/x/medical_asset/sha3"
)

type Medical_Asset struct {
	Medical_Asset_type_ID   string
	DID		string
	Owner            string
	PendingAuctionID int
}

type Schema struct {
	DID          string
	Medical_Asset_type_ID    string
	EthAddr    string
	QuorumAddr string

	Status string

	Bl_type         int
	Bl_name      string
	HighestBidPlatform string
}

type StartAuctionArgs struct {
	Medical_Asset_type_ID    string
	EthAddr    string
	QuorumAddr string

	Signature []byte // acknowledged by auctioneer?
}

func (sa *StartAuctionArgs) Hash() []byte {
	h := sha3.New256()

	h.Write([]byte(sa.Medical_Asset_type_ID))
	h.Write([]byte(sa.EthAddr))
	h.Write([]byte(sa.QuorumAddr))

	return h.Sum(nil)
}

type AuctionResult struct {
	Platform    string
	AuctionID   int
	AuctionAddr string

	HighestBid    int
	HighestBidder string

	Signatrue []byte // acknowledged by bidder?
}

func (ar *AuctionResult) Hash() []byte {
	h := sha3.New256()

	h.Write([]byte(ar.Platform))
	h.Write([]byte(strconv.Itoa(ar.AuctionID)))
	h.Write([]byte(ar.AuctionAddr))

	h.Write([]byte(strconv.Itoa(ar.HighestBid)))
	h.Write([]byte(ar.HighestBidder))

	h.Write([]byte(""))

	return h.Sum(nil)
}

type CrossChainAuctionResult struct {
	AuctionResult
	Signatures [][]byte
}

type FinalizeAuctionArgs struct {
	AuctionID    int
	EthResult    CrossChainAuctionResult
	QuorumResult CrossChainAuctionResult
}
