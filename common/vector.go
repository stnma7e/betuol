package common

type Vector struct {
	array []interface{}
	emptyIndices IntQueue
	Length int
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

func (vc *Vector) Push_back(data interface{}, resizeStep, checkDistance int) {
	if cap(vc.array) <= vc.Length + checkDistance {
		tmp := vc.array
		vc.Length = len(tmp)
		vc.array = make([]interface{}, vc.Length + resizeStep)
		for i := range tmp {
			vc.array[i] = tmp[i]
		}
	}
	vc.array[vc.Length] = data
}

func (vc *Vector) Insert(data interface{}) {
    space, err := vc.emptyIndices.Dequeue()
    if err != nil {
        vc.Push_back(data, 1,1)
    } else {
        vc.array[space] = data
    }
}

func (vc *Vector) Erase(index int) {
	vc.array[index] = nil
        vc.emptyIndices.Queue(index)
}

func (vc *Vector) Empty() {
    vc.array = make([]interface{}, 0)
    vc.Length = 0
}

func (vc *Vector) IsEmpty() bool {
	return vc.Length == 0
}
