package mlib

import (
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
// func GlmLookAt(eye *float64, center *float64, up *float64, R *M.GslMatrix) {
// 	var f, s, u, SEC, CUF [3]float64

// 	//tvec3<T, P> const f(normalize(center - eye));



// 	//tvec3<T, P> const s(normalize(cross(f, up)));

// 	//tvec3<T, P> const u(cross(s, f));

// }

//+normalize [x]
//+cross [x]
//+dot [x]
//+getVectorLength [x]

// ===== TRANSFORMATION =====

//set translation coefs based on XYZ vector, opt transponse (row order / transponse for column)

//set scale coefs based on XYZ vector, opt transponse (row order / transponse for column)

//set rotation by X coefs based on DEG, opt transponse (row order / transponse for column)

//set rotation by Y coefs based on DEG, opt transponse (row order / transponse for column)

//set rotation by Z coefs based on DEG, opt transponse (row order / transponse for column)

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


//"dot": scalar product of 2 Xd vectors

//substart v1 from v2 Xd both

//add v1 to v2 Xd both

//mul by N vec

//div by N vec

//get cross of two 3d vectors

// ===== DEBUG =====

//printf MxN matrix (debug)
