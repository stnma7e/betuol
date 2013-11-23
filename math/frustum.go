package math

// Represents the specific planes/sides of the frustum.
const (
	NEAR       byte = iota
	FAR        byte = iota
	TOP        byte = iota
	RIGHT      byte = iota
	BOTTOM     byte = iota
	LEFT       byte = iota
	SIDE_TOTAL byte = iota
)

// Represents the possible intersection attributes of a sphere or other solid object.
const (
	OUTSIDE   int = iota
	INSIDE    int = iota
	INTERSECT int = iota
)

// Frustum, as defined by wikipedia, (plural: frusta or frustums) is the portion of a solid (normally a cone or pyramid) that lies between two parallel planes cutting it.
// In this case it represents a section of a pyramid lying between two parallel planes.
// It can also be visualized as a cube that has been applied a perspective projection.
type Frustum struct {
	fov, aspect, nearDist, farDist float32
	lookAt, perspMat               Mat4x4
	sides                          [SIDE_TOTAL]Plane
}

// MakeFrustum returns a pointer to a Frustum.
func MakeFrustum(nearDist, farDist, fov, aspect float32) *Frustum {
	var frust Frustum
	frust.fov = fov
	frust.aspect = aspect
	frust.nearDist = nearDist
	frust.farDist = farDist
	frust.lookAt = Mat4x4{}
	frust.lookAt.MakeIdentity()

	frust.perspMat = MakePerspectiveMatrix(nearDist, farDist, fov, aspect)
	frust.init()

	return &frust
}

// init reinitializes the planes of the frustum to match changed lookat or projection matrices.
func (frust *Frustum) init() {
	worldToClip := Mult4m4m(frust.perspMat, frust.lookAt)
	for i := 0; i < 4; i++ {
		frust.sides[LEFT][i] = worldToClip[12+i] + worldToClip[i]
		frust.sides[RIGHT][i] = worldToClip[12+i] - worldToClip[i]
		frust.sides[BOTTOM][i] = worldToClip[12+i] + worldToClip[4+i]
		frust.sides[TOP][i] = worldToClip[12+i] - worldToClip[4+i]
		frust.sides[NEAR][i] = worldToClip[12+i] + worldToClip[8+i]
		frust.sides[FAR][i] = worldToClip[12+i] - worldToClip[8+i]
	}
	for i := range frust.sides {
		frust.sides[i].Normalize()
	}
}

// LookAt will orient the frustum's planes to accomodate a change in the target to face, location of the frustum in world space, or the direction of up in world space.
// This function takes parameters of the target location for the frustum to face in world space, the location of the frustum in world space, and a vector that represents the direction of up in the world space.
func (frust *Frustum) LookAt(target, eye, up Vec3) {
	u := Normalize3v(up)
	f := Normalize3v(Sub3v3v(target, eye))
	s := Normalize3v(Cross3v3v(f, u))
	u = Cross3v3v(s, f)

	frust.lookAt = Mat4x4{s[0], s[1], s[2], -Dot3v3v(s, eye),
		u[0], u[1], u[2], -Dot3v3v(u, eye),
		-f[0], -f[1], -f[2], Dot3v3v(f, eye),
		0, 0, 0, 1}
	frust.init()
}

// ContainsPoint returns true if a point lies within the bounds of the frustum.
// This is used for culling points for rendering.
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

// ContainsSphere returns one state of OUTSIDE, INTERSECT, or INSIDE that represents how a sphere interacts with the frustum.
// If no points within the sphere lie within the frustum, then the function returns OUTSIDE.
// If any point within the sphere is inside the frustum, and any one other point within the sphere is outside the frustum, then the function returns INTERSECT.
// If none of the previous conditions are met, then the sphere lies fully within the frustum, and the function returns INSIDE.
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

// LookAtMatrix returns the lookAt matrix of the frustum.
func (frust *Frustum) LookAtMatrix() Mat4x4 {
	return frust.lookAt
}

// Projection returns the projection matrix of the frustum.
func (frust *Frustum) Projection() Mat4x4 {
	return frust.perspMat
}
