package genmatrix

import (
    "testing"
    "math/big"
)

func decode(val interface{}, err error) (*big.Int, error) {
    return val.(*big.Int), err
}

func Compare(a, b Matrix, t *testing.T) {
    if a.Cols != b.Cols {
        t.Errorf("differing number of columns (%d and %d)", a.Cols, b.Cols)
    }
    if a.Rows != b.Rows {
        t.Errorf("differing number of columns (%d and %d)", a.Rows, b.Rows)
    }
    for i := 0; i < a.Rows; i += 1 {
        for j := 0; j < a.Cols; j += 1 {
            a_val, _ := decode(a.At(i, j))
            b_val, _ := decode(b.At(i, j))
            if a_val.Cmp(b_val) != 0 {
                t.Errorf("values differ at (%d, %d)", i, j)
            }
        }
    }
}

func TestValidNewMatrix(t *testing.T) {
    t.Run("vanilla", func(t *testing.T){
        var matrixData []interface{}
        var testData []*big.Int
        var dval int
        for dval = 1; dval <= 9; dval++ {
            matrixData = append(matrixData, big.NewInt(int64(dval)))
            testData = append(testData, big.NewInt(int64(dval)))
        }
        m, err := NewMatrix(3, 3, matrixData, bigint{})
        if err != nil {t.Error(err)}
        if m.Cols != 3 {
            t.Error("wrong column size")
        }
        if m.Rows != 3 {
            t.Error("wrong row size")
        }
        for row := 0; row < 3; row++ {
            for col := 0; col < 3; col++ {
                m_val, err := decode(m.At(row, col))
                if err != nil {t.Error(err)}
                exp := testData[3*row + col]
                if m_val.Cmp(exp) != 0 {
                    t.Errorf("wrong data: expected %d got %d", exp, m_val)
                }
            }
        }
    })
    t.Run("uninitialized data", func(t *testing.T) {
        m, err := NewMatrix(3, 3, nil, bigint{})
        if err != nil {t.Error(err)}
        for i := 0; i < 9; i++ {
            if m.values[i] != nil {
                t.Error("not initialized")
            }
        }
    })
}

func TestValidNewMatrixFromInt(t *testing.T) {
    t.Run("vanilla", func(t *testing.T){
        var matrixData []int
        var testData []*big.Int
        var dval int
        for dval = 1; dval <= 9; dval++ {
            matrixData = append(matrixData, dval)
            testData = append(testData, big.NewInt(int64(dval)))
        }
        m, err := NewMatrixFromInt(3, 3, matrixData)
        if err != nil {t.Error(err)}
        if m.Cols != 3 {
            t.Error("wrong column size")
        }
        if m.Rows != 3 {
            t.Error("wrong row size")
        }
        for i := 0; i < 9; i++ {
            if m.values[i].(*big.Int).Cmp(testData[i]) != 0 {
                t.Error("wrong data")
            }
        }
    })
    t.Run("uninitialized data", func(t *testing.T) {
        m, err := NewMatrixFromInt(3, 3, nil)
        if err != nil {t.Error(err)}
        for i := 0; i < 9; i++ {
            if m.values[i] != nil {
                t.Error("not initialized")
            }
        }
    })
}

func TestInvalidNewMatrix(t *testing.T) {
    var matrixData []interface{}
    var dval int64
    for dval = 1; dval <= 8; dval++ {
        matrixData = append(matrixData, big.NewInt(dval))
    }
    _, err := NewMatrix(3, 3, matrixData, bigint{})
    if err == nil {t.Error("no error on mismatched size")}
}

func TestAt(t *testing.T) {
    var matrixData []interface{}
    var testData []*big.Int
    var dval int64
    for dval = 1; dval <= 9; dval++ {
        matrixData = append(matrixData, big.NewInt(dval))
        testData = append(testData, big.NewInt(dval))
    }
    m, err := NewMatrix(3, 3, matrixData, bigint{})
    if err != nil {t.Error(err)}
    row, col := 0, 0
    for _, val := range testData {
        int_val, err := decode(m.At(row, col))
        if err != nil {t.Error(err)}
        if val.Cmp(int_val) != 0 {
            t.Error("malformed matrix")
        }
        if col == 2 {
            col = 0
            row += 1
        } else {
            col += 1
        }
    }
    t.Run("index out of bounds", func (t *testing.T) {
        _, err := m.At(0,3)
        if err == nil {t.Error("no error on index out of bounds")}
        _, err = m.At(3,0)
        if err == nil {t.Error("no error on index out of bounds")}
    })
}

func TestSet(t *testing.T) {
    a, err := NewMatrixFromInt(2, 2, []int{1,2,3,4})
    if err != nil {t.Error(err)}
    b, err := NewMatrixFromInt(2, 2, []int{1,2,5,4})
    if err != nil {t.Error(err)}
    a.Set(1,0,big.NewInt(5))
    Compare(a, b, t)
    err = a.Set(0,3,big.NewInt(10))
    if err == nil {t.Error("no error on index out of bounds")}
}

func TestMultiplication(t *testing.T) {
    a, err := NewMatrixFromInt(2, 2, []int{1,2,3,4})
    if err != nil {t.Error(err)}
    b, err := NewMatrixFromInt(2, 3, []int{1,2,3,4,5,6})
    if err != nil {t.Error(err)}
    t.Run("vanilla", func(t *testing.T) {
        c, err := a.Multiply(b)
        if err != nil {t.Error(err)}
        d, err := NewMatrixFromInt(2, 3, []int{9,12,15,19,26,33})
        if err != nil {t.Error(err)}
        Compare(c, d, t)
    })
    t.Run("dimension mismatch", func(t *testing.T) {
        _, err := b.Multiply(a)
        if err == nil {t.Error("no error on dimension mismatch")}
    })
}

func TestAddition(t *testing.T) {
    a, err := NewMatrixFromInt(2, 2, []int{1,2,3,4})
    if err != nil {t.Error(err)}
    doubleA, err := a.Add(a)
    if err != nil {t.Error(err)}
    correct, err := NewMatrixFromInt(2, 2, []int{2,4,6,8})
    if err != nil {t.Error(err)}

    t.Run("vanilla addition", func(t *testing.T) {
        Compare(doubleA, correct, t)
    })
    t.Run("row mismatch", func(t *testing.T) {
        d, err := NewMatrix(3, 2, nil, bigint{})
        if err != nil {t.Error(err)}
        _, err = a.Add(d)
        if err == nil {t.Error("addition of mismatched matrices passed")}
    })
    t.Run("column mismatch", func(t *testing.T) {
        e, err := NewMatrix(2, 3, nil, bigint{})
        if err != nil {t.Error(err)}
        _, err = a.Add(e)
        if err == nil {t.Error("addition of mismatched matrices passed")}
    })
}

func TestSubtraction(t *testing.T) {
    a, err := NewMatrixFromInt(2, 2, []int{5,3,7,9})
    if err != nil {t.Error(err)}
    b, err := NewMatrixFromInt(2, 2, []int{1,2,3,4})
    if err != nil {t.Error(err)}
    c, err := a.Subtract(b)
    if err != nil {t.Error(err)}
    correct, err := NewMatrixFromInt(2, 2, []int{4,1,4,5})
    if err != nil {t.Error(err)}
    t.Run("vanilla subtraction", func(t *testing.T) {
        Compare(c, correct, t)
    })

    t.Run("row mismatch", func(t *testing.T) {
        e, err := NewMatrix(3, 2, nil, bigint{})
        if err != nil {t.Error(err)}
        _, err = a.Subtract(e)
        if err == nil {t.Error("subtraction of mismatched matrices passed")}
    })
    t.Run("column mismatch", func(t *testing.T) {
        f, err := NewMatrix(2, 3, nil, bigint{})	
        if err != nil {t.Error(err)}
        _, err = a.Subtract(f)
        if err == nil {t.Error("subtraction of mismatched matrices passed")}
    })
}

func TestFactorMultiplication(t *testing.T) {
    a, err := NewMatrixFromInt(2, 3, []int{3, 4, 2, 1, 8, 5})
    if err != nil {t.Error(err)}
    b := big.NewInt(2)
    c, err := NewMatrixFromInt(2, 3, []int{6, 8, 4, 2, 16, 10})
    if err != nil {t.Error(err)}
    d, err := a.MultiplyScalar(b)
    if err != nil {t.Error(err)}
    Compare(c, d, t)
}

func TestConcatenation(t *testing.T) {
    a, err := NewMatrixFromInt(3, 2, []int{1, 2, 3, 4, 5, 6})
    if err != nil {t.Error(err)}
    t.Run("valid concatenation", func(t *testing.T) {
        b, err := NewMatrixFromInt(3, 2, []int{1, 2, 3, 4, 5, 6})
        if err != nil {t.Error(err)}
        correct, err := NewMatrixFromInt(3, 4, []int{1, 2, 1, 2, 3, 4, 3, 4, 5, 6, 5, 6})
        if err != nil {t.Error(err)}
        ab, err := a.Concatenate(b)
        if err != nil {t.Error(err)}
        Compare(correct, ab, t)
    })
    t.Run("invalid concatenation", func(t *testing.T) {
        b, err := NewMatrix(4, 4, nil, bigint{})
        if err != nil {t.Error(err)}
        _, err = a.Concatenate(b)
        if err == nil {t.Error("no error on mismatched dimensions")}
    })
}

func TestCrop(t *testing.T) {
    a, err := NewMatrixFromInt(3, 3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
    if err != nil {t.Error(err)}
    a = a.CropHorizontally(2)   
    correct, err := NewMatrixFromInt(3, 2, []int{2, 3, 5, 6, 8, 9})
    if err != nil {t.Error(err)}
    Compare(a, correct, t)
}

func TestMod(t *testing.T) {
    a, err := NewMatrixFromInt(3, 2, []int{9,4,6,3,8,6})
    if err != nil {t.Error(err)}
    b, err := NewMatrixFromInt(3, 2, []int{0,1,0,0,2,0})
    if err != nil {t.Error(err)}
    a, err = a.Apply(func(val interface{}) (interface{}, error) {return new(big.Int).Mod(val.(*big.Int), big.NewInt(3)), nil})
    if err != nil {t.Error(err)}
    for i := 0; i < 3; i += 1 {
        for j := 0; j < 2; j += 1 {
            a_int, err := decode(a.At(i,j))
            if err != nil {t.Error(err)}
            b_int, err := decode(b.At(i,j))
            if err != nil {t.Error(err)}
            if a_int.Cmp(b_int) != 0 {
                t.Errorf("(%d,%d): expected %d, got %d", i, j, a_int, b_int)
            }
        }
    }
}
