package common

import "fmt"

type Vector struct {
	array []interface{}
	emptyIndices IntQueue
	length uint
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
	if cap(vc.array) <= int(vc.length + checkDistance) {
		tmp := vc.array
		vc.length = uint(len(tmp))
		vc.array = make([]interface{}, vc.length + resizeStep)
		for i := range tmp {
			vc.array[i] = tmp[i]
		}
	}

	fmt.Println(vc.length, len(vc.array))
	vc.array[vc.length] = data
}

func (vc *Vector) Erase(index int) {
	vc.array[index] = nil
}