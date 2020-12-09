package bigmatrix

type space interface {

    // addition of two elements in the space
    Add(interface{}, interface{}) (sum interface{}, err error)

    // subtraction of two elements in the space
    Subtract(interface{}, interface{}) (diff interface{}, err error)

    // multiplication of two elements in the space
    Multiply(interface{}, interface{}) (product interface{}, err error)

    // scaling of an element by scalar factor
    Scale(spaced interface{}, factor interface{}) (product interface{}, err error)

    // return true if this space (matrix) consist of scalar factors
    // e.i. if Scale are to be used in matrix multiplication
    Scalarspace() bool
}
