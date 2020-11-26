package bigmatrix

import (
    "math/big"
)

type cryptosystem interface {
    Add(*big.Int, *big.Int) (*big.Int, error)
    MultiplyFactor(ciphertext *big.Int, plaintext *big.Int) (*big.Int, error)
    Multiply(*big.Int, *big.Int) (*big.Int, error)
}
