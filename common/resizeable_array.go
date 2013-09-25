package common

type Vector struct {
	array []interface{}
	emptyIndices IntQueue
	Length uint
}

func MakeVector() *Vector {
	return &Vector{
		make([]interface{},0),
		IntQueue{}, 0,
	}
}

func (vc *Vector) Array() []interface{} {
	return vc.array
}

func (vc *Vector) Push_back(data interface{}, resizeStep, checkDistance uint) {
	if cap(vc.array) <= int(vc.Length + checkDistance) {
		tmp := vc.array
		vc.Length = uint(len(tmp))
		vc.array = make([]interface{}, vc.Length + resizeStep)
		for i := range tmp {
			vc.array[i] = tmp[i]
		}
	}
	vc.array[vc.Length] = data
}

func (vc *Vector) Erase(index int) {
	vc.array[index] = nil
}