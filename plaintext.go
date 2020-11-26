package bigmatrix

import (
    "math/big"
)

type plain struct {}

func (p plain) Add(terms ...*big.Int) (*big.Int, error) {
    sum := big.NewInt(0)
    for _, term := range terms {
        sum.Add(sum, term)
    }
    return sum, nil
}

func (p plain) MultiplyFactor(ciphertext *big.Int, plaintext *big.Int) (*big.Int, error) {
    return p.Multiply(ciphertext, plaintext)
}

func (p plain) Multiply(terms ...*big.Int) (*big.Int, error) {
    product := big.NewInt(1)
    for _, term := range terms {
        product.Mul(product, term)
    }
    return product, nil
}