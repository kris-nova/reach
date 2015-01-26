package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestArgTargets(t *testing.T) {
	type pair struct {
		targets []Target
		err     error
	}

	ustr := []string{
		"http://example.com",
		"http://blog.example.com/",
		"://malformed@/",
	}

	l := len(ustr)

	urls := make([]*url.URL, l)
	errs := make([]error, l)

	for i := 0; i < l; i++ {
		urls[i], errs[i] = url.ParseRequestURI(ustr[i])
	}

	argss := [][]string{
		[]string{},
		[]string{ustr[0]},
		[]string{ustr[0], ustr[1]},
		[]string{ustr[2]},
	}

	expecteds := []pair{
		pair{[]Target{}, fmt.Errorf("Please supply at least one URL.")},
		pair{[]Target{Target{urls[0]}}, nil},
		pair{[]Target{Target{urls[0]}, Target{urls[1]}}, nil},
		pair{[]Target{}, errs[2]},
	}

	for i := range argss {
		a, e := argTargets(argss[i])
		act := pair{a, e}
		exp := expecteds[i]

		shouldError := exp.err != nil
		didError := act.err != nil

		if shouldError != didError {
			if shouldError {
				t.Errorf("expected error %q ", exp.err)
			} else {
				t.Errorf("did not expect error %q", act.err)
			}
		}

		if shouldError && didError {
			if act.err.Error() != exp.err.Error() {
				t.Errorf("expected error %q, got %q", exp.err, act.err)
			}
		}

		if len(exp.targets) != len(act.targets) {
			t.Errorf("size of exp and act targets mismatch on:\n\texp: %q\n\tact: %q", exp, act)
		}
		for j := range exp.targets {
			expVal := exp.targets[j].String()
			actVal := act.targets[j].String()
			if expVal != actVal {
				t.Errorf("expected actual value %q to equal expected value %q",
					act.targets[j], exp.targets[j])
			}
		}
	}
}

func TestDropEmpties(t *testing.T) {
	input := []string{"not empty", "", "also not empty"}
	expected := fmt.Sprintf("%q", []string{"not empty", "also not empty"})
	actual := fmt.Sprintf("%q", dropEmpties(input))
	if expected != actual {
		t.Errorf("dropEmpties failed!\nExpected:\n\t%s\nGot:\n\t%s\n",
			expected, actual)
	}
}
