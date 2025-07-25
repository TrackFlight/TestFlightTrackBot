package utils

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"errors"
)

var aesKey = []byte{
	0x3f, 0xa7, 0xc9, 0x1e,
	0x5b, 0x4d, 0x72, 0x9f,
	0x2d, 0x6e, 0x8a, 0xbc,
	0xef, 0x01, 0x23, 0x45,
}

func EncryptBlockECB(block, key []byte) ([]byte, error) {
	if len(block) != aes.BlockSize {
		return nil, errors.New("invalid block size")
	}
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, aes.BlockSize)
	cipher.Encrypt(out, block)
	return out, nil
}

func DecryptBlockECB(block, key []byte) ([]byte, error) {
	if len(block) != aes.BlockSize {
		return nil, errors.New("invalid block size")
	}
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, aes.BlockSize)
	cipher.Decrypt(out, block)
	return out, nil
}

func EncodeName(id int64) string {
	blk := make([]byte, aes.BlockSize)
	binary.BigEndian.PutUint64(blk[:8], uint64(id))
	ct, err := EncryptBlockECB(blk, aesKey)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(ct)
}

func DecodeName(encoded string) (int64, error) {
	ct, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return 0, err
	}
	if len(ct) != aes.BlockSize {
		return 0, errors.New("invalid encoded length")
	}
	pt, err := DecryptBlockECB(ct, aesKey)
	if err != nil {
		return 0, err
	}
	id := int64(binary.BigEndian.Uint64(pt[:8]))
	return id, nil
}
