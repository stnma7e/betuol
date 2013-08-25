package math

type Mat2x2 [4]float32

func (mat *Mat2x2) Mult(scalar float32) {
	for i := range mat {
		mat[i] = scalar * mat[i]
	}
}
func (mat *Mat2x2) Invert() {
	frac := 1 / (mat[0]*mat[3] - mat[1]*mat[2])
	inv := Mat2x2{mat[3],-mat[1],-mat[2],mat[3]}
	inv.Mult(frac)
	*mat = inv
}
func (segs *Mat2x2) Split() (a,b,c,d float32) {
	return segs[0], segs[1], segs[2], segs[3]
}