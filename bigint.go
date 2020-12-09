package genmatrix

import (
    "math/big"
)

type bigint struct {}

func (p bigint) Add(a, b interface{}) (interface{}, error) {
    return new(big.Int).Add(a.(*big.Int), b.(*big.Int)), nil
}

func (p bigint) Subtract(a, b interface{}) (interface{}, error) {
    return new(big.Int).Sub(a.(*big.Int), b.(*big.Int)), nil
}

func (p bigint) Multiply(a, b interface{}) (interface{}, error) {
    return new(big.Int).Mul(a.(*big.Int), b.(*big.Int)), nil
}

func (p bigint) Scale(ciphertext interface{}, plaintext interface{}) (interface{}, error) {
    return p.Multiply(ciphertext.(*big.Int), plaintext.(*big.Int))
}

func (p bigint) Scalarspace() bool {
    return true
}

// create a new Matrix from int values
func NewMatrixFromInt(rows, cols int, data []int) (Matrix, error) {
    if data == nil {
        return NewMatrix(rows, cols, nil, bigint{})
    }
    l := len(data)
    s := make([]interface{}, l)
    for i := 0; i < l; i += 1 {
        s[i] = big.NewInt(int64(data[i]))
    }
    return NewMatrix(rows, cols, s, bigint{})
}