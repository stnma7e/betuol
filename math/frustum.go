package math

const (
	NEAR    		byte = 0
	FAR     		byte = 1
	TOP 	 		byte = 2
	RIGHT 	 		byte = 3
	BOTTOM	 		byte = 4
	LEFT			byte = 5
	SIDE_TOTAL   	byte = 6

	TOPLEFT  		byte = 0
	TOPRIGHT		byte = 1
	DOWNRIGHT   	byte = 2
	DOWNLEFT 		byte = 3

	OUTSIDE			int = 0
	INSIDE 			int = 1
	INTERSECT 		int = 2
)

type frustum struct {
	fov, aspect, nearDist, farDist float32
	lookAt, perspMat *Mat4x4
	sides [SIDE_TOTAL]Plane
}

func MakeFrustum(nearDist, farDist, fov, aspect float32) *frustum {
	var frust frustum
	frust.fov 	 	 	= fov
	frust.aspect	 	= aspect
	frust.nearDist 	 	= nearDist
	frust.farDist 	 	= farDist
	frust.lookAt 		= &Mat4x4{}; frust.lookAt.MakeIdentity()

	perspMat     	   := MakePerspectiveMatrix(nearDist,farDist,fov,aspect)
	frust.perspMat 		= &perspMat
	frust.init()
	
	return &frust
}

func (frust *frustum) init() {
		worldToClip := frust.lookAt.Mult(frust.perspMat)
		for i := 0; i < 4; i++ {
		frust.sides[LEFT][i]     = worldToClip[12+i] + worldToClip[i]
		frust.sides[RIGHT][i]	 = worldToClip[12+i] - worldToClip[i]
		frust.sides[BOTTOM][i]   = worldToClip[12+i] + worldToClip[4+i]
		frust.sides[TOP][i]      = worldToClip[12+i] - worldToClip[4+i]
		frust.sides[NEAR][i] 	 = worldToClip[12+i] + worldToClip[8+i]
		frust.sides[FAR][i] 	 = worldToClip[12+i] - worldToClip[8+i]
	}
	for i := range frust.sides {
		frust.sides[i].Normalize()
	}
}
func (frust *frustum) LookAt(target, eye, up *Vec3) {
	uPrime    := up.Normalize()
	f  		  := target.Subtract(eye)
	f.Normalize()
	s 		  := f.Cross(uPrime)
	u 		  := s.Cross(&f)

	m := Mat4x4{ s[0], s[1], s[2],0,
				 u[0], u[1], u[2],0,
				-f[0],-f[1],-f[2],0,
				 0   , 0   , 0   ,1 }
	t := Mat4x4{ 1,0,0,-eye[0],
				 0,1,0,-eye[1],
				 0,0,1,-eye[2],
				 0,0,0, 1 }

	frust.lookAt = m.Mult(&t)
	frust.init()
}

func (frust *frustum) IsPointInside(vec Vec3) bool {
	worldToClip := frust.lookAt.Mult(frust.perspMat)
	vec = Mult(vec, worldToClip)
	lookAt := Mult(Vec3{1,1,1}, frust.lookAt)
	if lookAt.Distance(&vec) > frust.farDist {
		return false
	}
	for i := range frust.sides {
		if frust.sides[i].IsInside(&vec) != true {
			return false
		}
	}

	return true
}

func (frust *frustum) IsSphereInside(sp *Sphere) int {
	for i := range frust.sides {
		distance := frust.sides[i].Distance(&sp.Center)
		if distance < -sp.Radius {
			return OUTSIDE
		}
		if distance < sp.Radius {
			return INTERSECT
		}
	}
	return INSIDE
}