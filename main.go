package main

import (
	"fmt"
	"sort"
)

const (
	MAX_BLOB_SIZE_IN_BYTES = 1024 * 128 // 131072
)

func main() {
	// anything to do here?
}

func MergeBlobData(blobs ...[]byte) ([][]byte, error) {
	var result [][]byte
	var blob []byte
	result = append(result, blob)

	// TODO does this work fine?
	// TODO what to do in case of the same length of the []byte items in [][]byte?
	sort.Slice(blobs, func(i, j int) bool {
		return len(blobs[i]) > len(blobs[j])
	})

	for _, blob_data := range blobs {
		var flag bool = false
		if len(blob_data) > MAX_BLOB_SIZE_IN_BYTES {
			return nil, fmt.Errorf("One of the blob data is longer than max allowed length of %d bytes", MAX_BLOB_SIZE_IN_BYTES)
		}
		for i, blob := range result {
			if len(blob)+len(blob_data) <= MAX_BLOB_SIZE_IN_BYTES {
				flag = true
				result[i] = append(blob, blob_data...)
				break
			}
		}
		if !flag {
			var blob []byte
			copy(blob, blob_data)
			result = append(result, blob)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) > len(result[j])
	})

	return result, nil
}
