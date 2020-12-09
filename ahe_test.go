package genmatrix

import (
    "math/big"
    "testing"
    "github.com/niclabs/tcpaillier"
)

func EncryptMatrix(a Matrix, pk *tcpaillier.PubKey) (b Matrix, err error) {
    b_vals := make([]interface{}, len(a.values))
    for i := range a.values {
        b_vals[i], _, err = pk.Encrypt(a.values[i].(*big.Int))
        if err != nil {return}
    }
    return NewMatrix(a.Rows, a.Cols, b_vals, DJ_public_key{pk})
}

func DecryptMatrix(cipher Matrix, pk *tcpaillier.PubKey, sks []*tcpaillier.KeyShare) (plain Matrix, err error) {
    plain_vals := make([]interface{}, len(cipher.values))
    for i, enc_val := range cipher.values {
        part_dec := make([]*tcpaillier.DecryptionShare, len(sks))
        for j, sk := range sks {
            part_dec[j], err = sk.PartialDecrypt(enc_val.(*big.Int))
            if err != nil {return}
        }
        plain_vals[i], err = pk.CombineShares(part_dec...)
        if err != nil {return}
    }
    return NewMatrix(cipher.Rows, cipher.Cols, plain_vals, Bigint{})
}

func TestEncryptedMatrixAddition(t *testing.T) {
    cs, djsks, err := NewDJCryptosystem()
    if err != nil {t.Error(err)}
    a, err := NewMatrixFromInt(2, 3, []int{3, 4, 2, 1, 8, 5})
    if err != nil {t.Error(err)}
    b, err := NewMatrixFromInt(2, 3, []int{1, 2, 3, 4, 5, 6})
    if err != nil {t.Error(err)}
    c, err := NewMatrixFromInt(2, 3, []int{4, 6, 5, 5, 13, 11})
    if err != nil {t.Error(err)}
    a, err = EncryptMatrix(a, cs.PubKey)
    if err != nil {t.Error(err)}
    b, err = EncryptMatrix(b, cs.PubKey)
    if err != nil {t.Error(err)}
    sum, err := a.Add(b)
    if err != nil {t.Error(err)}
    sum, err = DecryptMatrix(sum, cs.PubKey, djsks)
    if err != nil {t.Error(err)}
    Compare(sum, c, t)
}

func TestEncryptedMatrixSubtraction(t *testing.T) {
    cs, djsks, err := NewDJCryptosystem()
    if err != nil {t.Error(err)}
    a, err := NewMatrixFromInt(2, 3, []int{3, 4, 2, 1, 8, 5})
    if err != nil {t.Error(err)}
    b, err := NewMatrixFromInt(2, 3, []int{1, 2, 2, 0, 4, 3})
    if err != nil {t.Error(err)}
    c, err := NewMatrixFromInt(2, 3, []int{2, 2, 0, 1, 4, 2})
    if err != nil {t.Error(err)}
    a, err = EncryptMatrix(a, cs.PubKey)
    if err != nil {t.Error(err)}
    b, err = EncryptMatrix(b, cs.PubKey)
    if err != nil {t.Error(err)}
    diff, err := a.Subtract(b)
    if err != nil {t.Error(err)}
    diff, err = DecryptMatrix(diff, cs.PubKey, djsks)
    if err != nil {t.Error(err)}
    Compare(diff, c, t)
}

func TestEncryptedMatrixMultiplication(t *testing.T) {
    a, err := NewMatrixFromInt(2, 3, []int{1,2,3,4,5,6})
    if err != nil {t.Error(err)}
    b, err := NewMatrixFromInt(3, 2, []int{1,2,3,4,5,6})
    if err != nil {t.Error(err)}
    cs, djsks, err := NewDJCryptosystem()
    if err != nil {t.Error(err)}
    ae, err := EncryptMatrix(a, cs.PubKey)
    if err != nil {t.Error(err)}
    t.Run("plaintext from right", func(t *testing.T) {  
        ab, err := ae.Multiply(b)
        if err != nil {t.Error(err)}
        correct, err := a.Multiply(b)
        if err != nil {t.Error(err)}
        ab, err = DecryptMatrix(ab, cs.PubKey, djsks)
        if err != nil {t.Error(err)}
        Compare(ab, correct, t)
    })
    t.Run("plaintext from left", func(t *testing.T) {  
        ba, err := b.Multiply(ae)
        if err != nil {t.Error(err)}
        correct, err := b.Multiply(a)
        if err != nil {t.Error(err)}
        ba, err = DecryptMatrix(ba, cs.PubKey, djsks)
        if err != nil {t.Error(err)}
        Compare(ba, correct, t)
    })
}

func TestMultiplyPlaintextFactor(t *testing.T) {
    a, err := NewMatrixFromInt(2, 3, []int{1,2,3,4,5,6})
    correct, err := NewMatrixFromInt(2, 3, []int{3,6,9,12,15,18})
    c := big.NewInt(3)
    cs, djsks, err := NewDJCryptosystem()
    if err != nil {t.Error(err)}
    a, err = EncryptMatrix(a, cs.PubKey)
    a, err = a.Scale(c)
    if err != nil {t.Error(err)}
    a, err = DecryptMatrix(a, cs.PubKey, djsks)
    Compare(a, correct, t)
}