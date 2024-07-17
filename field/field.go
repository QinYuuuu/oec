package field

type Field interface {
	Add(a, b []byte) []byte
	Subtract(a, b []byte) []byte
	Multiply(a, b []byte) []byte
	Divide(a, b []byte) []byte
	Inverse(a []byte) []byte
	Zero() []byte
	One() []byte
}
