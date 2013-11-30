// Package common is composed of data structures and other common structures to ease development.
package common

// Vector is a resizeable array. It takes its name from the C++ std::vector class.
type Vector struct {
	array        []interface{}
	emptyIndices Queue
	Length       int
}

// MakeVector returns a pointer to a Vector
func MakeVector() *Vector {
	return &Vector{
		make([]interface{}, 0),
		Queue{}, 0,
	}
}

// Array returns the underlying array used by the vector
func (vc *Vector) Array() []interface{} {
	return vc.array
}

// Push_back will append a data structure to the end of the vector.
// If the length of the vector is less than the space required for the new data structure, then the vector is extended and the data is copied into the new, bigger vector.
// Then the new data structure is appended to the end.
func (vc *Vector) Push_back(data interface{}, resizeStep, checkDistance int) {
	if cap(vc.array) <= vc.Length+checkDistance {
		tmp := vc.array
		vc.array = make([]interface{}, vc.Length+resizeStep)
		for i := range tmp {
			vc.array[i] = tmp[i]
		}
	}
	vc.array[vc.Length] = data
	vc.Length++
}

// Insert will attempt to insert a value into the vector using an empty slot.
// If no empty space is found, Insert resorts to Push_back to append the data structure to the end of the vector.
func (vc *Vector) Insert(data interface{}) int {
	space, err := vc.emptyIndices.Dequeue()
	if err != nil {
		vc.Push_back(data, 1, 1)
		return vc.Length - 1
	}

	vc.array[space.(int)] = data
	return space.(int)
}

// Difference will return a vector that is composed of the data structures that were unique in one of the two input vectors.
func (vec1 *Vector) Difference(vec2 *Vector) *Vector {
	ret := MakeVector()
	list1 := vec1.Array()
	list2 := vec2.Array()
	for i := range list2 {
		if func() bool {
			for j := range list1 {
				if list2[i] == list1[j] {
					return false
				}
			}
			return true
		}() {
			ret.Insert(list2[i])
		}
	}
	return ret
}

// Erase will delete the data from an array position in the vector.
func (vc *Vector) Erase(index int) {
	vc.array[index] = nil
	vc.emptyIndices.Queue(index)
	vc.Length--
}

// Empty will delete all data in the vector and essentially create a new vector.
func (vc *Vector) Empty() {
	vc.array = make([]interface{}, 0)
	vc.Length = 0
}

func (vc *Vector) IsEmpty() bool {
	return vc.Length < 1
}
