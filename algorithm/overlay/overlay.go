package overlay

import (
	"github.com/spatial-go/geoos/algorithm"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/algorithm/relate"
)

// Overlay  Computes the overlay of two geometries,either or both of which may be nil.
type Overlay interface {
	// Union  Computes the Union of two geometries,either or both of which may be nil.
	Union() (matrix.Steric, error)
	// Intersection  Computes the Intersection of two geometries,either or both of which may be nil.
	Intersection() (matrix.Steric, error)
	// Difference returns a geometry that represents that part of geometry A that does not intersect with geometry B.
	// One can think of this as GeometryA - Intersection(A,B).
	// If A is completely contained in B then an empty geometry collection is returned.
	Difference() (matrix.Steric, error)
	// SymDifference returns a geometry that represents the portions of A and B that do not intersect.
	// It is called a symmetric difference because SymDifference(A,B) = SymDifference(B,A).
	// One can think of this as Union(geomA,geomB) - Intersection(A,B).
	SymDifference() (matrix.Steric, error)
}

// PointOverlay  Computes the overlay of two geometries,either or both of which may be nil.
type PointOverlay struct {
	Subject, Clipping matrix.Steric
}

// Union  Computes the Union of two geometries,either or both of which may be nil.
func (p *PointOverlay) Union() (matrix.Steric, error) {
	if res, ok := p.unionCheck(); !ok {
		return res, nil
	}
	if s, ok := p.Subject.(matrix.Matrix); ok {
		if c, ok := p.Clipping.(matrix.Matrix); ok {
			if s.Equals(c) {
				return s, nil
			}
			return matrix.Collection{s, c}, nil
		}
	}
	return nil, algorithm.ErrNotMatchType
}

// Intersection  Computes the Intersection of two geometries,either or both of which may be nil.
func (p *PointOverlay) Intersection() (matrix.Steric, error) {
	if res, ok := p.intersectionCheck(); !ok {
		return res, nil
	}
	if _, ok := p.Subject.(matrix.Matrix); !ok {
		return nil, algorithm.ErrNotMatchType
	}
	switch c := p.Clipping.(type) {
	case matrix.Matrix:
		if p.Subject.(matrix.Matrix).Equals(c) {
			return p.Subject.(matrix.Matrix), nil
		}
		return nil, nil
	case matrix.LineMatrix:
		if mark := relate.InLineMatrix(p.Subject.(matrix.Matrix), c); mark {
			return p.Subject.(matrix.Matrix), nil
		}
		return nil, nil
	case matrix.PolygonMatrix:
		inter := envelope.Bound(c.Bound()).IsIntersects(envelope.Bound(p.Subject.Bound()))
		if mark := relate.IM(c, p.Subject.(matrix.Matrix), inter).IsCovers(); mark {
			return p.Subject.(matrix.Matrix), nil
		}
		return nil, nil
	}
	return nil, algorithm.ErrNotMatchType
}

// Difference returns a geometry that represents that part of geometry A that does not intersect with geometry B.
// One can think of this as GeometryA - Intersection(A,B).
// If A is completely contained in B then an empty geometry collection is returned.
func (p *PointOverlay) Difference() (matrix.Steric, error) {
	if res, ok := p.differenceCheck(); !ok {
		return res, nil
	}
	if s, ok := p.Subject.(matrix.Matrix); ok {
		if c, ok := p.Clipping.(matrix.Matrix); ok {
			if s.Equals(c) {
				return nil, nil
			}
			return s, nil
		}
	}
	return nil, algorithm.ErrNotMatchType
}

// DifferenceReverse returns a geometry that represents reverse that part of geometry A that does not intersect with geometry B .
// One can think of this as GeometryB - Intersection(A,B).
// If B is completely contained in A then an empty geometry collection is returned.
func (p *PointOverlay) DifferenceReverse() (matrix.Steric, error) {
	newP := &PointOverlay{Subject: p.Clipping, Clipping: p.Subject}
	return newP.Difference()
}

// SymDifference returns a geometry that represents the portions of A and B that do not intersect.
// It is called a symmetric difference because SymDifference(A,B) = SymDifference(B,A).
// One can think of this as Union(geomA,geomB) - Intersection(A,B).
func (p *PointOverlay) SymDifference() (matrix.Steric, error) {
	result := matrix.Collection{}
	if res, err := p.Difference(); err == nil {
		result = append(result, res)
	}
	if res, err := p.DifferenceReverse(); err == nil {
		result = append(result, res)
	}
	return result, nil
}

// unionCheck  Computes the Union of two geometries,either or both of which may be null.
func (p *PointOverlay) unionCheck() (matrix.Steric, bool) {

	if p.Subject == nil && p.Clipping == nil {
		return nil, false
	}
	if p.Subject == nil {
		return p.Clipping, false
	}

	if p.Clipping == nil {
		return p.Subject, false
	}

	return nil, true
}

// intersectionCheck  Computes the Union of two geometries,either or both of which may be null.
func (p *PointOverlay) intersectionCheck() (matrix.Steric, bool) {

	if p.Subject == nil && p.Clipping == nil {
		return nil, false
	}
	if p.Subject == nil {
		return nil, false
	}

	if p.Clipping == nil {
		return nil, false
	}

	return nil, true
}

// differenceCheck  Computes the Union of two geometries,either or both of which may be null.
func (p *PointOverlay) differenceCheck() (matrix.Steric, bool) {

	if p.Subject == nil && p.Clipping == nil {
		return nil, false
	}
	if p.Subject == nil {
		return nil, false
	}

	if p.Clipping == nil {
		return p.Subject, false
	}

	return nil, true
}
