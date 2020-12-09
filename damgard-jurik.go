package genmatrix

import (
    "math/big"
    "github.com/niclabs/tcpaillier"
    "errors"
)

type dj_public_key struct {
    *tcpaillier.PubKey
}

func (pk dj_public_key) Add(a, b interface{}) (sum interface{}, err error) {
    return pk.PubKey.Add(a.(*big.Int), b.(*big.Int))
}

func (pk dj_public_key) Subtract(a, b interface{}) (diff interface{}, err error) {
    neg, _, err := pk.PubKey.Multiply(b.(*big.Int), new(big.Int).SetInt64(-1))
    if err != nil {return nil, err}
    return pk.PubKey.Add(a.(*big.Int), neg)
}

func (pk dj_public_key) Scale(ciphertext, factor interface{}) (product interface{}, err error) {
    product, _, err = pk.PubKey.Multiply(ciphertext.(*big.Int), factor.(*big.Int))
    return
}

func (pk dj_public_key) Multiply(a, b interface{}) (interface{}, error) {
    return nil, errors.New("multiplication not supported")
}

func (pk dj_public_key) Scalarspace() bool {
    return false
}

func NewDJCryptosystem() (public_key dj_public_key, secret_keys []*tcpaillier.KeyShare, err error) {
    secret_keys, djpk, err := tcpaillier.NewKey(128, 1, 3, 3)
    if err != nil {return}
    public_key = dj_public_key{djpk}
    return
}