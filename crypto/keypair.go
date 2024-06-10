package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"math/big"

	"go-blockchain/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{
		key: key,
	}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.key.PublicKey,
	}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	S, R *big.Int
}

func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}

// Serialization methods

func (k PrivateKey) GobEncode() ([]byte, error) {
	x := k.key.X.Bytes()
	y := k.key.Y.Bytes()
	d := k.key.D.Bytes()
	return encodePrivateKey(x, y, d)
}

func (k *PrivateKey) GobDecode(data []byte) error {
	x, y, d, err := decodePrivateKey(data)
	if err != nil {
		return err
	}
	k.key = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		},
		D: d,
	}
	return nil
}

func (k PublicKey) GobEncode() ([]byte, error) {
	x := k.Key.X.Bytes()
	y := k.Key.Y.Bytes()
	return encodePublicKey(x, y)
}

func (k *PublicKey) GobDecode(data []byte) error {
	x, y, err := decodePublicKey(data)
	if err != nil {
		return err
	}
	k.Key = &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return nil
}

func encodePrivateKey(x, y, d []byte) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(x); err != nil {
		return nil, err
	}
	if err := encoder.Encode(y); err != nil {
		return nil, err
	}
	if err := encoder.Encode(d); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func decodePrivateKey(data []byte) (*big.Int, *big.Int, *big.Int, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	var xBytes, yBytes, dBytes []byte
	if err := decoder.Decode(&xBytes); err != nil {
		return nil, nil, nil, err
	}
	if err := decoder.Decode(&yBytes); err != nil {
		return nil, nil, nil, err
	}
	if err := decoder.Decode(&dBytes); err != nil {
		return nil, nil, nil, err
	}
	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)
	d := new(big.Int).SetBytes(dBytes)
	return x, y, d, nil
}

func encodePublicKey(x, y []byte) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(x); err != nil {
		return nil, err
	}
	if err := encoder.Encode(y); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func decodePublicKey(data []byte) (*big.Int, *big.Int, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	var xBytes, yBytes []byte
	if err := decoder.Decode(&xBytes); err != nil {
		return nil, nil, err
	}
	if err := decoder.Decode(&yBytes); err != nil {
		return nil, nil, err
	}
	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)
	return x, y, nil
}
