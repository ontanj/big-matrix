package bigmatrix

import (
    "testing"
    "math/big"
)

func Compare(a, b BigMatrix, t *testing.T) {
    if a.cols != b.cols {
        t.Errorf("differing number of columns (%d and %d)", a.cols, b.cols)
    }
    if a.rows != b.rows {
        t.Errorf("differing number of columns (%d and %d)", a.rows, b.rows)
    }
    for i := 0; i < a.rows; i += 1 {
        for j := 0; j < a.cols; j += 1 {
            if a.At(i, j).(*big.Int).Cmp(b.At(i, j).(*big.Int)) != 0 {
                t.Errorf("values differ at (%d, %d)", i, j)
            }
        }
    }
}

func TestValidNewBigMatrix(t *testing.T) {
    t.Run("vanilla", func(t *testing.T){
        var matrixData []interface{}
        var testData []*big.Int
        var dval int
        for dval = 1; dval <= 9; dval++ {
            matrixData = append(matrixData, big.NewInt(int64(dval)))
            testData = append(testData, big.NewInt(int64(dval)))
        }
        m := NewBigMatrix(3, 3, matrixData, nil)
        if m.cols != 3 {
            t.Error("wrong column size")
        }
        if m.rows != 3 {
            t.Error("wrong row size")
        }
        for i := 0; i < 9; i++ {
            if m.values[i].(*big.Int).Cmp(testData[i]) != 0 {
                t.Error("wrong data")
            }
        }
    })
    t.Run("uninitialized data", func(t *testing.T) {
        m := NewBigMatrix(3, 3, nil, nil)
        for i := 0; i < 9; i++ {
            if m.values[i] != nil {
                t.Error("not initialized")
            }
        }
    })
}

func TestValidNewBigMatrixFromInt(t *testing.T) {
    t.Run("vanilla", func(t *testing.T){
        var matrixData []int
        var testData []*big.Int
        var dval int
        for dval = 1; dval <= 9; dval++ {
            matrixData = append(matrixData, dval)
            testData = append(testData, big.NewInt(int64(dval)))
        }
        m := NewBigMatrixFromInt(3, 3, matrixData)
        if m.cols != 3 {
            t.Error("wrong column size")
        }
        if m.rows != 3 {
            t.Error("wrong row size")
        }
        for i := 0; i < 9; i++ {
            if m.values[i].(*big.Int).Cmp(testData[i]) != 0 {
                t.Error("wrong data")
            }
        }
    })
    t.Run("uninitialized data", func(t *testing.T) {
        m := NewBigMatrixFromInt(3, 3, nil)
        for i := 0; i < 9; i++ {
            if m.values[i] != nil {
                t.Error("not initialized")
            }
        }
    })
}

func TestInvalidNewBigMatrix(t *testing.T) {
    defer func() {
        if recover() == nil {
            t.Error("contructor did not panic on mismatched size")
        }
    }()
    var matrixData []interface{}
    var dval int64
    for dval = 1; dval <= 8; dval++ {
        matrixData = append(matrixData, big.NewInt(dval))
    }
    NewBigMatrix(3, 3, matrixData, nil)
}

func TestAt(t *testing.T) {
    var matrixData []interface{}
    var testData []*big.Int
    var dval int64
    for dval = 1; dval <= 9; dval++ {
        matrixData = append(matrixData, big.NewInt(dval))
        testData = append(testData, big.NewInt(dval))
    }
    m := NewBigMatrix(3, 3, matrixData, nil)
    row, col := 0, 0
    for _, val := range testData {
        if val.Cmp(m.At(row, col).(*big.Int)) != 0 {
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
        defer func() {
            if recover() == nil {
                t.Error("didn't panic on index out of bounds")
            }
        }()
        m.At(0,3)
        m.At(3,0)
    })
}

func TestSet(t *testing.T) {
    a := NewBigMatrixFromInt(2, 2, []int{1,2,3,4})
    b := NewBigMatrixFromInt(2, 2, []int{1,2,5,4})
    a.Set(1,0,big.NewInt(5))
    Compare(a, b, t)
    defer func() {
        if recover() == nil {
            t.Error("didn't panic on index out of bounds")
        }
    }()
    a.Set(0,3,big.NewInt(10))
}

func TestMultiplication(t *testing.T) {
    a := NewBigMatrixFromInt(2, 2, []int{1,2,3,4})
    b := NewBigMatrixFromInt(2, 3, []int{1,2,3,4,5,6})
    t.Run("vanilla", func(t *testing.T) {
        c, err := a.Multiply(b)
        if err != nil {t.Error(err)}
        d := NewBigMatrixFromInt(2, 3, []int{9,12,15,19,26,33})
        Compare(c, d, t)
    })
    t.Run("dimension mismatch", func(t *testing.T) {
        defer func() {
            if recover() == nil {
                t.Error("multiplication of mismatched matrices passed")
            }
        }()
        b.Multiply(a)
    })
}

func TestAddition(t *testing.T) {
    a := NewBigMatrixFromInt(2, 2, []int{1,2,3,4})
    doubleA, err := a.Add(a)
    if err != nil {t.Error(err)}
    correct := NewBigMatrixFromInt(2, 2, []int{2,4,6,8})

    t.Run("vanilla addition", func(t *testing.T) {
        Compare(doubleA, correct, t)
    })
    t.Run("row mismatch", func(t *testing.T) {
        defer func() {
            if recover() == nil {
                t.Error("addition of mismatched matrices passed")
            }
        }()
        d := NewBigMatrix(3,2,nil,nil)
        a.Add(d)
    })
    t.Run("column mismatch", func(t *testing.T) {
        defer func() {
            if recover() == nil {
                t.Error("addition of mismatched matrices passed")
            }
        }()
        e := NewBigMatrix(2,3,nil,nil)
        a.Add(e)
    })
}

func TestSubtraction(t *testing.T) {
    a := NewBigMatrixFromInt(2, 2, []int{5,3,7,9})
    b := NewBigMatrixFromInt(2, 2, []int{1,2,3,4})
    c, err := a.Subtract(b)
    if err != nil {t.Error(err)}
    correct := NewBigMatrixFromInt(2, 2, []int{4,1,4,5})
    t.Run("vanilla subtraction", func(t *testing.T) {
        Compare(c, correct, t)
    })

    t.Run("row mismatch", func(t *testing.T) {
        defer func() {
            if recover() == nil {
                t.Error("addition of mismatched matrices passed")
            }
        }()
        e := NewBigMatrix(3,2,nil,nil)
        a.Subtract(e)
    })
    t.Run("column mismatch", func(t *testing.T) {
        defer func() {
            if recover() == nil {
                t.Error("addition of mismatched matrices passed")
            }
            }()
        f := NewBigMatrix(2,3,nil,nil)	
        a.Subtract(f)
    })
}

func TestFactorMultiplication(t *testing.T) {
    a := NewBigMatrixFromInt(2, 3, []int{3, 4, 2, 1, 8, 5})
    b := big.NewInt(2)
    c := NewBigMatrixFromInt(2, 3, []int{6, 8, 4, 2, 16, 10})
    d, err := a.MultiplyFactor(b)
    if err != nil {t.Error(err)}
    Compare(c, d, t)
}

func TestConcatenation(t *testing.T) {
    a := NewBigMatrixFromInt(3, 2, []int{1, 2, 3, 4, 5, 6})
    t.Run("valid concatenation", func(t *testing.T) {
        b := NewBigMatrixFromInt(3, 2, []int{1, 2, 3, 4, 5, 6})
        correct := NewBigMatrixFromInt(3, 4, []int{1, 2, 1, 2, 3, 4, 3, 4, 5, 6, 5, 6})
        ab := a.Concatenate(b)
        Compare(correct, ab, t)
    })
    t.Run("invalid concatenation", func(t *testing.T) {
        b := NewBigMatrix(4, 4, nil, nil)
        defer func() {
            if recover() == nil {
                t.Error("invalid concatenation did not panic")
            }
        }()
        a.Concatenate(b)
    })
}

func TestCrop(t *testing.T) {
    a := NewBigMatrixFromInt(3, 3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
    a = a.CropHorizontally(2)   
    correct := NewBigMatrixFromInt(3, 2, []int{2, 3, 5, 6, 8, 9})
    Compare(a, correct, t)
}

func TestMod(t *testing.T) {
    a := NewBigMatrixFromInt(3, 2, []int{9,4,6,3,8,6})
    b := NewBigMatrixFromInt(3, 2, []int{0,1,0,0,2,0})
    a, err := a.Apply(func(val interface{}) (interface{}, error) {return new(big.Int).Mod(val.(*big.Int), big.NewInt(3)), nil})
    if err != nil {t.Error(err)}
    for i := 0; i < 3; i += 1 {
        for j := 0; j < 2; j += 1 {
            if a.At(i,j).(*big.Int).Cmp(b.At(i,j).(*big.Int)) != 0 {
                t.Errorf("(%d,%d): expected %d, got %d", i, j, a.At(i,j), b.At(i,j))
            }
        }
    }
}
