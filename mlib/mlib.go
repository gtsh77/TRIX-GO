package mlib

import (
	"fmt"
	"math"
	M "github.com/dtromb/gogsl/matrix"
	B "github.com/dtromb/gogsl/blas"
)

// ===== INITIALIZATION =====

//new empty matrix MxN
func NewMatrix(m int, n int) *M.GslMatrix {
	var ma *M.GslMatrix = M.MatrixAlloc(m, n)

	M.SetZero(ma)

	return ma
}

//new diag matrix MxM
func NewDiagMatrix(m int) *M.GslMatrix {
	var ma *M.GslMatrix

	M.SetZero(ma)
	for i := 0; i < m; i++ {
		M.Set(ma, i, i, 1.0)
	}

	return ma
}

//reset to diag matrix
func ResetDiagMatrix(m *M.GslMatrix, size int) {	
	M.SetZero(m)
	for i := 0; i < size; i++ {
		M.Set(m, i, i, 1.0)
	}
}

//set matrix elements to X
func SetAll(m *M.GslMatrix, x float64) {
	M.SetAll(m,x)
}

// ===== GLM LIKE FUNCS =====

//glm.perspective
func GlmPerspective(fovy float64, aspect float64, zNear float64, zFar float64, R *M.GslMatrix) {
	var tanHalfFovy float64 = math.Tan(fovy/2.0)

	M.Set(R,0,0,1.0/(aspect * tanHalfFovy))
	M.Set(R,1,1,1.0/(tanHalfFovy))
	M.Set(R,2,2, - 1.0)
	M.Set(R,2,2, - (zFar + zNear) / (zFar - zNear))
	M.Set(R,3,2, - (2.0 * zFar * zNear) / (zFar - zNear))
}

//glm.lookAt (right-handed)
func GlmLookAt(eye *[]float64, center *[]float64, up *[]float64, R *M.GslMatrix) {
	var f, s, u, SEC, CUF []float64 = []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{0, 0, 0}

	//tvec3<T, P> const f(normalize(center - eye));
	SubV(center, eye, 3, &SEC)
	Normalize(&SEC, 3, &f)
	//tvec3<T, P> const s(normalize(cross(f, up)));
	CrossV3(&f, up, &CUF)
	Normalize(&CUF, 3, &s)
	//tvec3<T, P> const u(cross(s, f));
	CrossV3(&s, &f, &u)

	M.Set(R, 0, 0, s[0])
	M.Set(R, 1, 0, s[1])
	M.Set(R, 2, 0, s[2])

	M.Set(R, 0, 1, u[0])
	M.Set(R, 1, 1, u[1])
	M.Set(R, 2, 1, u[2])

	M.Set(R, 0, 2, -f[0])
	M.Set(R, 1, 2, -f[1])
	M.Set(R, 2, 2, -f[2])

	M.Set(R, 3, 0, -Scalar(&s,eye,3))
	M.Set(R, 3, 1, -Scalar(&u,eye,3))
	M.Set(R, 3, 2, Scalar(&f,eye,3))

	M.Set(R, 3, 3, 1)
}

//+normalize [x]
//+cross [x]
//+dot [x]
//+getVectorLength [x]

// ===== TRANSFORMATION =====

//set translation coefs based on XYZ vector, opt transponse (row order / transponse for column)
func SetT(m *M.GslMatrix, x float64, y float64, z float64, t bool){
	M.Set(m, 3, 0, x)
	M.Set(m, 3, 1, y)
	M.Set(m, 3, 2, z)
	if t == true {
		M.Transpose(m)
	}
}

//set scale coefs based on XYZ vector, opt transponse (row order / transponse for column)
func SetSc(m *M.GslMatrix, x float64, y float64, z float64, t bool){
	M.Set(m, 0, 0, x)
	M.Set(m, 1, 1, y)
	M.Set(m, 2, 2, z)
	if t == true {
		M.Transpose(m)
	}
}

type Degree float64
type Radian float64

func (d Degree) RAD() float64 {
	return float64(d * (math.Pi / 180.0))
}

func (r Radian) ZeroCheck() float64 {
	if (r == 0) {
		return float64(0)
	} else {
		return float64(r)
	}
}

//set rotation by X coefs based on DEG, opt transponse (row order / transponse for column)
func SetRx(m *M.GslMatrix, d Degree, t bool){
	M.Set(m, 1, 1, Radian(math.Cos(d.RAD())).ZeroCheck())
	M.Set(m, 2, 1, Radian(math.Sin(d.RAD())).ZeroCheck())
	M.Set(m, 1, 2, Radian(-math.Sin(d.RAD())).ZeroCheck())
	M.Set(m, 2, 2, Radian(math.Cos(d.RAD())).ZeroCheck())
	if t == true {
		M.Transpose(m)
	}
}

//set rotation by Y coefs based on DEG, opt transponse (row order / transponse for column)
func SetRy(m *M.GslMatrix, d Degree, t bool){
	M.Set(m, 0, 0, Radian(math.Cos(d.RAD())).ZeroCheck())
	M.Set(m, 0, 2, Radian(math.Sin(d.RAD())).ZeroCheck())
	M.Set(m, 2, 0, Radian(-math.Sin(d.RAD())).ZeroCheck())
	M.Set(m, 2, 2, Radian(math.Cos(d.RAD())).ZeroCheck())
	if t == true {
		M.Transpose(m)
	}
}

//set rotation by Z coefs based on DEG, opt transponse (row order / transponse for column)
func SetRz(m *M.GslMatrix, d Degree, t bool){
	M.Set(m, 0, 0, Radian(math.Cos(d.RAD())).ZeroCheck())
	M.Set(m, 0, 1, Radian(math.Sin(d.RAD())).ZeroCheck())
	M.Set(m, 1, 0, Radian(-math.Sin(d.RAD())).ZeroCheck())
	M.Set(m, 1, 1, Radian(math.Cos(d.RAD())).ZeroCheck())
	if t == true {
		M.Transpose(m)
	}
}

// ===== OPERATIONS =====

//multiply matrices m1,m2, store result inro R
func Mul(m1 *M.GslMatrix, m2 *M.GslMatrix, R *M.GslMatrix){
	B.Dgemm(B.NoTrans, B.NoTrans, 1.0, m1, m2, 0.0, R)
}

//export matrix from gsl_matrix struct as double array
func Array(ma *M.GslMatrix, m int, n int, array *[16]float64){
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			array[i*m+j] = M.Get(ma,i,j)
		}
	}
}

//get single vector length Xd
func GetVLen(vec *[]float64, size int) float64 {
	var sum float64
	for i := 0; i<size; i++ {
		sum += (*vec)[i] * (*vec)[i]
	}

	return math.Sqrt(sum)
}

//normalize single vector Xd
func Normalize(vec *[]float64, size int, r *[]float64){
	var len float64 = GetVLen(vec, size)

	for i := 0; i < size; i++ {
		if len == 0 {
			(*r)[i] = 0 
		} else {
			(*r)[i] = (*vec)[i]/len
		}
	}
}

//"dot": scalar product of 2 Xd vectors
func Scalar(vec1 *[]float64, vec2 *[]float64, size int) float64 {
	var sum float64

	for i := 0; i < size; i++ {
		sum += (*vec1)[i] * (*vec2)[i]
	}

	return sum
}

//substart v1 from v2 Xd both
func SubV(vec1 *[]float64, vec2 *[]float64, size int, r *[]float64){
	for i := 0; i < size; i++ {
		(*r)[i] = (*vec1)[i] - (*vec2)[i]
	}	
}

//add v1 to v2 Xd both
func AddV(vec1 *[]float64, vec2 *[]float64, size int, r *[]float64){
	for i := 0; i < size; i++ {
		(*r)[i] = (*vec1)[i] + (*vec2)[i]
	}	
}

//mul by N vec
func MulV(vec1 *[]float64, N float64, size int, r *[]float64){
	for i := 0; i < size; i++ {
		(*r)[i] = (*vec1)[i] * N
	}	
}

//div by N vec
func DivV(vec1 *[]float64, N float64, size int, r *[]float64){
	for i := 0; i < size; i++ {
		(*r)[i] = (*vec1)[i] / N
	}	
}

//get cross of two 3d vectors
func CrossV3(vec1 *[]float64, vec2 *[]float64, r *[]float64){
	(*r)[0] = (((*vec1)[1] * (*vec2)[2]) - ((*vec1)[2] * (*vec2)[1]));
	(*r)[1] = (((*vec1)[2] * (*vec2)[0]) - ((*vec1)[0] * (*vec2)[2]));
	(*r)[2] = (((*vec1)[0] * (*vec2)[1]) - ((*vec1)[1] * (*vec2)[0]));
}

// ===== DEBUG =====

//printf MxN matrix (debug)
func PrintMatrix(ma *M.GslMatrix, m int, n int){
	for i := 0; i < n; i++ {
		fmt.Printf("\n")
		for j := 0; j < m; j++ {
			fmt.Printf("%.12f  ",M.Get(ma,i,j))
		}
	}
}