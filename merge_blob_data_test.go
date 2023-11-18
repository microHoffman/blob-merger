package main

import (
	crand "crypto/rand"
	"math/big"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func pickRandomNumber() int {
	mrand.Seed(time.Now().UnixNano())
	return mrand.Intn(MAX_BLOB_SIZE_IN_BYTES + 1)
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
	_, err := MergeBlobData(randomByteSlice)
	assert.NotNil(t, err, "Did not throw on blob data that are bigger than the max.")
}

func TestMergeBlobDataOkWithRandomValidInput(t *testing.T) {
	blobDatasCount := 20
	blobDatas := make([][]byte, blobDatasCount)
	for i := 0; i < blobDatasCount; i++ {
		randomByteSlice, _ := generateRandomByteArray(pickRandomNumber())
		blobDatas[i] = randomByteSlice
	}
	_, err := MergeBlobData(blobDatas...)
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
	result, err := MergeBlobData(allBlobs...)
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
	result, err := MergeBlobData(allBlobs...)
	assert.Nil(t, err, "Threw error on valid specific data.")
	assert.Len(t, result, 1, "Result does not have 1 blob.")
	assert.Len(t, result[0], 131000, "Blob should have 131000 length.")
}
