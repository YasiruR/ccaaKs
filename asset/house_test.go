package asset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/tryfix/log"
	"strconv"
	"testing"
)

var (
	testHouse = Asset{Color: "brown", ID: 88, Owner: "Arnold", Value: 989}
)

func TestSmartContractGetAllHouses(t *testing.T) {
	stub := newMockStub()
	testInitLedger(stub, t)

	res := stub.MockInvoke(`2`, [][]byte{[]byte("GetAllHouses")})
	if res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	in := marshalHouses()
	if !bytes.Equal(in, res.Payload) {
		t.Fatalf(errExpect, in, res.Payload)
	}
}

func TestSmartContractGetHouse(t *testing.T) {
	stub := newMockStub()
	testCreateHouse(stub, t)

	res := stub.MockInvoke(`1`, [][]byte{[]byte("GetHouse"), []byte(strconv.Itoa(testHouse.ID))})
	if res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	in := marshalHouse()
	if !bytes.Equal(in, res.Payload) {
		t.Fatalf(errExpect, in, res.Payload)
	}
}

func TestSmartContractUpdateHouse(t *testing.T) {
	stub := newMockStub()
	testCreateHouse(stub, t)

	testHouse.Color = clrBlue
	testHouse.Value = 1500
	testUpdateHouse(stub, t)

	in := marshalHouse()
	out := getHouseState(stub, testHouse.ID, t)

	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractTransferHouse(t *testing.T) {
	stub := newMockStub()
	testCreateHouse(stub, t)

	if res := stub.MockInvoke(`7`, [][]byte{
		[]byte("TransferHouse"), []byte(strconv.Itoa(testHouse.ID)), []byte(ownrDavid),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testHouse.Owner = ownrDavid
	out := getHouseState(stub, testHouse.ID, t)
	in := marshalHouse()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractChangeHouseColour(t *testing.T) {
	stub := newMockStub()
	testCreateHouse(stub, t)

	if res := stub.MockInvoke(`77`, [][]byte{
		[]byte("ChangeHouseColour"), []byte(strconv.Itoa(testHouse.ID)), []byte(clrBrown),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testHouse.Color = clrBrown
	out := getHouseState(stub, testHouse.ID, t)
	in := marshalHouse()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractCreateHouse(t *testing.T) {
	stub := newMockStub()
	testCreateHouse(stub, t)

	in := marshalHouse()
	out := getHouseState(stub, testHouse.ID, t)

	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractChangeHouseValue(t *testing.T) {
	stub := newMockStub()
	testCreateHouse(stub, t)

	if res := stub.MockInvoke(`77`, [][]byte{
		[]byte("ChangeHouseValue"), []byte(strconv.Itoa(testHouse.ID)), []byte("888"),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testHouse.Value = 888
	out := getHouseState(stub, testHouse.ID, t)
	in := marshalHouse()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractDeleteHouse(t *testing.T) {
	stub := newMockStub()
	testCreateHouse(stub, t)

	if res := stub.MockInvoke(`6`, [][]byte{
		[]byte("DeleteHouse"), []byte(strconv.Itoa(testHouse.ID)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	out := getHouseState(stub, testHouse.ID, t)
	if out != nil {
		t.Fatalf(`state db should be empty after mock delete (%s)`, string(out))
	}
}

func getHouseState(stub *shimtest.MockStub, id int, t *testing.T) []byte {
	out, err := stub.GetState(strconv.Itoa(id))
	if err != nil {
		t.Fatalf("failed to retrieve asset info - %s", err.Error())
	}

	return out
}

func marshalHouses() []byte {
	byts, err := json.Marshal(assets)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to marshal assets - %s", err.Error()))
	}

	return byts
}

func testCreateHouse(stub *shimtest.MockStub, t *testing.T) {
	if res := stub.MockInvoke(`4`, [][]byte{
		[]byte("CreateHouse"), []byte(testHouse.Color), []byte(strconv.Itoa(testHouse.ID)), []byte(testHouse.Owner), []byte(strconv.Itoa(testHouse.Value)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}
}

func marshalHouse() []byte {
	byts, err := json.Marshal(testHouse)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to marshal asset info - %s", err.Error()))
	}

	return byts
}

func testUpdateHouse(stub *shimtest.MockStub, t *testing.T) {
	if res := stub.MockInvoke(`5`, [][]byte{
		[]byte("UpdateHouse"), []byte(testHouse.Color), []byte(strconv.Itoa(testHouse.ID)), []byte(testHouse.Owner), []byte(strconv.Itoa(testHouse.Value)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}
}
