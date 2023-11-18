package main

import (
	crand "crypto/rand"
	"math/big"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	TEST_ADDRESS_ONE   = "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
	TEST_ADDRESS_TWO   = "0x758b8178A9A4B7206D1f648c4a77C515CbaC7000"
	TEST_ADDRESS_THREE = "0x814FaE9f487206471B6B0D713cD51a2D35980000"
	TEST_ADDRESS_FOUR  = "0x763c396673F9c391DCe3361A9A71C8E161388000"
	TEST_ADDRESS_FIVE  = "0xd4E96eF8eee8678dBFf4d535E033Ed1a4F7605b7"
)

var ALL_TEST_ADDRESSES []string

func init() {
	ALL_TEST_ADDRESSES = []string{
		TEST_ADDRESS_ONE,
		TEST_ADDRESS_TWO,
		TEST_ADDRESS_THREE,
		TEST_ADDRESS_FOUR,
		TEST_ADDRESS_FIVE,
	}
}

func pickRandomNumber(maxNumber int) int {
	mrand.Seed(time.Now().UnixNano())
	return mrand.Intn(maxNumber)
}

func pickRandomBlobDataSize() int {
	return pickRandomNumber(MAX_BLOB_SIZE_IN_BYTES + 1)
}

func generateRandomByteArray(length int) ([]byte, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)

	for i := range randomBytes {
		n, err := crand.Int(crand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return nil, err
		}
		randomBytes[i] = charset[n.Int64()]
	}

	return randomBytes, nil
}

func TestMergeBlobDataThrowOnInvalidInput(t *testing.T) {
	randomByteSlice, _ := generateRandomByteArray(MAX_BLOB_SIZE_IN_BYTES + 1)
	_, err := MergeBlobData([][]byte{[]byte(TEST_ADDRESS_ONE)}, [][]byte{randomByteSlice})
	assert.NotNil(t, err, "Did not throw on blob data that are bigger than the max.")
}

func TestMergeBlobDataOkWithRandomValidInput(t *testing.T) {
	blobDatasCount := 20
	blobDatas := make([][]byte, blobDatasCount)
	toAddresses := make([][]byte, blobDatasCount)
	for i := 0; i < blobDatasCount; i++ {
		randomByteSlice, _ := generateRandomByteArray(pickRandomBlobDataSize())
		blobDatas[i] = randomByteSlice
		toAddresses[i] = []byte(ALL_TEST_ADDRESSES[pickRandomNumber(len(ALL_TEST_ADDRESSES))])
	}
	_, err := MergeBlobData(toAddresses, blobDatas)
	assert.Nil(t, err, "Incorrectly throwed error on blob data that are valid.")
}

func TestMergeBlobDataWithMultipleBlobsResult(t *testing.T) {
	// max blob size: 131072
	// [100000, 30000, 120000, 1000, 35000] -> [130000, 121000, 35000]
	// flow:
	//  1) take highest number (120000), check if we can find any blob we can merge to
	//  2) we can merge 120000 + 1000 => 121000 blob
	//  3) remove 120000 + 1000 blobs
	//  4) take next remaining highest number (100000), check if we can find any blob we can merge to
	//  5) we can merge 100000 + 30000 => 130000 blob
	//  6) remove 100000 + 30000 blobs
	//  7) only 35000 blob remains, add a new blob
	//  8) done
	a, _ := generateRandomByteArray(100000)
	b, _ := generateRandomByteArray(30000)
	c, _ := generateRandomByteArray(120000)
	d, _ := generateRandomByteArray(1000)
	e, _ := generateRandomByteArray(35000)
	allBlobs := [][]byte{a, b, c, d, e}
	allToAddresses := [][]byte{
		[]byte(TEST_ADDRESS_ONE),
		[]byte(TEST_ADDRESS_TWO),
		[]byte(TEST_ADDRESS_THREE),
		[]byte(TEST_ADDRESS_FOUR),
		[]byte(TEST_ADDRESS_FIVE),
	}
	result, err := MergeBlobData(allToAddresses, allBlobs)
	assert.Nil(t, err, "Threw error on valid specific data.")
	assert.Equal(t, 3, len(result), "Result does not have 3 blobs.")
	assert.Equal(t, 130000, len(result[0]), "Largest blob should have 131000 length.")
	assert.Equal(t, 121000, len(result[1]), "Second largest blob should have 120000 length.")
	assert.Equal(t, 35000, len(result[2]), "Third largest blob should have 35000 length.")
}

func TestMergeBlobDataWithSingleBlobResult(t *testing.T) {
	// max blob size: 131072
	// [100000, 30000, 1000] -> [131000]
	a, _ := generateRandomByteArray(100000)
	b, _ := generateRandomByteArray(30000)
	c, _ := generateRandomByteArray(1000)
	allBlobs := [][]byte{a, b, c}
	allToAddresses := [][]byte{
		[]byte(TEST_ADDRESS_ONE),
		[]byte(TEST_ADDRESS_TWO),
		[]byte(TEST_ADDRESS_THREE),
	}
	result, err := MergeBlobData(allToAddresses, allBlobs)
	assert.Nil(t, err, "Threw error on valid specific data.")
	assert.Len(t, result, 1, "Result does not have 1 blob.")
	assert.Len(t, result[0], 131000, "Blob should have 131000 length.")
}
