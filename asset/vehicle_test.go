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
	testVehicle = Asset{Color: "brown", ID: 88, Owner: "Arnold", Value: 989}
)

func TestSmartContractCreateVehicle(t *testing.T) {
	stub := newMockStub()
	testCreateVehicle(stub, t)

	in := marshalVehicle()
	out := getVehicleState(stub, testVehicle.ID, t)

	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractUpdateVehicle(t *testing.T) {
	stub := newMockStub()
	testCreateVehicle(stub, t)

	testVehicle.Color = clrBlue
	testVehicle.Value = 1500
	testUpdateVehicle(stub, t)

	in := marshalVehicle()
	out := getVehicleState(stub, testVehicle.ID, t)

	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractGetVehicle(t *testing.T) {
	stub := newMockStub()
	testCreateVehicle(stub, t)

	res := stub.MockInvoke(`1`, [][]byte{[]byte("GetVehicle"), []byte(strconv.Itoa(testVehicle.ID))})
	if res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	in := marshalVehicle()
	if !bytes.Equal(in, res.Payload) {
		t.Fatalf(errExpect, in, res.Payload)
	}
}

func TestSmartContractGetAllVehicles(t *testing.T) {
	stub := newMockStub()
	testInitLedger(stub, t)

	res := stub.MockInvoke(`2`, [][]byte{[]byte("GetAllVehicles")})
	if res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	in := marshalVehicles()
	if !bytes.Equal(in, res.Payload) {
		t.Fatalf(errExpect, in, res.Payload)
	}
}

func TestSmartContractDeleteVehicle(t *testing.T) {
	stub := newMockStub()
	testCreateVehicle(stub, t)

	if res := stub.MockInvoke(`6`, [][]byte{
		[]byte("DeleteVehicle"), []byte(strconv.Itoa(testVehicle.ID)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	out := getVehicleState(stub, testVehicle.ID, t)
	if out != nil {
		t.Fatalf(`state db should be empty after mock delete (%s)`, string(out))
	}
}

func TestSmartContractTransferVehicle(t *testing.T) {
	stub := newMockStub()
	testCreateVehicle(stub, t)

	if res := stub.MockInvoke(`7`, [][]byte{
		[]byte("TransferVehicle"), []byte(strconv.Itoa(testVehicle.ID)), []byte(ownrDavid),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testVehicle.Owner = ownrDavid
	out := getVehicleState(stub, testVehicle.ID, t)
	in := marshalVehicle()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractChangeVehicleColour(t *testing.T) {
	stub := newMockStub()
	testCreateVehicle(stub, t)

	if res := stub.MockInvoke(`77`, [][]byte{
		[]byte("ChangeVehicleColour"), []byte(strconv.Itoa(testVehicle.ID)), []byte(clrBrown),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testVehicle.Color = clrBrown
	out := getVehicleState(stub, testVehicle.ID, t)
	in := marshalVehicle()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractChangeVehicleValue(t *testing.T) {
	stub := newMockStub()
	testCreateVehicle(stub, t)

	if res := stub.MockInvoke(`77`, [][]byte{
		[]byte("ChangeVehicleValue"), []byte(strconv.Itoa(testVehicle.ID)), []byte("888"),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testVehicle.Value = 888
	out := getVehicleState(stub, testVehicle.ID, t)
	in := marshalVehicle()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func testCreateVehicle(stub *shimtest.MockStub, t *testing.T) {
	if res := stub.MockInvoke(`4`, [][]byte{
		[]byte("CreateVehicle"), []byte(testVehicle.Color), []byte(strconv.Itoa(testVehicle.ID)), []byte(testVehicle.Owner), []byte(strconv.Itoa(testVehicle.Value)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}
}

func testUpdateVehicle(stub *shimtest.MockStub, t *testing.T) {
	if res := stub.MockInvoke(`5`, [][]byte{
		[]byte("UpdateVehicle"), []byte(testVehicle.Color), []byte(strconv.Itoa(testVehicle.ID)), []byte(testVehicle.Owner), []byte(strconv.Itoa(testVehicle.Value)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}
}

func getVehicleState(stub *shimtest.MockStub, id int, t *testing.T) []byte {
	out, err := stub.GetState(strconv.Itoa(id))
	if err != nil {
		t.Fatalf("failed to retrieve asset info - %s", err.Error())
	}

	return out
}

func marshalVehicle() []byte {
	byts, err := json.Marshal(testVehicle)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to marshal asset info - %s", err.Error()))
	}

	return byts
}

func marshalVehicles() []byte {
	byts, err := json.Marshal(assets)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to marshal assets - %s", err.Error()))
	}

	return byts
}
