package common

type Vector struct {
	array        []interface{}
	emptyIndices IntQueue
	Length       int
}

func MakeVector() *Vector {
	return &Vector{
		make([]interface{}, 0),
		IntQueue{}, 0,
	}
}

func (vc *Vector) Array() []interface{} {
	return vc.array
}

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

func (vc *Vector) Insert(data interface{}) int {
	space, err := vc.emptyIndices.Dequeue()
	if err != nil {
		vc.Push_back(data, 1, 1)
		return vc.Length - 1
	} else {
		vc.array[space] = data
		return space
	}
}

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
			ret.Insert(i)
		}
	}
	return ret
}

func (vc *Vector) Erase(index int) {
	vc.array[index] = nil
	vc.emptyIndices.Queue(index)
	vc.Length--
}

func (vc *Vector) Empty() {
	vc.array = make([]interface{}, 0)
	vc.Length = 0
}

func (vc *Vector) IsEmpty() bool {
	return vc.Length < 1
}
