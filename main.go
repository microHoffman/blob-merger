package main

const (
	MAX_BLOB_SIZE_IN_BYTES = 1024 * 128
)

func Merge_blob_data(blobs ...[]byte) [][]byte {
	var result [][]byte
	var blob []byte
	result = append(result, blob)
	for _, blob_data := range blobs {
		var flag bool = false
		if len(blob_data) > MAX_BLOB_SIZE_IN_BYTES {
			continue
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
	return result
}
