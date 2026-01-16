package main

import "testing"

func TestCleanInput(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected []string
	}{
		{
			Input:    "command",
			Expected: []string{"command"},
		},
		{
			Input:    "command",
			Expected: []string{"command"},
		},
		{
			Input:    "command arg1",
			Expected: []string{"command", "arg1"},
		},
		{
			Input:    "command arg1 arg2",
			Expected: []string{"command", "arg1", "arg2"},
		},
		{
			Input:    "    command arg1 arg2    ",
			Expected: []string{"command", "arg1", "arg2"},
		},
		{
			Input:    "ComMand aRG1",
			Expected: []string{"command", "arg1"},
		},
		{
			Input:    "ComMand      aRG1",
			Expected: []string{"command", "arg1"},
		},
	}

	for _, test := range testCases {
		actual := cleanInput(test.Input)

		if len(actual) != len(test.Expected) {
			t.Errorf("actual and expected lists are different sizes. actual %d != expected %d", len(actual), len(test.Expected))
			continue
		}

		for j, word := range actual {
			if word != test.Expected[j] {
				t.Errorf("result words do not match. actual %s != expected %s", word, test.Expected[j])
				break
			}
		}
	}
}
