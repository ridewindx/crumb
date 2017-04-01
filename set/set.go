package set

import (
	"fmt"
	"strings"
)

type Set map[interface{}]struct{}

func NewSet() Set {
	return make(Set)
}

func (s Set) Add(item interface{}) {
	s[item] = struct{}{}
}

func (s Set) Remove(item interface{}) {
	delete(s, item)
}

func (s Set) Contains(item interface{}) bool {
	_, exists := s[item]
	return exists
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Clear() {
	for item := range s {
		delete(s, item)
	}
}

func (s Set) ToSlice() []interface{} {
	items := make([]interface{}, 0, s.Size())
	for item := range s {
		items = append(items, item)
	}
	return items
}

func (s Set) Clone() Set {
	cloned := NewSet()
	for item := range s {
		cloned.Add(item)
	}
	return cloned
}

func (s Set) String() string {
	items := make([]string, 0, s.Size())
	for item := range s {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ","))
}

func (s Set) Equal(other Set) bool {
	if s.Size() != other.Size() {
		return false
	}
	for item := range s {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

func (s Set) IsSubset(other Set) bool {
	for item := range s {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

func (s Set) IsSuperset(other Set) bool {
	return other.IsSubset(s)
}

func (s Set) Union(other Set) Set {
	unioned := s.Clone()
	for item := range other {
		unioned.Add(item)
	}
	return unioned
}

func (s Set) Intersect(other Set) Set {
	intersected := NewSet()
	if s.Size() > other.Size() {
		for item := range other {
			if s.Contains(item) {
				intersected.Add(item)
			}
		}
	} else {
		for item := range s {
			if other.Contains(item) {
				intersected.Add(item)
			}
		}
	}
	return intersected
}

func (s Set) Difference(other Set) Set {
	diff := NewSet()
	for item := range s {
		if !other.Contains(item) {
			diff.Add(item)
		}
	}
	return diff
}

func (s Set) SymmetricDifference(other Set) Set {
	symmetric := s.Difference(other)
	for item := range other.Difference(s) {
		symmetric.Add(item)
	}
	return symmetric
}
