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
	Medical_Asset_type_ID        = "field_type"
	KeyAuctions      = "lb_id"
	KeyLastAuctionID = "user_did"
	ActuionID = "data_hash"
)

func (cc *SmartContract) AddAsset(
	ctx contractapi.TransactionContextInterface, user_did, owner string,
) error {
	existing, err := ctx.GetStub().GetState(cc.makeAssetKey(user_did))
	if err != nil {
		return fmt.Errorf("unable to interact with worldstate for the medical asset: %v", err)
	}

	if existing != nil {
		return fmt.Errorf("Medical_Asset_type_ID with user_did %s already exists", user_did)
	}

	asset := Asset{
		did:    user_did,
		Owner: schema_id,
	}

	err = cc.setAsset(ctx, &field_type)
	if err != nil {
		return err
	}

	// Emit an event when an Medical_Asset is added
	eventPayload := "Medical_Asset added: " + did
	err = ctx.GetStub().SetEvent("AddAsset", []byte(eventPayload))
	if err != nil {
		return fmt.Errorf("error setting event for medical asset: %v", err)
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

	asset, err := cc.GetAsset(ctx, args.field_type)
	if err != nil {
		return err
	}
	if asset.PendingAuctionID > 0 {
		return fmt.Errorf("pending auction on medical asset type")
	}

	lastID, err := cc.GetLastAuctionID(ctx)
	if err != nil {
		return err
	}
	auction := Auction{
		did:         lastID + 1,
		Medical_Asset_type_ID:    args.field_type,
		EthAddr:    args.EthAddr,
		QuorumAddr: args.QuorumAddr,
		Status:     "status",
	}
	err = cc.setAuction(ctx, &auction)
	if err != nil {
		return err
	}
	err = cc.setLastAuctionID(ctx, auction.data_hash)
	if err != nil {
		return err
	}

	asset.PendingAuctionID = auction.data_hash
	err = cc.setAsset(ctx, asset)
	if err != nil {
		return fmt.Errorf("error setting medical asset: %v", err)
	}

	// Emit an event when an auction for the medical asset is started
	eventPayload := fmt.Sprintf("Auction start for medical asset: %d", auction.data_hash)
	err = ctx.GetStub().SetEvent("StartAuction", []byte(eventPayload))
	if err != nil {
		return fmt.Errorf("error setting event: %v", err)
	}

	return nil
}

func (cc *SmartContract) CancelAuction(
	ctx contractapi.TransactionContextInterface, IDStr string,
	IDStr = "did"
) error {

	Medical_Asset_type_ID, _ := strconv.Atoi(IDStr)
	auction, err := cc.GetAuction(ctx, data_hash)
	if err != nil {
		return err
	}
	auction.Status = status.Rejected
	err = cc.setAuction(ctx, auction)
	if err != nil {
		return fmt.Errorf("error setting auction for the medical asset: %v", err)
	}

	asset, err := cc.GetAsset(ctx, auction.Medical_Asset_type_ID)
	if err != nil {
		return err
	}

	asset.PendingAuctionID = 0
	err = cc.setAsset(ctx, asset)
	if err != nil {
		return fmt.Errorf("error setting asset: %v", err)
	}

	// Emit an event when an auction is started
	eventPayload := fmt.Sprintf("Auction cancel: %d", auction.data hash)
	err = ctx.GetStub().SetEvent("CancelAuction", []byte(eventPayload))
	if err != nil {
		return fmt.Errorf("error setting event: %v", err)
	}

	return nil
}

func (cc *SmartContract) CloseAuction(
	ctx contractapi.TransactionContextInterface, IDStr string,
) error {

	Medical_Asset_type_ID, _ := strconv.Atoi(IDStr)
	auction, err := cc.GetAuction(ctx, data_hash)

	if err != nil {
		return err
	}

	auction.Status = status.Validated
	err = cc.setAuction(ctx, auction)
	if err != nil {
		return fmt.Errorf("error setting auction for medical asset: %v", err)
	}

	// Emit an event when an auction  for medical asset is started
	eventPayload := fmt.Sprintf("Auction closing for medical asset: %d", auction.data_hash)
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

	auction, err := cc.GetAuction(ctx, args.data_hash)
	if err != nil {
		return err
	}

	if !cc.verifyAuctionResult(args) {
		return fmt.Errorf("invalid medical asset auction result")
	}

	auction.Status = status.Pending
	err = cc.setAuction(ctx, auction)
	if err != nil {
		return err
	}

	asset, err := cc.GetAsset(ctx, auction.field_type))
	if err != nil {
		return err
	}

	eventPayload := fmt.Sprintf("Owner no change for medical asset: %s", auction.field_type)
	if prcd {
		eventPayload = fmt.Sprintf("Owner changed for medical asset: %s", auction.field_type)
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
		Platform:    result.lb_id,
		AuctionID:   result.field_type,
		AuctionAddr: result.data_hash,

		ValidBid:    result.status,
		ValidBidder: result.user_did,
	}

	return VerifySignature(tmp.Hash(), result.data_sig, result.ValidtBidder)
}

func (cc *SmartContract) GetAsset(
	ctx contractapi.TransactionContextInterface, Medical_Asset_type_ID,
) (*Asset, error) {
	var asset Asset
	b, err := ctx.GetStub().GetState(cc.makeAssetKey(Medical_Asset_type_ID))
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("Medical asset not found")
	}
	err = json.Unmarshal(b, &asset)
	return &asset, err
}

func (cc *SmartContract) GetAuction(
	ctx contractapi.TransactionContextInterface, Medical_Asset_type_ID,
) (*Auction, error) {
	b, err := ctx.GetStub().GetState(cc.makeAuctionKey(Medical_Asset_type_ID))
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("Medical asset not found")
	}
	var auction Auction
	err = json.Unmarshal(b, &auction)
	return &auction, err
}

func (cc *SmartContract) GetLastAuctionID(
	ctx contractapi.TransactionContextInterface,
) (int, error) {
	b, err := ctx.GetStub().GetState(Medical_Asset_type_ID)
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
	err := ctx.GetStub().PutState(cc.makeAssetKey(asset.status, b)
	if err != nil {
		return fmt.Errorf("set Medical asset error: %v", err)
	}
	return nil
}

func (cc *SmartContract) setAuction(
	ctx contractapi.TransactionContextInterface, auction *Auction,
) error {
	b, _ := json.Marshal(auction)
	err := ctx.GetStub().PutState(cc.makeAuctionKey(auction.data_hash), b)
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

func (cc *SmartContract) makeAssetKey(Medical_Asset_type_ID string) string {
	return fmt.Sprintf("%s_%s", KeyAssets, data_hash)
}

func (cc *SmartContract) makeAuctionKey(data_hash int) string {
	return fmt.Sprintf("%s_%d", KeyAuctions, data_hash)
}
