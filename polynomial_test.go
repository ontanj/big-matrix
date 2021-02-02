package genmatrix

import (
    "testing"
    "math/big"
)

func ComparePoly(a, b Polynomial, t *testing.T) {
    if a.Size != b.Size {
        t.Errorf("different sizes (%d and %d)", a.Size, b.Size)
    }
    for i := 0; i < a.Size; i += 1 {
        a_val, _ := decode(a.At(i))
        b_val, _ := decode(b.At(i))
        if a_val.Cmp(b_val) != 0 {
            t.Errorf("values differ at %d", i)
        }
    }
}

func TestPolyMultiplication(t *testing.T) {
    a, err := NewPolyFromInt([]int{1,2,3,4})
    if err != nil {t.Error(err)}
    b, err := NewPolyFromInt([]int{1,2,3})
    if err != nil {t.Error(err)}
    t.Run("a bigger", func(t *testing.T) {
        c, err := a.Multiply(b)
        if err != nil {t.Error(err)}
        d, err := NewPolyFromInt([]int{1,4,10,16,17,12})
        if err != nil {t.Error(err)}
        ComparePoly(c, d, t)
    })
    t.Run("b bigger", func(t *testing.T) {
        c, err := b.Multiply(a)
        if err != nil {t.Error(err)}
        d, err := NewPolyFromInt([]int{1,4,10,16,17,12})
        if err != nil {t.Error(err)}
        ComparePoly(c, d, t)
    })
}

func TestPolyAddition(t *testing.T) {
    a, err := NewPolyFromInt([]int{1,2,3,4})
    if err != nil {t.Error(err)}
    doubleA, err := a.Add(a)
    if err != nil {t.Error(err)}
    correct, err := NewPolyFromInt([]int{2,4,6,8})
    if err != nil {t.Error(err)}
    ComparePoly(doubleA, correct, t)
}

func TestPolySubtraction(t *testing.T) {
    a, err := NewPolyFromInt([]int{5,3,7,9})
    if err != nil {t.Error(err)}
    b, err := NewPolyFromInt([]int{1,2,3})
    if err != nil {t.Error(err)}
    t.Run("a bigger", func(t *testing.T) {
        c, err := a.Subtract(b)
        if err != nil {t.Error(err)}
        correct, err := NewPolyFromInt([]int{4,1,4,9})
        if err != nil {t.Error(err)}
        ComparePoly(c, correct, t)
    })
    t.Run("b bigger", func(t *testing.T) {
        c, err := b.Subtract(a)
        if err != nil {t.Error(err)}
        correct, err := NewPolyFromInt([]int{-4,-1,-4,-9})
        if err != nil {t.Error(err)}
        ComparePoly(c, correct, t)
    })
}

func TestPolyFactorMultiplication(t *testing.T) {
    a, err := NewPolyFromInt([]int{3, 4, 2, 1, 8, 5})
    if err != nil {t.Error(err)}
    b := big.NewInt(2)
    c, err := NewPolyFromInt([]int{6, 8, 4, 2, 16, 10})
    if err != nil {t.Error(err)}
    d, err := a.MultiplyScalar(b)
    if err != nil {t.Error(err)}
    ComparePoly(c, d, t)
}

func TestPolyMod(t *testing.T) {
    a, err := NewPolyFromInt([]int{9,4,6,3,8,6})
    if err != nil {t.Error(err)}
    b, err := NewPolyFromInt([]int{0,1,0,0,2,0})
    if err != nil {t.Error(err)}
    a, err = a.Apply(func(val interface{}) (interface{}, error) {return new(big.Int).Mod(val.(*big.Int), big.NewInt(3)), nil})
    if err != nil {t.Error(err)}
    for i := 0; i < 3; i += 1 {
        a_int, err := decode(a.At(i))
        if err != nil {t.Error(err)}
        b_int, err := decode(b.At(i))
        if err != nil {t.Error(err)}
        if a_int.Cmp(b_int) != 0 {
            t.Errorf("%d: expected %d, got %d", i, a_int, b_int)
        }
    }
}
