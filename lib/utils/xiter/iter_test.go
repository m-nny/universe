package xiter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJoin(t *testing.T) {
	testCases := []struct {
		name string
		arr  []string
		sep  string
		want string
	}{
		{
			name: "empty array with separator",
			arr:  []string{},
			sep:  "|",
			want: "",
		},
		{
			name: "one element array with separator",
			arr:  []string{"abc"},
			sep:  "|",
			want: "abc",
		},
		{
			name: "multiple elements array with separator",
			arr:  []string{"abc", "123", "xyz"},
			sep:  "|",
			want: "abc|123|xyz",
		},
		{
			name: "empty array without separator",
			arr:  []string{},
			sep:  "",
			want: "",
		},
		{
			name: "one element array without separator",
			arr:  []string{"abc"},
			sep:  "",
			want: "abc",
		},
		{
			name: "multiple elements array without separator",
			arr:  []string{"abc", "123", "xyz"},
			sep:  "",
			want: "abc123xyz",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			seq := Iter(tc.arr)
			got := Join(seq, tc.sep)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("Join() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSliceJoinFn(t *testing.T) {
	type Item struct{ value string }
	mapper := func(item Item) string { return item.value }
	testCases := []struct {
		name string
		arr  []Item
		sep  string
		want string
	}{
		{
			name: "empty array with separator",
			arr:  []Item{},
			sep:  "|",
			want: "",
		},
		{
			name: "one element array with separator",
			arr:  []Item{{"abc"}},
			sep:  "|",
			want: "abc",
		},
		{
			name: "multiple elements array with separator",
			arr:  []Item{{"abc"}, {"123"}, {"xyz"}},
			sep:  "|",
			want: "abc|123|xyz",
		},
		{
			name: "empty array without separator",
			arr:  []Item{},
			sep:  "",
			want: "",
		},
		{
			name: "one element array without separator",
			arr:  []Item{{"abc"}},
			sep:  "",
			want: "abc",
		},
		{
			name: "multiple elements array without separator",
			arr:  []Item{{"abc"}, {"123"}, {"xyz"}},
			sep:  "",
			want: "abc123xyz",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := SliceJoinFn(tc.arr, tc.sep, mapper)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("Join() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
