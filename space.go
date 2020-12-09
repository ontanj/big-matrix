package bigmatrix

type space interface {
    Add(interface{}, interface{}) (sum interface{}, err error)
    Subtract(interface{}, interface{}) (diff interface{}, err error)
    Multiply(interface{}, interface{}) (product interface{}, err error)
    MultiplyScalar(spaced interface{}, plaintext interface{}) (product interface{}, err error)
    IsPlaintext() bool
}
