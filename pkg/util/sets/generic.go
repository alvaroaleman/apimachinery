package sets

import (
	"reflect"
	"sort"
)

type Set[T comparable] map[T]struct{}

func (s Set[T]) Insert(keys ...T) {
	for _, k := range keys {
		s[k] = struct{}{}
	}
}

func (s Set[T]) Delete(keys ...T) {
	for _, k := range keys {
		delete(s, k)
	}
}

func (s Set[T]) Has(k T) bool {
	_, ok := s[k]
	return ok
}

func (s Set[T]) HasAll(keys ...T) bool {
	for _, k := range keys {
		if !s.Has(k) {
			return false
		}
	}

	return true
}

func (s Set[T]) IsSuperset(s2 Set[T]) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}

	return true
}

func (s Set[T]) List() []T {
	// using `make` here fails compilation: "invalid argument: cannot make T: no core type"
	var res []T
	for k := range s {
		res = append(res, k)
	}

	var less func(i, j int) bool

	switch reflect.TypeOf(s).Key().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		less = func(i, j int) bool {
			return reflect.ValueOf(res[i]).Int() < reflect.ValueOf(res[j]).Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		less = func(i, j int) bool {
			return reflect.ValueOf(res[i]).Uint() < reflect.ValueOf(res[j]).Uint()
		}
	case reflect.Float32, reflect.Float64:
		less = func(i, j int) bool {
			return reflect.ValueOf(res[i]).Float() < reflect.ValueOf(res[j]).Float()
		}
	case reflect.String:
		less = func(i, j int) bool {
			return reflect.ValueOf(res[i]).String() < reflect.ValueOf(res[j]).String()
		}
	}

	if less != nil {
		sort.Slice(res, less)
	}

	return res
}

func (s Set[T]) Difference(s2 Set[T]) Set[T] {
	result := New[T]()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}

	return result
}

func (s Set[T]) HasAny(items ...T) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}

	return false
}

func (s Set[T]) Equal(s2 Set[T]) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

func (s Set[T]) Union(s2 Set[T]) Set[T] {
	result := New[T]()
	for key := range s {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}

	return result
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Intersection(s2 Set[T]) Set[T] {
	var walk, other Set[T]
	result := New[T]()
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}

	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}

	return result
}

func New[T comparable](keys ...T) Set[T] {
	s := Set[T]{}
	for _, k := range keys {
		s.Insert(k)
	}

	return s
}
