package asset

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

const (
	KeyAssets        = "assets"
	KeyAuctions      = "auctions"
	KeyLastAuctionID = "lastAuction"
)

func (cc *SmartContract) AddAsset(
	ctx contractapi.TransactionContextInterface, id, owner string,
) error {
	existing, err := ctx.GetStub().GetState(cc.makeAssetKey(id))
	if err != nil {
		return fmt.Errorf("unable to interact with worldstate: %v", err)
	}

	if existing != nil {
		return fmt.Errorf("asset with ID %s already exists", id)
	}

	asset := Asset{
		ID:    id,
		Owner: owner,
	}

	err = cc.setAsset(ctx, &asset)
	if err != nil {
		return err
	}

	// Emit an event when an asset is added
	eventPayload := "Asset added: " + id
	err = ctx.GetStub().SetEvent("AddAsset", []byte(eventPayload))
	if err != nil {
		return fmt.Errorf("error setting event: %v", err)
	}
	return nil
}

func (cc *SmartContract) StartAuction(
	ctx contractapi.TransactionContextInterface, argjson string,
) error {
	var args StartAuctionArgs
	err := json.Unmarshal([]byte(argjson), &args)
	if err != nil {
		return err
	}

	asset, err := cc.GetAsset(ctx, args.AssetID)
	if err != nil {
		return err
	}
	if asset.PendingAuctionID > 0 {
		return fmt.Errorf("pending auction on asset")
	}

	lastID, err := cc.GetLastAuctionID(ctx)
	if err != nil {
		return err
	}
	auction := Auction{
		ID:         lastID + 1,
		AssetID:    args.AssetID,
		EthAddr:    args.EthAddr,
		QuorumAddr: args.QuorumAddr,
		Status:     "open",
	}
	err = cc.setAuction(ctx, &auction)
	if err != nil {
		return err
	}
	err = cc.setLastAuctionID(ctx, auction.ID)
	if err != nil {
		return err
	}

	asset.PendingAuctionID = auction.ID
	err = cc.setAsset(ctx, asset)
	if err != nil {
		return fmt.Errorf("error setting asset: %v", err)
	}

	// Emit an event when an auction is started
	eventPayload := fmt.Sprintf("Auction start: %d", auction.ID)
	err = ctx.GetStub().SetEvent("StartAuction", []byte(eventPayload))
	if err != nil {
		return fmt.Errorf("error setting event: %v", err)
	}

	return nil
}

func (cc *SmartContract) CancelAuction(
	ctx contractapi.TransactionContextInterface, IDStr string,
) error {

	ID, _ := strconv.Atoi(IDStr)
	auction, err := cc.GetAuction(ctx, ID)
	if err != nil {
		return err
	}
	auction.Status = "closed"
	err = cc.setAuction(ctx, auction)
	if err != nil {
		return fmt.Errorf("error setting auction: %v", err)
	}

	asset, err := cc.GetAsset(ctx, auction.AssetID)
	if err != nil {
		return err
	}

	asset.PendingAuctionID = 0
	err = cc.setAsset(ctx, asset)
	if err != nil {
		return fmt.Errorf("error setting asset: %v", err)
	}

	// Emit an event when an auction is started
	eventPayload := fmt.Sprintf("Auction cancel: %d", auction.ID)
	err = ctx.GetStub().SetEvent("CancelAuction", []byte(eventPayload))
	if err != nil {
		return fmt.Errorf("error setting event: %v", err)
	}

	return nil
}

func (cc *SmartContract) CloseAuction(
	ctx contractapi.TransactionContextInterface, IDStr string,
) error {

	ID, _ := strconv.Atoi(IDStr)
	auction, err := cc.GetAuction(ctx, ID)

	if err != nil {
		return err
	}

	auction.Status = "closing"
	err = cc.setAuction(ctx, auction)
	if err != nil {
		return fmt.Errorf("error setting auction: %v", err)
	}

	// Emit an event when an auction is started
	eventPayload := fmt.Sprintf("Auction closing: %d", auction.ID)
	err = ctx.GetStub().SetEvent("CloseAuction", []byte(eventPayload))
	if err != nil {
		return fmt.Errorf("error setting event: %v", err)
	}

	return nil
}

func (cc *SmartContract) FinAuction(
	ctx contractapi.TransactionContextInterface, argjson string, prcdStr string,
) error {
	// only owner or admin can call this

	var args AuctionResult
	err := json.Unmarshal([]byte(argjson), &args)
	if err != nil {
		return err
	}

	prcd, err := strconv.ParseBool(prcdStr)
	if err != nil {
		return err
	}

	auction, err := cc.GetAuction(ctx, args.AuctionID)
	if err != nil {
		return err
	}

	if !cc.verifyAuctionResult(args) {
		return fmt.Errorf("invalid auction result")
	}

	auction.Status = "closed"
	err = cc.setAuction(ctx, auction)
	if err != nil {
		return err
	}

	asset, err := cc.GetAsset(ctx, auction.AssetID)
	if err != nil {
		return err
	}

	eventPayload := fmt.Sprintf("Owner no change for asset: %s", auction.AssetID)
	if prcd {
		eventPayload = fmt.Sprintf("Owner changed for asset: %s", auction.AssetID)
		asset.Owner = auction.HighestBidder
	}

	asset.PendingAuctionID = 0
	err = cc.setAsset(ctx, asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().SetEvent("AuctionClosed", []byte(eventPayload))
	if err != nil {
		return fmt.Errorf("error setting event: %v", err)
	}

	return nil
}

// can add some mech to check if bidder has DID creditional
func (cc *SmartContract) verifyAuctionResult(result AuctionResult) bool {

	tmp := &AuctionResult{
		Platform:    result.Platform,
		AuctionID:   result.AuctionID,
		AuctionAddr: result.AuctionAddr,

		HighestBid:    result.HighestBid,
		HighestBidder: result.HighestBidder,
	}

	return VerifySignature(tmp.Hash(), result.Signatrue, result.HighestBidder)
}

func (cc *SmartContract) GetAsset(
	ctx contractapi.TransactionContextInterface, assetID string,
) (*Asset, error) {
	var asset Asset
	b, err := ctx.GetStub().GetState(cc.makeAssetKey(assetID))
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("asset not found")
	}
	err = json.Unmarshal(b, &asset)
	return &asset, err
}

func (cc *SmartContract) GetAuction(
	ctx contractapi.TransactionContextInterface, auctionID int,
) (*Auction, error) {
	b, err := ctx.GetStub().GetState(cc.makeAuctionKey(auctionID))
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("auction not found")
	}
	var auction Auction
	err = json.Unmarshal(b, &auction)
	return &auction, err
}

func (cc *SmartContract) GetLastAuctionID(
	ctx contractapi.TransactionContextInterface,
) (int, error) {
	b, err := ctx.GetStub().GetState(KeyLastAuctionID)
	if err != nil {
		return 0, err
	}
	var count int
	json.Unmarshal(b, &count)
	return count, nil
}

func (cc *SmartContract) setAsset(
	ctx contractapi.TransactionContextInterface, asset *Asset,
) error {
	b, _ := json.Marshal(asset)
	err := ctx.GetStub().PutState(cc.makeAssetKey(asset.ID), b)
	if err != nil {
		return fmt.Errorf("set asset error: %v", err)
	}
	return nil
}

func (cc *SmartContract) setAuction(
	ctx contractapi.TransactionContextInterface, auction *Auction,
) error {
	b, _ := json.Marshal(auction)
	err := ctx.GetStub().PutState(cc.makeAuctionKey(auction.ID), b)
	if err != nil {
		return fmt.Errorf("set auction error: %v", err)
	}
	return nil
}

func (cc *SmartContract) setLastAuctionID(
	ctx contractapi.TransactionContextInterface, id int,
) error {
	b, _ := json.Marshal(id)
	return ctx.GetStub().PutState(KeyLastAuctionID, b)
}

func (cc *SmartContract) makeAssetKey(assetID string) string {
	return fmt.Sprintf("%s_%s", KeyAssets, assetID)
}

func (cc *SmartContract) makeAuctionKey(auctionID int) string {
	return fmt.Sprintf("%s_%d", KeyAuctions, auctionID)
}
