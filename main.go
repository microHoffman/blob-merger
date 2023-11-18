package main

import (
	"bytes"
	"fmt"
	"sort"
)

const (
	MAX_BLOB_SIZE_IN_BYTES = 1024 * 128 // 131072
)

// TODO add offsets + test

// TODO add json rpc endpoint on suave-geth & call it

// TODO add the precompile function to solidity library
//  https://github.com/flashbots/suave-geth/blob/main/suave/sol/libraries/Suave.sol#L184

// TODO best blob sorting ->
//  1) prioritize maximalization of blob size
//  2) start by highest and check if we find any other blobs that we can merge together, however this means that the blobs could have less filled data after merge

func main() {
	// TODO anything to do here?
}

func MergeBlobData(blobs ...[]byte) ([][]byte, error) {
	for _, blobData := range blobs {
		if len(blobData) > MAX_BLOB_SIZE_IN_BYTES {
			return nil, fmt.Errorf("One of the blob data is longer than max allowed length of %d bytes", MAX_BLOB_SIZE_IN_BYTES)
		}
	}

	// Sort blobs by length in descending order
	sort.Slice(blobs, func(i, j int) bool {
		return len(blobs[i]) > len(blobs[j])
	})

	var result [][]byte
	for len(blobs) > 0 {
		var blobSize int = 0
		var blobData []byte
		var removedBlobs [][]byte

		for _, blob := range blobs {
			if (blobSize + len(blob)) <= MAX_BLOB_SIZE_IN_BYTES {
				blobSize += len(blob)
				removedBlobs = append(removedBlobs, blob)
				blobData = append(blobData, blob...)
			}
		}

		// Remove processed blobs from stackOfBlobs
		blobs = subtractSets(blobs, removedBlobs)
		result = append(result, blobData)
	}

	// Sort blobs by length in descending order
	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) > len(result[j])
	})

	return result, nil
}

func subtractSets(slice, toRemove [][]byte) [][]byte {
	var result [][]byte

OuterLoop:
	for _, s := range slice {
		for _, r := range toRemove {
			if bytes.Equal(s, r) {
				continue OuterLoop
			}
		}
		result = append(result, s)
	}

	return result
}
