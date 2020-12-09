# Generic Matrix

This library provides basic matrix operations for matrices with elements from any space.

## Space

The library is built around the `interface space`. It defines the element-wise operations needed for the matrix operations to work. Two examples are implemented in `bigint.go` and `damgard-jurik.go` where the operations are defined for `*big.Int` from the standard library, and the [additive homomorphic cryptosystem](https://www.researchgate.net/publication/225753264_A_generalization_of_Paillier%27s_public-key_system_with_applications_to_electronic_voting) described by Damgård and Jurik and implemented in [tcpaillier](https://github.com/niclabs/tcpaillier).

## Matrix structure

The matrices are defined as
```go
type Matrix struct {
    values []interface{}
    rows, cols int
    space space
}
```
where values are stored in row-major order and `space` stores the evaluation space for the matrix.
