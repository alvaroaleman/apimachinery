/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sets

import (
	"reflect"
	"testing"
)

type stringSetInterface interface {
	Insert(...string) stringSetInterface
	Delete(...string) stringSetInterface
	Has(string) bool
	HasAll(...string) bool
	HasAny(...string) bool
	List() []string
	IsSuperset(stringSetInterface) bool
	Difference(stringSetInterface) stringSetInterface
	Equal(stringSetInterface) bool
	Union(stringSetInterface) stringSetInterface
	Intersection(stringSetInterface) stringSetInterface

	// Not actually implemented by the sets, but no other way to access this
	Len() int
}

type stringSetInterfaceAdapter struct {
	String
}

func (s *stringSetInterfaceAdapter) Insert(items ...string) stringSetInterface {
	s.String.Insert(items...)
	return s
}

func (s *stringSetInterfaceAdapter) Delete(items ...string) stringSetInterface {
	s.String.Delete(items...)
	return s
}

func (s *stringSetInterfaceAdapter) IsSuperset(other stringSetInterface) bool {
	return s.String.IsSuperset(other.(*stringSetInterfaceAdapter).String)
}

func (s *stringSetInterfaceAdapter) Difference(other stringSetInterface) stringSetInterface {
	return &stringSetInterfaceAdapter{s.String.Difference(other.(*stringSetInterfaceAdapter).String)}
}

func (s *stringSetInterfaceAdapter) Equal(other stringSetInterface) bool {
	return s.String.Equal(other.(*stringSetInterfaceAdapter).String)
}

func (s *stringSetInterfaceAdapter) Union(other stringSetInterface) stringSetInterface {
	return &stringSetInterfaceAdapter{s.String.Union(other.(*stringSetInterfaceAdapter).String)}
}
func (s *stringSetInterfaceAdapter) Intersection(other stringSetInterface) stringSetInterface {
	return &stringSetInterfaceAdapter{s.String.Intersection(other.(*stringSetInterfaceAdapter).String)}
}

type genericStringSetInterfaceAdapter struct {
	Set[string]
}

func (s *genericStringSetInterfaceAdapter) Insert(items ...string) stringSetInterface {
	s.Set.Insert(items...)
	return s
}

func (s *genericStringSetInterfaceAdapter) Delete(items ...string) stringSetInterface {
	s.Set.Delete(items...)
	return s
}

func (s *genericStringSetInterfaceAdapter) IsSuperset(other stringSetInterface) bool {
	return s.Set.IsSuperset(other.(*genericStringSetInterfaceAdapter).Set)
}

func (s *genericStringSetInterfaceAdapter) Difference(other stringSetInterface) stringSetInterface {
	return &genericStringSetInterfaceAdapter{s.Set.Difference(other.(*genericStringSetInterfaceAdapter).Set)}
}

func (s *genericStringSetInterfaceAdapter) Equal(other stringSetInterface) bool {
	return s.Set.Equal(other.(*genericStringSetInterfaceAdapter).Set)
}

func (s *genericStringSetInterfaceAdapter) Union(other stringSetInterface) stringSetInterface {
	return &genericStringSetInterfaceAdapter{s.Set.Union(other.(*genericStringSetInterfaceAdapter).Set)}
}

func (s *genericStringSetInterfaceAdapter) Intersection(other stringSetInterface) stringSetInterface {
	return &genericStringSetInterfaceAdapter{s.Set.Intersection(other.(*genericStringSetInterfaceAdapter).Set)}
}

func (s *genericStringSetInterfaceAdapter) Len() int {
	return len(s.Set)
}

type stringSetConstructor func(...string) stringSetInterface

// TestStringSet runs all tests for both the String and the Set implementation to make
// sure they behave the same.
func TestStringSet(t *testing.T) {
	implementations := []struct {
		name        string
		constructor stringSetConstructor
	}{
		{
			name: "generated",
			constructor: func(items ...string) stringSetInterface {
				return &stringSetInterfaceAdapter{NewString(items...)}
			},
		},
		{
			name: "generic",
			constructor: func(items ...string) stringSetInterface {
				return &genericStringSetInterfaceAdapter{New(items...)}
			},
		},
	}

	tests := []func(*testing.T, stringSetConstructor){
		testStringSet,
		testStringSetDeleteMultiples,
		testNewStringSet,
		testStringSetList,
		testStringSetDifference,
		testStringSetHasAny,
		testStringSetEquals,
		testStringUnion,
		testStringIntersection,
	}

	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(reflect.TypeOf(test).Name(), func(t *testing.T) {
					test(t, implementation.constructor)
				})
			}
		})
	}
}

func testStringSet(t *testing.T, constructor stringSetConstructor) {
	s := constructor()
	s2 := constructor()
	if s.Len() != 0 {
		t.Errorf("Expected len=0: %d", s.Len())
	}
	s.Insert("a", "b")
	if s.Len() != 2 {
		t.Errorf("Expected len=2: %d", s.Len())
	}
	s.Insert("c")
	if s.Has("d") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Has("a") {
		t.Errorf("Missing contents: %#v", s)
	}
	s.Delete("a")
	if s.Has("a") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s.Insert("a")
	if s.HasAll("a", "b", "d") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.HasAll("a", "b") {
		t.Errorf("Missing contents: %#v", s)
	}
	s2.Insert("a", "b", "d")
	if s.IsSuperset(s2) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s2.Delete("d")
	if !s.IsSuperset(s2) {
		t.Errorf("Missing contents: %#v", s)
	}
}

func testStringSetDeleteMultiples(t *testing.T, constructor stringSetConstructor) {
	s := constructor()
	s.Insert("a", "b", "c")
	if s.Len() != 3 {
		t.Errorf("Expected len=3: %d", s.Len())
	}

	s.Delete("a", "c")
	if s.Len() != 1 {
		t.Errorf("Expected len=1: %d", s.Len())
	}
	if s.Has("a") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if s.Has("c") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Has("b") {
		t.Errorf("Missing contents: %#v", s)
	}

}

func testNewStringSet(t *testing.T, constructor stringSetConstructor) {
	s := constructor("a", "b", "c")
	if s.Len() != 3 {
		t.Errorf("Expected len=3: %d", s.Len())
	}
	if !s.Has("a") || !s.Has("b") || !s.Has("c") {
		t.Errorf("Unexpected contents: %#v", s)
	}
}

func testStringSetList(t *testing.T, constructor stringSetConstructor) {
	s := constructor("z", "y", "x", "a")
	if !reflect.DeepEqual(s.List(), []string{"a", "x", "y", "z"}) {
		t.Errorf("List gave unexpected result: %#v", s.List())
	}
}

func testStringSetDifference(t *testing.T, constructor stringSetConstructor) {
	a := constructor("1", "2", "3")
	b := constructor("1", "2", "4", "5")
	c := a.Difference(b)
	d := b.Difference(a)
	if c.Len() != 1 {
		t.Errorf("Expected len=1: %d", c.Len())
	}
	if !c.Has("3") {
		t.Errorf("Unexpected contents: %#v", c.List())
	}
	if d.Len() != 2 {
		t.Errorf("Expected len=2: %d", d.Len())
	}
	if !d.Has("4") || !d.Has("5") {
		t.Errorf("Unexpected contents: %#v", d.List())
	}
}

func testStringSetHasAny(t *testing.T, constructor stringSetConstructor) {
	a := constructor("1", "2", "3")

	if !a.HasAny("1", "4") {
		t.Errorf("expected true, got false")
	}

	if a.HasAny("0", "4") {
		t.Errorf("expected false, got true")
	}
}

func testStringSetEquals(t *testing.T, constructor stringSetConstructor) {
	// Simple case (order doesn't matter)
	a := constructor("1", "2")
	b := constructor("2", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// It is a set; duplicates are ignored
	b = constructor("2", "2", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// Edge cases around empty sets / empty strings
	a = constructor()
	b = constructor()
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	b = constructor("1", "2", "3")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	b = constructor("1", "2", "")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	// Check for equality after mutation
	a = constructor()
	a.Insert("1")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert("2")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert("")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	a.Delete("")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}
}

func testStringUnion(t *testing.T, constuctor stringSetConstructor) {
	tests := []struct {
		s1       stringSetInterface
		s2       stringSetInterface
		expected stringSetInterface
	}{
		{
			constuctor("1", "2", "3", "4"),
			constuctor("3", "4", "5", "6"),
			constuctor("1", "2", "3", "4", "5", "6"),
		},
		{
			constuctor("1", "2", "3", "4"),
			constuctor(),
			constuctor("1", "2", "3", "4"),
		},
		{
			constuctor(),
			constuctor("1", "2", "3", "4"),
			constuctor("1", "2", "3", "4"),
		},
		{
			constuctor(),
			constuctor(),
			constuctor(),
		},
	}

	for _, test := range tests {
		union := test.s1.Union(test.s2)
		if union.Len() != test.expected.Len() {
			t.Errorf("Expected union.Len()=%d but got %d", test.expected.Len(), union.Len())
		}

		if !union.Equal(test.expected) {
			t.Errorf("Expected union.Equal(expected) but not true.  union:%v expected:%v", union.List(), test.expected.List())
		}
	}
}

func testStringIntersection(t *testing.T, constuctor stringSetConstructor) {
	tests := []struct {
		s1       stringSetInterface
		s2       stringSetInterface
		expected stringSetInterface
	}{
		{
			constuctor("1", "2", "3", "4"),
			constuctor("3", "4", "5", "6"),
			constuctor("3", "4"),
		},
		{
			constuctor("1", "2", "3", "4"),
			constuctor("1", "2", "3", "4"),
			constuctor("1", "2", "3", "4"),
		},
		{
			constuctor("1", "2", "3", "4"),
			constuctor(),
			constuctor(),
		},
		{
			constuctor(),
			constuctor("1", "2", "3", "4"),
			constuctor(),
		},
		{
			constuctor(),
			constuctor(),
			constuctor(),
		},
	}

	for _, test := range tests {
		intersection := test.s1.Intersection(test.s2)
		if intersection.Len() != test.expected.Len() {
			t.Errorf("Expected intersection.Len()=%d but got %d", test.expected.Len(), intersection.Len())
		}

		if !intersection.Equal(test.expected) {
			t.Errorf("Expected intersection.Equal(expected) but not true.  intersection:%v expected:%v", intersection.List(), test.expected.List())
		}
	}
}
