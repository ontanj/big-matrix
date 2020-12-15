package genmatrix

import (
    "math/big"
    "github.com/niclabs/tcpaillier"
    "errors"
)

type DJ_public_key struct {
    *tcpaillier.PubKey
}

func (pk DJ_public_key) Add(a, b interface{}) (sum interface{}, err error) {
    err = assertBigint(a, b)
    if err != nil {return}
    return pk.PubKey.Add(a.(*big.Int), b.(*big.Int))
}

func (pk DJ_public_key) Subtract(a, b interface{}) (diff interface{}, err error) {
    err = assertBigint(a, b)
    if err != nil {return}
    neg, _, err := pk.PubKey.Multiply(b.(*big.Int), new(big.Int).SetInt64(-1))
    if err != nil {return nil, err}
    return pk.PubKey.Add(a.(*big.Int), neg)
}

func (pk DJ_public_key) Scale(ciphertext, factor interface{}) (product interface{}, err error) {
    err = assertBigint(ciphertext, factor)
    if err != nil {return}
    product, _, err = pk.PubKey.Multiply(ciphertext.(*big.Int), factor.(*big.Int))
    return
}

func (pk DJ_public_key) Multiply(a, b interface{}) (interface{}, error) {
    err := assertBigint(a, b)
    if err != nil {return nil, err}
    return nil, errors.New("multiplication not supported")
}

func (pk DJ_public_key) Scalarspace() bool {
    return false
}

func NewDJCryptosystem() (public_key DJ_public_key, secret_keys []*tcpaillier.KeyShare, err error) {
    secret_keys, djpk, err := tcpaillier.NewKey(128, 1, 3, 3)
    if err != nil {return}
    public_key = DJ_public_key{djpk}
    return
}