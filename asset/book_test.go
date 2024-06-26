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
	testBook = Asset{Color: "brown", ID: 88, Owner: "Arnold", Value: 989}
)

func marshalBook() []byte {
	byts, err := json.Marshal(testBook)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to marshal asset info - %s", err.Error()))
	}

	return byts
}

func getBookState(stub *shimtest.MockStub, id int, t *testing.T) []byte {
	out, err := stub.GetState(strconv.Itoa(id))
	if err != nil {
		t.Fatalf("failed to retrieve asset info - %s", err.Error())
	}

	return out
}

func testCreateBook(stub *shimtest.MockStub, t *testing.T) {
	if res := stub.MockInvoke(`4`, [][]byte{
		[]byte("CreateBook"), []byte(testBook.Color), []byte(strconv.Itoa(testBook.ID)), []byte(testBook.Owner), []byte(strconv.Itoa(testBook.Value)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}
}

func marshalBooks() []byte {
	byts, err := json.Marshal(assets)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to marshal assets - %s", err.Error()))
	}

	return byts
}

func testUpdateBook(stub *shimtest.MockStub, t *testing.T) {
	if res := stub.MockInvoke(`5`, [][]byte{
		[]byte("UpdateBook"), []byte(testBook.Color), []byte(strconv.Itoa(testBook.ID)), []byte(testBook.Owner), []byte(strconv.Itoa(testBook.Value)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}
}

func TestSmartContractGetBook(t *testing.T) {
	stub := newMockStub()
	testCreateBook(stub, t)

	res := stub.MockInvoke(`1`, [][]byte{[]byte("GetBook"), []byte(strconv.Itoa(testBook.ID))})
	if res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	in := marshalBook()
	if !bytes.Equal(in, res.Payload) {
		t.Fatalf(errExpect, in, res.Payload)
	}
}

func TestSmartContractUpdateBook(t *testing.T) {
	stub := newMockStub()
	testCreateBook(stub, t)

	testBook.Color = clrBlue
	testBook.Value = 1500
	testUpdateBook(stub, t)

	in := marshalBook()
	out := getBookState(stub, testBook.ID, t)

	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractChangeBookValue(t *testing.T) {
	stub := newMockStub()
	testCreateBook(stub, t)

	if res := stub.MockInvoke(`77`, [][]byte{
		[]byte("ChangeBookValue"), []byte(strconv.Itoa(testBook.ID)), []byte("888"),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testBook.Value = 888
	out := getBookState(stub, testBook.ID, t)
	in := marshalBook()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractChangeBookColour(t *testing.T) {
	stub := newMockStub()
	testCreateBook(stub, t)

	if res := stub.MockInvoke(`77`, [][]byte{
		[]byte("ChangeBookColour"), []byte(strconv.Itoa(testBook.ID)), []byte(clrBrown),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testBook.Color = clrBrown
	out := getBookState(stub, testBook.ID, t)
	in := marshalBook()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractGetAllBooks(t *testing.T) {
	stub := newMockStub()
	testInitLedger(stub, t)

	res := stub.MockInvoke(`2`, [][]byte{[]byte("GetAllBooks")})
	if res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	in := marshalBooks()
	if !bytes.Equal(in, res.Payload) {
		t.Fatalf(errExpect, in, res.Payload)
	}
}

func TestSmartContractDeleteBook(t *testing.T) {
	stub := newMockStub()
	testCreateBook(stub, t)

	if res := stub.MockInvoke(`6`, [][]byte{
		[]byte("DeleteBook"), []byte(strconv.Itoa(testBook.ID)),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	out := getBookState(stub, testBook.ID, t)
	if out != nil {
		t.Fatalf(`state db should be empty after mock delete (%s)`, string(out))
	}
}

func TestSmartContractCreateBook(t *testing.T) {
	stub := newMockStub()
	testCreateBook(stub, t)

	in := marshalBook()
	out := getBookState(stub, testBook.ID, t)

	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}

func TestSmartContractTransferBook(t *testing.T) {
	stub := newMockStub()
	testCreateBook(stub, t)

	if res := stub.MockInvoke(`7`, [][]byte{
		[]byte("TransferBook"), []byte(strconv.Itoa(testBook.ID)), []byte(ownrDavid),
	}); res.Status != shim.OK {
		t.Fatalf(errOK, res.Status, res.Message)
	}

	testBook.Owner = ownrDavid
	out := getBookState(stub, testBook.ID, t)
	in := marshalBook()
	if !bytes.Equal(in, out) {
		t.Fatalf(errExpect, in, out)
	}
}
