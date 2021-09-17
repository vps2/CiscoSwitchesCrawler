package set_test

import (
	"reflect"
	"testing"

	"github.com/vps2/cisco-switches-crawler/pkg/set"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		in   func(s *set.Set[int])
		out  []int
	}{
		{
			"Add unique values",
			func(s *set.Set[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
			},
			[]int{1, 2, 3},
		}, {
			"Add duplicates",
			func(s *set.Set[int]) {
				s.Add(1)
				s.Add(1)
				s.Add(2)
				s.Add(2)
				s.Add(3)
				s.Add(3)
			},
			[]int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := set.New[int]()
			tt.in(set)
			got := set.ToSlice()
			if ok := reflect.DeepEqual(got, tt.out); !ok {
				t.Errorf("Add() = %v, want %v", got, tt.out)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	set := set.New[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)

	set.Remove(1)

	got := set.ToSlice()
	want := []int{2, 3, 4}

	if ok := reflect.DeepEqual(got, want); !ok {
		t.Errorf("Remove() = %v, want %v", got, want)
	}
}

func TestContains(t *testing.T) {
	set := set.New[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	got := set.Has(2)
	want := true
	if got != want {
		t.Errorf("Has() = %v, want %v", got, want)
	}

	got = set.Has(4)
	want = false
	if got != want {
		t.Errorf("Has() = %v, want %v", got, want)
	}

}
