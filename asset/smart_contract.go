package asset

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
)

/* This is a sample chaincode implemented as per the Fabric documentation */

var (
	assets = []Asset{
		{ID: 1, Color: "blue", Owner: "John Doe", Value: 500},
		{ID: 2, Color: "red", Owner: "Jane Doe", Value: 600},
		{ID: 3, Color: "yellow", Owner: "Bill", Value: 450},
	}
)

type SmartContract struct {
	contractapi.Contract
}

// Asset attributes are defined in alphabetical order to make JSON struct deterministic
type Asset struct {
	Color string `json:"color"`
	ID    int    `json:"id"`
	Owner string `json:"owner"`
	Value int    `json:"value"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	for _, a := range assets {
		aByts, err := json.Marshal(a)
		if err != nil {
			return fmt.Errorf(`marshal asset failed for asset %d - %w`, a.ID, err)
		}

		if err = ctx.GetStub().PutState(strconv.Itoa(a.ID), aByts); err != nil {
			return fmt.Errorf(`put asset failed for asset %d - %w`, a.ID, err)
		}
	}

	return nil
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, color string, id int, owner string, val int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`create asset failed - %w`, err)
	}

	if exists {
		return fmt.Errorf(`asset with id %d already exists`, id)
	}

	asset := Asset{
		Color: color,
		ID:    id,
		Owner: owner,
		Value: val,
	}

	aByts, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(asset.ID), aByts)
}

func (s *SmartContract) GetAsset(ctx contractapi.TransactionContextInterface, id int) (*Asset, error) {
	aByts, err := ctx.GetStub().GetState(strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf(`get state failed for asset %d - %w`, id, err)
	}

	if aByts == nil {
		return nil, fmt.Errorf(`asset does not exist for id %d`, id)
	}

	var a Asset
	if err = json.Unmarshal(aByts, &a); err != nil {
		return nil, fmt.Errorf(`unmarshal asset failed for asset %d - %w`, id, err)
	}

	return &a, nil
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, color string, id int, owner string, val int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`checking asset existence failed - %w`, err)
	}

	if !exists {
		return fmt.Errorf(`asset with id %d does not exist`, id)
	}

	a := Asset{
		Color: color,
		ID:    id,
		Owner: owner,
		Value: val,
	}

	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(a.ID), aByts)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`checking asset existence failed - %w`, err)
	}

	if !exists {
		return fmt.Errorf(`asset with id %d does not exist`, id)
	}

	return ctx.GetStub().DelState(strconv.Itoa(id))
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id int, newOwner string) error {
	a, err := s.GetAsset(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Owner = newOwner
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// ranging with empty start and end keys returns all assets in chaincode namespace
	itr, err := ctx.GetStub().GetStateByRange(``, ``)
	if err != nil {
		return nil, fmt.Errorf(`get state by range faied - %w`, err)
	}
	defer itr.Close()

	var ats []*Asset
	for itr.HasNext() {
		res, err := itr.Next()
		if err != nil {
			return nil, fmt.Errorf(`iterating next query result failed - %w`, err)
		}

		var a Asset
		if err = json.Unmarshal(res.Value, &a); err != nil {
			return nil, fmt.Errorf(`unmarshal failed - %w`, err)
		}

		ats = append(ats, &a)
	}

	return ats, nil
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id int) (bool, error) {
	aByts, err := ctx.GetStub().GetState(strconv.Itoa(id))
	if err != nil {
		return false, fmt.Errorf(`get state failed for asset %d - %w`, id, err)
	}

	return aByts != nil, nil
}

func (s *SmartContract) ChangeAssetColour(ctx contractapi.TransactionContextInterface, id int, clr string) error {
	a, err := s.GetAsset(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Color = clr
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

func (s *SmartContract) ChangeAssetValue(ctx contractapi.TransactionContextInterface, id int, val int) error {
	a, err := s.GetAsset(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Value = val
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

// Vehicle functions

func (s *SmartContract) CreateVehicle(ctx contractapi.TransactionContextInterface, color string, id int, owner string, val int) error {
	exists, err := s.VehicleExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`create vehicle failed - %w`, err)
	}

	if exists {
		return fmt.Errorf(`vehicle with id %d already exists`, id)
	}

	vehicle := Asset{
		Color: color,
		ID:    id,
		Owner: owner,
		Value: val,
	}

	aByts, err := json.Marshal(vehicle)
	if err != nil {
		return fmt.Errorf(`marshal vehicle failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(vehicle.ID), aByts)
}

func (s *SmartContract) GetVehicle(ctx contractapi.TransactionContextInterface, id int) (*Asset, error) {
	aByts, err := ctx.GetStub().GetState(strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf(`get state failed for vehicle %d - %w`, id, err)
	}

	if aByts == nil {
		return nil, fmt.Errorf(`vehicle does not exist for id %d`, id)
	}

	var a Asset
	if err = json.Unmarshal(aByts, &a); err != nil {
		return nil, fmt.Errorf(`unmarshal vehicle failed for vehicle %d - %w`, id, err)
	}

	return &a, nil
}

func (s *SmartContract) UpdateVehicle(ctx contractapi.TransactionContextInterface, color string, id int, owner string, val int) error {
	exists, err := s.VehicleExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`checking vehicle existence failed - %w`, err)
	}

	if !exists {
		return fmt.Errorf(`vehicle with id %d does not exist`, id)
	}

	a := Asset{
		Color: color,
		ID:    id,
		Owner: owner,
		Value: val,
	}

	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal vehicle failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(a.ID), aByts)
}

func (s *SmartContract) DeleteVehicle(ctx contractapi.TransactionContextInterface, id int) error {
	exists, err := s.VehicleExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`checking vehicle existence failed - %w`, err)
	}

	if !exists {
		return fmt.Errorf(`vehicle with id %d does not exist`, id)
	}

	return ctx.GetStub().DelState(strconv.Itoa(id))
}

func (s *SmartContract) TransferVehicle(ctx contractapi.TransactionContextInterface, id int, newOwner string) error {
	a, err := s.GetVehicle(ctx, id)
	if err != nil {
		return fmt.Errorf(`get vehicle failed - %w`, err)
	}

	a.Owner = newOwner
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal vehicle failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

func (s *SmartContract) GetAllVehicles(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// ranging with empty start and end keys returns all vehicles in chaincode namespace
	// an identical function of get all assets
	itr, err := ctx.GetStub().GetStateByRange(``, ``)
	if err != nil {
		return nil, fmt.Errorf(`get vehicle state by range faied - %w`, err)
	}
	defer itr.Close()

	var ats []*Asset
	for itr.HasNext() {
		res, err := itr.Next()
		if err != nil {
			return nil, fmt.Errorf(`iterating next vehicle query result failed - %w`, err)
		}

		var a Asset
		if err = json.Unmarshal(res.Value, &a); err != nil {
			return nil, fmt.Errorf(`unmarshal of vehicle failed - %w`, err)
		}

		ats = append(ats, &a)
	}

	return ats, nil
}

func (s *SmartContract) VehicleExists(ctx contractapi.TransactionContextInterface, id int) (bool, error) {
	aByts, err := ctx.GetStub().GetState(strconv.Itoa(id))
	if err != nil {
		return false, fmt.Errorf(`get state failed for vehicle %d - %w`, id, err)
	}

	return aByts != nil, nil
}

func (s *SmartContract) ChangeVehicleColour(ctx contractapi.TransactionContextInterface, id int, clr string) error {
	a, err := s.GetVehicle(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Color = clr
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

func (s *SmartContract) ChangeVehicleValue(ctx contractapi.TransactionContextInterface, id int, val int) error {
	a, err := s.GetVehicle(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Value = val
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

// Book functions

func (s *SmartContract) CreateBook(ctx contractapi.TransactionContextInterface, color string, id int, owner string, val int) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`create book failed - %w`, err)
	}

	if exists {
		return fmt.Errorf(`book with id %d already exists`, id)
	}

	book := Asset{
		Color: color,
		ID:    id,
		Owner: owner,
		Value: val,
	}

	aByts, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf(`marshal book failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(book.ID), aByts)
}

func (s *SmartContract) GetBook(ctx contractapi.TransactionContextInterface, id int) (*Asset, error) {
	aByts, err := ctx.GetStub().GetState(strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf(`get state failed for book %d - %w`, id, err)
	}

	if aByts == nil {
		return nil, fmt.Errorf(`book does not exist for id %d`, id)
	}

	var a Asset
	if err = json.Unmarshal(aByts, &a); err != nil {
		return nil, fmt.Errorf(`unmarshal book failed for book %d - %w`, id, err)
	}

	return &a, nil
}

func (s *SmartContract) UpdateBook(ctx contractapi.TransactionContextInterface, color string, id int, owner string, val int) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`checking book existence failed - %w`, err)
	}

	if !exists {
		return fmt.Errorf(`book with id %d does not exist`, id)
	}

	a := Asset{
		Color: color,
		ID:    id,
		Owner: owner,
		Value: val,
	}

	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal book failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(a.ID), aByts)
}

func (s *SmartContract) DeleteBook(ctx contractapi.TransactionContextInterface, id int) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`checking book existence failed - %w`, err)
	}

	if !exists {
		return fmt.Errorf(`book with id %d does not exist`, id)
	}

	return ctx.GetStub().DelState(strconv.Itoa(id))
}

func (s *SmartContract) TransferBook(ctx contractapi.TransactionContextInterface, id int, newOwner string) error {
	a, err := s.GetBook(ctx, id)
	if err != nil {
		return fmt.Errorf(`get book failed - %w`, err)
	}

	a.Owner = newOwner
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal book failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

func (s *SmartContract) GetAllBooks(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// ranging with empty start and end keys returns all books in chaincode namespace
	// an identical function of get all assets
	itr, err := ctx.GetStub().GetStateByRange(``, ``)
	if err != nil {
		return nil, fmt.Errorf(`get book state by range faied - %w`, err)
	}
	defer itr.Close()

	var ats []*Asset
	for itr.HasNext() {
		res, err := itr.Next()
		if err != nil {
			return nil, fmt.Errorf(`iterating next book query result failed - %w`, err)
		}

		var a Asset
		if err = json.Unmarshal(res.Value, &a); err != nil {
			return nil, fmt.Errorf(`unmarshal of book failed - %w`, err)
		}

		ats = append(ats, &a)
	}

	return ats, nil
}

func (s *SmartContract) BookExists(ctx contractapi.TransactionContextInterface, id int) (bool, error) {
	aByts, err := ctx.GetStub().GetState(strconv.Itoa(id))
	if err != nil {
		return false, fmt.Errorf(`get state failed for book %d - %w`, id, err)
	}

	return aByts != nil, nil
}

func (s *SmartContract) ChangeBookColour(ctx contractapi.TransactionContextInterface, id int, clr string) error {
	a, err := s.GetBook(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Color = clr
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

func (s *SmartContract) ChangeBookValue(ctx contractapi.TransactionContextInterface, id int, val int) error {
	a, err := s.GetBook(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Value = val
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

// House functions

func (s *SmartContract) CreateHouse(ctx contractapi.TransactionContextInterface, color string, id int, owner string, val int) error {
	exists, err := s.HouseExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`create house failed - %w`, err)
	}

	if exists {
		return fmt.Errorf(`house with id %d already exists`, id)
	}

	house := Asset{
		Color: color,
		ID:    id,
		Owner: owner,
		Value: val,
	}

	aByts, err := json.Marshal(house)
	if err != nil {
		return fmt.Errorf(`marshal house failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(house.ID), aByts)
}

func (s *SmartContract) GetHouse(ctx contractapi.TransactionContextInterface, id int) (*Asset, error) {
	aByts, err := ctx.GetStub().GetState(strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf(`get state failed for house %d - %w`, id, err)
	}

	if aByts == nil {
		return nil, fmt.Errorf(`house does not exist for id %d`, id)
	}

	var a Asset
	if err = json.Unmarshal(aByts, &a); err != nil {
		return nil, fmt.Errorf(`unmarshal house failed for house %d - %w`, id, err)
	}

	return &a, nil
}

func (s *SmartContract) UpdateHouse(ctx contractapi.TransactionContextInterface, color string, id int, owner string, val int) error {
	exists, err := s.HouseExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`checking house existence failed - %w`, err)
	}

	if !exists {
		return fmt.Errorf(`house with id %d does not exist`, id)
	}

	a := Asset{
		Color: color,
		ID:    id,
		Owner: owner,
		Value: val,
	}

	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal house failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(a.ID), aByts)
}

func (s *SmartContract) DeleteHouse(ctx contractapi.TransactionContextInterface, id int) error {
	exists, err := s.HouseExists(ctx, id)
	if err != nil {
		return fmt.Errorf(`checking house existence failed - %w`, err)
	}

	if !exists {
		return fmt.Errorf(`house with id %d does not exist`, id)
	}

	return ctx.GetStub().DelState(strconv.Itoa(id))
}

func (s *SmartContract) TransferHouse(ctx contractapi.TransactionContextInterface, id int, newOwner string) error {
	a, err := s.GetHouse(ctx, id)
	if err != nil {
		return fmt.Errorf(`get house failed - %w`, err)
	}

	a.Owner = newOwner
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal house failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

func (s *SmartContract) GetAllHouses(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// ranging with empty start and end keys returns all houses in chaincode namespace
	// an identical function of get all assets
	itr, err := ctx.GetStub().GetStateByRange(``, ``)
	if err != nil {
		return nil, fmt.Errorf(`get house state by range faied - %w`, err)
	}
	defer itr.Close()

	var ats []*Asset
	for itr.HasNext() {
		res, err := itr.Next()
		if err != nil {
			return nil, fmt.Errorf(`iterating next house query result failed - %w`, err)
		}

		var a Asset
		if err = json.Unmarshal(res.Value, &a); err != nil {
			return nil, fmt.Errorf(`unmarshal of house failed - %w`, err)
		}

		ats = append(ats, &a)
	}

	return ats, nil
}

func (s *SmartContract) HouseExists(ctx contractapi.TransactionContextInterface, id int) (bool, error) {
	aByts, err := ctx.GetStub().GetState(strconv.Itoa(id))
	if err != nil {
		return false, fmt.Errorf(`get state failed for house %d - %w`, id, err)
	}

	return aByts != nil, nil
}

func (s *SmartContract) ChangeHouseColour(ctx contractapi.TransactionContextInterface, id int, clr string) error {
	a, err := s.GetHouse(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Color = clr
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}

func (s *SmartContract) ChangeHouseValue(ctx contractapi.TransactionContextInterface, id int, val int) error {
	a, err := s.GetHouse(ctx, id)
	if err != nil {
		return fmt.Errorf(`get asset failed - %w`, err)
	}

	a.Value = val
	aByts, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf(`marshal asset failed - %w`, err)
	}

	return ctx.GetStub().PutState(strconv.Itoa(id), aByts)
}
