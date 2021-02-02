package genmatrix

import (
    "fmt"
)

type Polynomial struct {
    values []interface{}
    Size int
    Space Space
}

// create a new Polynomial with data[0] being the constant term
func NewPolynomial(data []interface{}, space Space) (m Polynomial, err error) {
    if data == nil {
        err = fmt.Errorf("data can't be nil")
        return
    }
    if space == nil {
        err = fmt.Errorf("space can't be nil")
        return
    }
    m.values = data
    m.Size = len(data)
    m.Space = space
    return
}

// get coefficient at index, where index 0 is constant term
func (m Polynomial) At(index int) (interface{}, error) {
    if index >= m.Size {
        return nil, fmt.Errorf("Index out of bounds: %d", index)
    }
    return m.values[index], nil
}

// set value at (row, col), where first row/col is 0.
func (m Polynomial) Set(index int, value interface{}) error {
    if index != m.Size {
        return fmt.Errorf("Index out of bounds: %d", index)
    }
    m.values[index] = value
    return nil
}

// multiply a * b
// also handles multiplication of scalar * non-scalar polynomial and vice versa
// if a and b are non-scalar in different spaces, the space of a is used
func (a Polynomial) Multiply(b Polynomial) (c Polynomial, err error) { //TODO
    l := b.Size + a.Size - 1
    values := make([]interface{}, l)
    for i := 0; i < l; i += 1 {
        if b.Space.Scalarspace() {
            values[i], err = a.Space.Scale(a.values[i], b.values[0])
            if err != nil {return}
        } else {
            values[i], err = a.Space.Multiply(a.values[i], b.values[0])
            if err != nil {return}
        }
        for j := 1; j < i; j += 1 {
            var contrib interface{}
            if b.Space.Scalarspace() {
                contrib, err = a.Space.Scale(a.values[i], b.values[0])
                if err != nil {return}
            } else {
                contrib, err = a.Space.Multiply(a.values[i], b.values[0])
                if err != nil {return}
            }
            values[i], err = a.Space.Add(values[i], contrib)
            if err != nil {return}
        }
    }
    return NewPolynomial(values, a.Space)
}

// multiplication of a by a scalar
// assumes polynomial and factor is in same space, otherwise use Scale
func (a Polynomial) MultiplyScalar(scalar interface{}) (Polynomial, error) {
    return polyScalarMultiplication(a.Space.Multiply, a, scalar)
}

// scale a according to scalar
// to be used if factor is in a scalar space wile a is not
func (a Polynomial) Scale(factor interface{}) (Polynomial, error) {
    return polyScalarMultiplication(a.Space.Scale, a, factor)
}

func polyScalarMultiplication(mulfunc func(interface{}, interface{}) (interface{}, error), a Polynomial, b interface{}) (Polynomial, error) {
    c_vals := make([]interface{}, a.Size)
    var err error
    for i := range c_vals {
        c_vals[i], err = mulfunc(a.values[i], b)
        if err != nil {return a, err}
    }
    return NewPolynomial(c_vals, a.Space)
}

// polynomial addition
func (a Polynomial) Add(b Polynomial) (c Polynomial, err error) {
    var lng Polynomial
    var sml Polynomial
    if a.Size >= b.Size {
        lng = a
        sml = b
    } else {
        lng = b
        sml = a
    }
    c_vals := make([]interface{}, lng.Size)
    i := 0
    for ; i < sml.Size; i += 1 {
        c_vals[i], err = a.Space.Add(lng.values[i], sml.values[i])
        if err != nil {return a, err}
    }
    for ; i < lng.Size; i += 1 {
        c_vals[i] = lng.values[i]
    }
    return NewPolynomial(c_vals, a.Space)
}

// matrix subtraction
func (a Polynomial) Subtract(b Polynomial) (c Polynomial, err error) {
    var l int
    if b.Size >= a.Size {
        l = b.Size
    } else {
        l = a.Size
    }
    c_vals := make([]interface{}, l)
    for i := range c_vals {
        if i < a.Size && i < b.Size {
            c_vals[i], err = a.Space.Subtract(a.values[i], b.values[i])
            if err != nil {return a, err}
        } else if i < a.Size {
            c_vals[i] = a.values[i]
        } else {
            c_vals[i], err = b.Space.Negate(b.values[i])
            if err != nil {return a, err}
        }
    }
    return NewPolynomial(c_vals, a.Space)
}

// evaluate polynomial at point x
// if x belongs to a different space, use EvaluateInSpace
func (a Polynomial) Evaluate(x interface{}) (eval interface{}, err error) {
    return a.EvaluateInSpace(x, a.Space)
}

// evaluate polynomial at point x
// powers of x is taken in custom space
func (a Polynomial) EvaluateInSpace(x interface{}, space Space) (eval interface{}, err error) {
    eval = a.values[0]
    xn := x
    for i := 1; i < a.Size; i += 1 {
        eval, err = a.Space.Multiply(a.values[i], xn)
        if err != nil {return}
        xn, err = space.Multiply(xn, x)
        if err != nil {return}
    }
    return
}

// apply function f to all matrix elements
func (a Polynomial) Apply(f func(interface{}) (interface{}, error)) (b Polynomial, err error) {
    b_vals := make([]interface{}, len(a.values))
    for i, v := range a.values {
        b_vals[i], err = f(v)
        if err != nil {return}
    }
    return NewPolynomial(b_vals, a.Space)
}