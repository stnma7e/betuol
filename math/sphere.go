package math

type Sphere struct {
	Center Vec3
	Radius float32
}

func (sp *Sphere) Intersects(sp2 Sphere) bool {
	radius := sp.Radius + sp2.Radius
	if sp.Center.DistanceSqrd(sp2.Center) > radius*radius {
		return false
	}

	return true
}