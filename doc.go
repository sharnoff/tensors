// Tensors is a package designed for facilitating the use of many-dimensional mathematical tensors.
// Tensors is designed for its primary use case in github.com/sharnoff/badstudent, but is open to
// improvement.
//
// The main unit of this package is the Interpreter, which will suffice for most use cases.
// Additionally, many methods have two other, 'Safe' and 'Fast' variants, which offer increased or
// decreased error checking respectively.
//
// Much of this is quite self-explanatory -- the only relevant information is that, even though the
// package is named tensors, the fundamental unit is the Interpreter. Interpreters do the job of
// handling interactions with tensors; they 'interpret' individual and sets of indices to convert
// between them.
package tensors
