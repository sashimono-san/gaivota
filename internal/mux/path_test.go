package mux

import (
	"fmt"
	"testing"
)

func shouldPanic(t *testing.T, f func(), message string) {
	defer func() { recover() }()
	f()
	t.Error(message)
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path   Path
		prefix Path

		ok bool
	}{
		{path: Path("/"), prefix: Path("/"), ok: true},
		{path: Path("/:id"), prefix: Path("/"), ok: true},
		{path: Path("/positions"), prefix: Path("/"), ok: true},
		{path: Path("/:id"), prefix: Path("/123"), ok: true},
		{path: Path("/123"), prefix: Path("/:id"), ok: true},
		{path: Path("/positions"), prefix: Path("/positions"), ok: true},
		{path: Path("/positions/:id"), prefix: Path("/positions"), ok: true},
		{path: Path("/investments/some-id/positions/other-id"), prefix: Path("/investments/:id/positions"), ok: true},
		{path: Path("/positions"), prefix: Path("/positions/123123"), ok: false},
		{path: Path("/:id"), prefix: Path("/positions/123123"), ok: false},
	}

	for _, tc := range testCases {
		t.Run(string(tc.path), func(t *testing.T) {
			got := tc.path.HasPrefix(tc.prefix)

			if got != tc.ok {
				t.Errorf("Mismatch for path '%s' and prefix '%v' - got %v instead of %v", tc.path, tc.prefix, got, tc.ok)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path      Path
		candidate Path

		ok bool
	}{
		{path: Path("/"), candidate: Path("/"), ok: true},
		{path: Path("/:id"), candidate: Path("/"), ok: false},
		{path: Path("/positions"), candidate: Path("/"), ok: false},
		{path: Path("/:id"), candidate: Path("/123"), ok: true},
		{path: Path("/positions"), candidate: Path("/positions"), ok: true},
		{path: Path("/positions"), candidate: Path("positions/123"), ok: false},
	}

	for _, tc := range testCases {
		t.Run(string(tc.path), func(t *testing.T) {
			got := tc.path.Match(tc.candidate)

			if got != tc.ok {
				t.Errorf("Mismatch for path '%s' and candidate '%v' - got %v instead of %v", tc.path, tc.candidate, got, tc.ok)
			}
		})
	}
}

func TestCleanPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path   string
		wanted string
	}{
		{path: "", wanted: "/"},
		{path: "/", wanted: "/"},
		{path: "/////", wanted: "/"},
		{path: "/:id/", wanted: "/:id"},
		{path: "positions", wanted: "/positions"},
		{path: "positions///", wanted: "/positions"},
		{path: "/positions////:id", wanted: "/positions/:id"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			got := cleanPath(tc.path)

			if got != tc.wanted {
				t.Errorf("Failed to clean path %s, got %s instead of %s", tc.path, got, tc.wanted)
			}
		})
	}
}

func TestParamsPosExtract(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path Path

		shouldPanic   bool
		wantParamsPos ParamsPos
	}{
		{path: Path("/"), wantParamsPos: ParamsPos{}},
		{path: Path("/positions"), wantParamsPos: ParamsPos{}},
		{path: Path("/positions/:id"), wantParamsPos: ParamsPos{"id": 1}},
		{path: Path("/:positions/:id"), wantParamsPos: ParamsPos{"positions": 0, "id": 1}},
		{path: Path("/investments/:uuid/positions/:weirdName"), wantParamsPos: ParamsPos{"uuid": 1, "weirdName": 3}},
		{path: Path("/investments/:uuid/positions/:uuid"), shouldPanic: true},
	}

	for _, tc := range testCases {
		if tc.shouldPanic {
			shouldPanic(t, func() {
				tc.path.extractParamsPos()
			}, fmt.Sprintf("Did not panic for duplicated param in path '%s'", tc.path))

			// If panic is expected, ignore the rest of the testcase
			continue
		}

		t.Run(string(tc.path), func(t *testing.T) {
			paramsPos := tc.path.extractParamsPos()
			for param, pos := range paramsPos {
				// Test if all extracted params are wanted
				if wantPos, ok := tc.wantParamsPos[param]; ok {
					if wantPos != pos {
						t.Errorf("Param '%s' fount at position '%v' instead of '%v' for path '%s'", param, pos, wantPos, tc.path)
					}

					continue
				}

				// Only get here if found unexpected param
				t.Errorf("Found unrequested param '%s' in path '%s'", param, tc.path)
			}
		})
	}
}
