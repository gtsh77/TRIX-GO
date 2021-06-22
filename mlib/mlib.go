package mlib

import (
	M "github.com/dtromb/gogsl/matrix"
)

// ===== INITIALIZATION =====

//new empty matrix MxN
func m_new (m uint32, n uint32) *M.GslMatrix {
	var ma *ma.GslMatrix = M.MatrixAlloc(m,n)
	M.SetZero(ma)

	return ma
}

//new diag matrix MxM

//reset to diag matrix

//set matrix elements to X

// ===== GLM LIKE FUNCS =====

//glm.perspective

//glm.lookAt (right-handed)

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

//export matrix from gsl_matrix struct as double array

//get single vector length Xd

//normalize single vector Xd

//"dot": scalar product of 2 Xd vectors

//substart v1 from v2 Xd both

//add v1 to v2 Xd both

//mul by N vec

//div by N vec

//get cross of two 3d vectors

// ===== DEBUG =====

//printf MxN matrix (debug)