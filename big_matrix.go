package bigmatrix

import (
    "math/big"
    "fmt"
    "reflect"
)

type BigMatrix struct {
    values []*big.Int
    rows, cols int
    cryptosystem cryptosystem
}

// create a new BigMatrix with the given size and data, possible encrypted by cryptosystem cs (or nil)
func NewBigMatrix(rows, cols int, data []*big.Int, cs cryptosystem) BigMatrix {
    if data == nil {
        data = make([]*big.Int, rows*cols)
        for i := range data {
            data[i] = big.NewInt(0)
        }
    } else if rows * cols != len(data) {
        panic("Data structure not matching matrix size")
    }
    var m BigMatrix
    m.values = data
    m.rows = rows
    m.cols = cols
    if cs == nil {
        m.cryptosystem = plain{}
    } else {
        m.cryptosystem = cs
    }
    return m
}

// create a new BigMatrix from int values
func NewBigMatrixFromInt(rows, cols int, data []int) BigMatrix {
    if data == nil {
        return NewBigMatrix(rows, cols, nil, nil)
    }
    l := len(data)
    s := make([]*big.Int, l)
    for i := 0; i < l; i += 1 {
        s[i] = big.NewInt(int64(data[i]))
    }
    return NewBigMatrix(rows, cols, s, nil)
}

// get value at (row, col), where first row/col is 0.
func (m BigMatrix) At(row, col int) *big.Int {
    if row >= m.rows || col >= m.cols || row < 0 || col < 0{
        panic(fmt.Sprintf("Index out of bounds: (%d, %d)", row, col))
    }
    valueIndex := m.cols*row + col
    return m.values[valueIndex]
}

// set value at (row, col), where first row/col is 0.
func (m BigMatrix) Set(row, col int, value *big.Int) {
    if row >= m.rows || col >= m.cols || row < 0 || col < 0 {
        panic(fmt.Sprintf("Index out of bounds: (%d, %d)", row, col))
    }
    m.values[m.cols*row + col] = value
}

// multiply a * b
// also handles multiplication of encrypted * unecrypted matrices or vice versa
// if a and b are encrypted under different cryptosystems, the cryptosystem of a is used
func (a BigMatrix) Multiply(b BigMatrix) (BigMatrix, error) {
    if a.cols != b.rows {
        panic("matrices a and b are not compatible")
    }
    cRows, cCols := a.rows, b.cols
    values := make([]*big.Int, cRows*cCols)
    var r *big.Int
    var err error
    for i := 0; i < cRows; i += 1 {
        for j := 0; j < cCols; j += 1 {
            var sum *big.Int
            for k := 0; k < a.cols; k += 1 {
                if reflect.TypeOf(a.cryptosystem) == reflect.TypeOf(plain{}) {
                    r, err = b.cryptosystem.MultiplyFactor(b.At(k, j), a.At(i, k))
                    if err != nil {return a, err}
                    if sum == nil {
                        sum = r
                    } else {
                        sum, err = b.cryptosystem.Add(r, sum)
                        if err != nil {return a, err}
                    }
                } else {
                    if reflect.TypeOf(b.cryptosystem) == reflect.TypeOf(plain{}) {
                        r, err = a.cryptosystem.MultiplyFactor(a.At(i, k), b.At(k, j))
                        if err != nil {return a, err}
                    } else {
                        r, err = a.cryptosystem.Multiply(a.At(i, k), b.At(k, j))
                        if err != nil {return a, err}
                    }
                    if sum == nil {
                        sum = r
                    } else {
                        sum, err = a.cryptosystem.Add(r, sum)
                        if err != nil {return a, err}
                    }
                }
            }
            values[i*cCols+j] = sum
            sum = nil
        }
    }
    return NewBigMatrix(cRows, cCols, values, a.cryptosystem), nil
}

// multiplication of a by a factor
// assumes matrix and factor is in same space, otherwise use MultiplyPlaintextFactor
func (a BigMatrix) MultiplyFactor(factor *big.Int) (BigMatrix, error) {
    two_val_mul := func (t1, t2 *big.Int) (*big.Int, error) {return a.cryptosystem.Multiply(t1, t2)}
    return scalarMultiplication(two_val_mul, a, factor)
}

// multiplication of a by a factor
// to be used if a is encrypted while factor is not
func (a BigMatrix) MultiplyPlaintextFactor(factor *big.Int) (BigMatrix, error) {
    return scalarMultiplication(a.cryptosystem.MultiplyFactor, a, factor)
}

func scalarMultiplication(mulfunc func(*big.Int, *big.Int) (*big.Int, error), a BigMatrix, b *big.Int) (BigMatrix, error) {
    c_vals := make([]*big.Int, len(a.values))
    var err error
    for i := range c_vals {
        c_vals[i], err = mulfunc(a.values[i], b)
        if err != nil {return a, err}
    }
    c := NewBigMatrix(a.rows, a.cols, c_vals, nil)
    return c, nil
}

// matrix addition
func (a BigMatrix) Add(b BigMatrix) (BigMatrix, error) {
    if a.rows != b.rows {
        panic("row mismatch in addition")
    } else if a.cols != b.cols {
        panic("column mismatch in addition")
    }
    c_vals := make([]*big.Int, len(a.values))
    var err error
    for i := range c_vals {
        c_vals[i], err = a.cryptosystem.Add(a.values[i], b.values[i])
        if err != nil {return a, err}
    }
    c := NewBigMatrix(a.rows, a.cols, c_vals, a.cryptosystem)
    return c, nil
}

// matrix subtraction
func (a BigMatrix) Subtract(b BigMatrix) (BigMatrix, error) {
    if a.rows != b.rows {
        panic("row mismatch in subtraction")
    } else if a.cols != b.cols {
        panic("column mismatch in subtraction")
    }
    c_vals := make([]*big.Int, len(a.values))
    neg := big.NewInt(-1)
    for i := range c_vals {
        b_neg, err := a.cryptosystem.MultiplyFactor(b.values[i], neg)
        if err != nil {return a, err}
        c_vals[i], err = a.cryptosystem.Add(a.values[i], b_neg)
        if err != nil {return a, err}
    }
    c := NewBigMatrix(a.rows, a.cols, c_vals, a.cryptosystem)
    return c, nil
}

// concatenate matrices as A|B
func (a BigMatrix) Concatenate(b BigMatrix) BigMatrix {
    if a.rows != b.rows {
        panic("matrices not compatible for concatenation")
    }
    vals := make([]*big.Int, 0, (a.cols + b.cols) * a.rows)
    for i := 0; i < a.rows; i += 1 {
        vals = append(vals, a.values[i*a.cols:(i+1)*a.cols]...)
        vals = append(vals, b.values[i*b.cols:(i+1)*b.cols]...)
    }
    return NewBigMatrix(a.rows, a.cols + b.cols, vals, a.cryptosystem)
}

// create a new matrix from last k columns of a
func (a BigMatrix) CropHorizontally(k int) BigMatrix {
    vals := make([]*big.Int, 0, k*a.rows)
    d := a.cols - k
    for i := 0; i < a.rows; i += 1 {
        vals = append(vals, a.values[i*a.cols+d:(i+1)*a.cols]...)
    }
    return NewBigMatrix(a.rows, k, vals, a.cryptosystem)
}

// apply Mod for all matrix elements
func (a BigMatrix) Mod(mod *big.Int) BigMatrix {
    b_vals := make([]*big.Int, len(a.values))
    for i := range a.values {
        b_vals[i] = new(big.Int).Mod(a.values[i], mod)
    }
    return NewBigMatrix(a.rows, a.cols, b_vals, a.cryptosystem)
}