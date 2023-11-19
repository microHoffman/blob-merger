package main

import (
	"bytes"
	"fmt"
	"sort"
)

const (
	MAX_BLOB_SIZE_IN_BYTES = 1024 * 128 // 131072
)

// TODO add json rpc endpoint on suave-geth & call it

// TODO add the precompile function to solidity library
//  https://github.com/flashbots/suave-geth/blob/main/suave/sol/libraries/Suave.sol#L184

// TODO best blob sorting ->
//  1) prioritize maximalization of blob size
//  2) start by highest and check if we find any other blobs that we can merge together, however this means that the blobs could have less filled data after merge

func main() {
	// TODO anything to do here?
}

func blobDataLengthToBytes(length int) []byte {
	return []byte{
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}
}

func MergeBlobData(toAddresses [][]byte, blobs [][]byte) ([][]byte, error) {
	fmt.Println("=====")
	fmt.Println()
	fmt.Printf("MergeBlobData - number of separate blobs data: %v\n", len(blobs))
	fmt.Println("MergeBlobData - start")

	if len(toAddresses) != len(blobs) {
		return nil, fmt.Errorf("toAddresses and blobs parameter should have the same length")
	}

	for _, toAddress := range toAddresses {
		if len(toAddress) != 42 {
			return nil, fmt.Errorf("To address must have 42 length.")
		} else if string(toAddress) == "0x0000000000000000000000000000000000000000" { // TODO shall we rather use some functions from go-ethereum here?
			return nil, fmt.Errorf("To address can't be null in the blob tx.")
		}
	}

	for _, blobData := range blobs {
		if len(blobData) > MAX_BLOB_SIZE_IN_BYTES {
			return nil, fmt.Errorf("One of the blob data is longer than max allowed length of %d bytes", MAX_BLOB_SIZE_IN_BYTES)
		}
	}

	fmt.Printf("Size of each blob on inputc: [")
	for _, blobData := range blobs {
		fmt.Print(len(blobData))
		fmt.Print(",")
	}
	fmt.Printf("]")
	fmt.Println()

	fmt.Println("MergeBlobData - checks passed")

	// Sort blobs by length in descending order
	sort.Slice(blobs, func(i, j int) bool {
		return len(blobs[i]) > len(blobs[j])
	})

	var result [][]byte
	const (
		ADDRESS_SIZE          = 42
		BLOB_DATA_LENGTH_SIZE = 3
	)
	for len(blobs) > 0 {
		var mergedBlobSize int = 0
		var mergedBlobData []byte
		var removedBlobs [][]byte

		for i, iteratedBlob := range blobs {
			iteratedBlobSize := len(iteratedBlob) + ADDRESS_SIZE + BLOB_DATA_LENGTH_SIZE
			if (mergedBlobSize + iteratedBlobSize) <= MAX_BLOB_SIZE_IN_BYTES {
				// adding toAddress
				mergedBlobData = append(mergedBlobData, toAddresses[i]...)
				// adding length of the blob data, the length takes always 3 bytes
				mergedBlobData = append(mergedBlobData, blobDataLengthToBytes(iteratedBlobSize)...)
				// add the blob data itself
				mergedBlobData = append(mergedBlobData, iteratedBlob...)

				mergedBlobSize += iteratedBlobSize
				removedBlobs = append(removedBlobs, iteratedBlob)
			}
		}

		blobs = removeUsedBlobs(blobs, removedBlobs)
		result = append(result, mergedBlobData)
	}

	// Sort blobs by length in descending order
	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) > len(result[j])
	})

	fmt.Println("MergeBlobData - end")
	fmt.Printf("MergeBlobData - number of merged blobs: %v\n", len(result))

	return result, nil
}

func removeUsedBlobs(allBlobs, blobsToRemove [][]byte) [][]byte {
	var result [][]byte

OuterLoop:
	for _, s := range allBlobs {
		for _, r := range blobsToRemove {
			if bytes.Equal(s, r) {
				continue OuterLoop
			}
		}
		result = append(result, s)
	}

	return result
}
