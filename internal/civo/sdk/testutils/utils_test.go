package testutils

import "testing"

func Test_InjectIntoSlice(t *testing.T) {
	t.Run("inject a value into a slice", func(t *testing.T) {
		t.Parallel()

		slice := []int{1, 2, 3}
		value := []int{4}

		InjectIntoSlice(t, &slice, value)

		AssertEqual(t, len(slice), 4)
		AssertEqual(t, slice[3], 4)
		t.Logf("slice: %#v", slice)
	})

	t.Run("inject a value into an empty slice", func(t *testing.T) {
		t.Parallel()

		slice := []int{}
		value := []int{4}

		InjectIntoSlice(t, &slice, value)

		AssertEqual(t, len(slice), 1)
		AssertEqual(t, slice[0], 4)
		t.Logf("slice: %#v", slice)
	})

	t.Run("inject a value into a nil slice", func(t *testing.T) {
		t.Parallel()

		var slice []int
		value := []int{4}

		InjectIntoSlice(t, &slice, value)

		AssertEqual(t, len(slice), 1)
		AssertEqual(t, slice[0], 4)
		t.Logf("slice: %#v", slice)
	})
}
