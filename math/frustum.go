package math

const (
	NEAR		        byte = 0
	FAR		        byte = 1
	TOP			byte = 2
	RIGHT			byte = 3
	BOTTOM			byte = 4
	LEFT			byte = 5
	SIDE_TOTAL	        byte = 6

	TOPLEFT		        byte = 0
	TOPRIGHT		byte = 1
	DOWNRIGHT	        byte = 2
	DOWNLEFT		byte = 3

	OUTSIDE			int = 0
	INSIDE			int = 1
	INTERSECT		int = 2
)

type Frustum struct {
	fov, aspect, nearDist, farDist float32
	lookAt, perspMat Mat4x4
	sides [SIDE_TOTAL]Plane
}

func MakeFrustum(nearDist, farDist, fov, aspect float32) *Frustum {
	var frust Frustum
	frust.fov		= fov
	frust.aspect		= aspect
	frust.nearDist		= nearDist
	frust.farDist		= farDist
	frust.lookAt		= Mat4x4{}; frust.lookAt.MakeIdentity()

	perspMat	   := MakePerspectiveMatrix(nearDist,farDist,fov,aspect)
	frust.perspMat	    = perspMat
	frust.init()

	return &frust
}

func (frust *Frustum) init() {
		worldToClip := Mult4m4m(frust.perspMat, frust.lookAt)
		for i := 0; i < 4; i++ {
		frust.sides[LEFT][i]     = worldToClip[12+i] + worldToClip[i]
		frust.sides[RIGHT][i]	 = worldToClip[12+i] - worldToClip[i]
		frust.sides[BOTTOM][i]   = worldToClip[12+i] + worldToClip[4+i]
		frust.sides[TOP][i]      = worldToClip[12+i] - worldToClip[4+i]
		frust.sides[NEAR][i]	 = worldToClip[12+i] + worldToClip[8+i]
		frust.sides[FAR][i]	 = worldToClip[12+i] - worldToClip[8+i]
	}
	for i := range frust.sides {
		frust.sides[i].Normalize()
	}
}

func (frust *Frustum) LookAt(target, eye, up Vec3) {
	u := Normalize3v(up)
	f := Normalize3v(Sub3v3v(target, eye))
	s := Normalize3v(Cross3v3v(f, u))
	u  = Cross3v3v(s, f)

	frust.lookAt = Mat4x4{ s[0], s[1], s[2],-Dot3v3v(s,eye),
				 		   u[0], u[1], u[2],-Dot3v3v(u,eye),
				          -f[0],-f[1],-f[2], Dot3v3v(f,eye),
				 		   0,    0,    0,    1 }
	frust.init()
}

func (frust *Frustum) ContainsPoint(vec Vec3) bool {
	worldToClip := Mult4m4m(frust.lookAt, frust.perspMat)
	vec = Mult4m3v(worldToClip, vec)
	lookAt := Mult4m3v(frust.lookAt, Vec3{})
	if Dist3v3v(lookAt, vec) > frust.farDist {
		return false
	}
	for i := range frust.sides {
		if frust.sides[i].IsInside(vec) != true {
			return false
		}
	}

	return true
}

func (frust *Frustum) ContainsSphere(sp Sphere) int {
	for i := range frust.sides {
		distance := frust.sides[i].Distance(sp.Center)
		if distance < -sp.Radius {
			return OUTSIDE
		}
		if distance < sp.Radius {
			return INTERSECT
		}
	}
	return INSIDE
}

func (frust *Frustum) LookAtMatrix() Mat4x4 {
	return frust.lookAt
}
func (frust *Frustum) Projection() Mat4x4 {
	return frust.perspMat
}
