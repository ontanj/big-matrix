package bigmatrix

import (
	"math/big"
	"github.com/niclabs/tcpaillier"
)

type dj_public_key struct {
    *tcpaillier.PubKey
}

func (pk dj_public_key) Add(terms ...*big.Int) (sum *big.Int, err error) {
    return pk.PubKey.Add(terms...)
}

func (pk dj_public_key) MultiplyFactor(ciphertext, constant *big.Int) (product *big.Int, err error) {
    product, _, err = pk.PubKey.Multiply(ciphertext, constant)
    return
}

func (pk dj_public_key) Multiply(...*big.Int) (*big.Int, error) {
    panic("Not supported for Damgård-Jurik cryptosystem.")
}

func NewDJCryptosystem() (public_key dj_public_key, secret_keys []*tcpaillier.KeyShare, err error) {
    secret_keys, djpk, err := tcpaillier.NewKey(128, 1, 3, 3)
    if err != nil {return}
    public_key = dj_public_key{djpk}
    return
}