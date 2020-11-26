package bigmatrix

import (
    "math/big"
)

type plain struct {}

func (p plain) Add(a, b *big.Int) (*big.Int, error) {
    return new(big.Int).Add(a, b), nil
}

func (p plain) MultiplyFactor(ciphertext *big.Int, plaintext *big.Int) (*big.Int, error) {
    return p.Multiply(ciphertext, plaintext)
}

func (p plain) Multiply(a, b *big.Int) (*big.Int, error) {
    return new(big.Int).Mul(a, b), nil
}