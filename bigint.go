package genmatrix

import (
    "math/big"
    "fmt"
    "reflect"
)

type Bigint struct {}

func assertBigint(a, b interface{}) error {
    if reflect.TypeOf(a) != reflect.TypeOf(new(big.Int)) {
        return fmt.Errorf("first operator is not *big.Int, but %T", a)
    } 
    if reflect.TypeOf(b) != reflect.TypeOf(new(big.Int)) {
        return fmt.Errorf("second operator is not *big.Int, but %T", b)    
    }
    return nil
}

func (p Bigint) Add(a, b interface{}) (interface{}, error) {
    err := assertBigint(a, b)
    if err != nil {return nil, err}
    return new(big.Int).Add(a.(*big.Int), b.(*big.Int)), nil
}

func (p Bigint) Subtract(a, b interface{}) (interface{}, error) {
    err := assertBigint(a, b)
    if err != nil {return nil, err}
    return new(big.Int).Sub(a.(*big.Int), b.(*big.Int)), nil
}

func (p Bigint) Negate(a interface{}) (interface{}, error) {
    return new(big.Int).Mul(big.NewInt(-1), a.(*big.Int)), nil
}

func (p Bigint) Multiply(a, b interface{}) (interface{}, error) {
    err := assertBigint(a, b)
    if err != nil {return nil, err}
    return new(big.Int).Mul(a.(*big.Int), b.(*big.Int)), nil
}

func (p Bigint) Scale(a interface{}, b interface{}) (interface{}, error) {
    err := assertBigint(a, b)
    if err != nil {return nil, err}
    return p.Multiply(a.(*big.Int), b.(*big.Int))
}

func (p Bigint) Scalarspace() bool {
    return true
}

// create a new Matrix from int values
func NewMatrixFromInt(rows, cols int, data []int) (Matrix, error) {
    if data == nil {
        return NewMatrix(rows, cols, nil, Bigint{})
    }
    l := len(data)
    s := make([]interface{}, l)
    for i := 0; i < l; i += 1 {
        s[i] = big.NewInt(int64(data[i]))
    }
    return NewMatrix(rows, cols, s, Bigint{})
}

// create a new polynomial from int values
func NewPolyFromInt(data []int) (Polynomial, error) {
    l := len(data)
    s := make([]interface{}, l)
    for i := 0; i < l; i += 1 {
        s[i] = big.NewInt(int64(data[i]))
    }
    return NewPolynomial(s, Bigint{})
}