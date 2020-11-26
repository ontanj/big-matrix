package bigmatrix

import (
    "math/big"
    "testing"
    "github.com/niclabs/tcpaillier"
)

func EncryptMatrix(a BigMatrix, pk *tcpaillier.PubKey) (b BigMatrix) {
    b_vals := make([]*big.Int, len(a.values))
    var err error
    for i := range a.values {
        b_vals[i], _, err = pk.Encrypt(a.values[i])
        if err != nil {panic(err)}
    }
    return NewBigMatrix(a.rows, a.cols, b_vals, dj_public_key{pk})
}

func DecryptMatrix(cipher BigMatrix, pk *tcpaillier.PubKey, sks []*tcpaillier.KeyShare) (plain BigMatrix) {
    plain_vals := make([]*big.Int, len(cipher.values))
    var err error
    for i, enc_val := range cipher.values {
        part_dec := make([]*tcpaillier.DecryptionShare, len(sks))
        for j, sk := range sks {
            part_dec[j], err = sk.PartialDecrypt(enc_val)
            if err != nil {panic(err)}
        }
        plain_vals[i], err = pk.CombineShares(part_dec...)
        if err != nil {panic(err)}
    }
    return NewBigMatrix(cipher.rows, cipher.cols, plain_vals, nil)
}

func TestEncryptedMatrixAddition(t *testing.T) {
    cs, djsks, _ := NewDJCryptosystem()
    a := NewBigMatrixFromInt(2, 3, []int{3, 4, 2, 1, 8, 5}, nil)
    b := NewBigMatrixFromInt(2, 3, []int{1, 2, 3, 4, 5, 6}, nil)
    c := NewBigMatrixFromInt(2, 3, []int{4, 6, 5, 5, 13, 11}, nil)
    a = EncryptMatrix(a, cs.PubKey)
    b = EncryptMatrix(b, cs.PubKey)
    sum, err := a.Add(b)
    if err != nil {t.Error(err)}
    sum = DecryptMatrix(sum, cs.PubKey, djsks)
    Compare(sum, c, t)
}

func TestEncryptedMatrixSubtraction(t *testing.T) {
    cs, djsks, _ := NewDJCryptosystem()
    a := NewBigMatrixFromInt(2, 3, []int{3, 4, 2, 1, 8, 5}, nil)
    b := NewBigMatrixFromInt(2, 3, []int{1, 2, 2, 0, 4, 3}, nil)
    c := NewBigMatrixFromInt(2, 3, []int{2, 2, 0, 1, 4, 2}, nil)
    a = EncryptMatrix(a, cs.PubKey)
    b = EncryptMatrix(b, cs.PubKey)
    diff, err := a.Subtract(b)
    if err != nil {t.Error(err)}
    diff = DecryptMatrix(diff, cs.PubKey, djsks)
    Compare(diff, c, t)
}

func TestEncryptedMatrixMultiplication(t *testing.T) {
    a := NewBigMatrixFromInt(2, 3, []int{1,2,3,4,5,6}, nil)
    b := NewBigMatrixFromInt(3, 2, []int{1,2,3,4,5,6}, nil)
    cs, djsks, _ := NewDJCryptosystem()
    ae := EncryptMatrix(a, cs.PubKey)
    t.Run("plaintext from right", func(t *testing.T) {  
        ab, err := ae.Multiply(b)
        if err != nil {t.Error(err)}
        correct, err := a.Multiply(b)
        if err != nil {t.Error(err)}
        ab = DecryptMatrix(ab, cs.PubKey, djsks)
        Compare(ab, correct, t)
    })
    t.Run("plaintext from left", func(t *testing.T) {  
        ba, err := b.Multiply(ae)
        if err != nil {t.Error(err)}
        correct, err := b.Multiply(a)
        if err != nil {t.Error(err)}
        ba = DecryptMatrix(ba, cs.PubKey, djsks)
        Compare(ba, correct, t)
    })
}

func TestMultiplyPlaintextFactor(t *testing.T) {
    a := NewBigMatrixFromInt(2, 3, []int{1,2,3,4,5,6}, nil)
    correct := NewBigMatrixFromInt(2, 3, []int{3,6,9,12,15,18}, nil)
    c := big.NewInt(3)
    cs, djsks, _ := NewDJCryptosystem()
    a = EncryptMatrix(a, cs.PubKey)
    a, err := a.MultiplyPlaintextFactor(c)
    if err != nil {t.Error(err)}
    a = DecryptMatrix(a, cs.PubKey, djsks)
    Compare(a, correct, t)
}