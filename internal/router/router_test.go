package router

import "testing"

func TestExtractParams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		route    *Route
		path     string
		expected map[string]string
	}{
		{
			route: &Route{
				ParamsPos: ParamsPos{},
			},
			path:     "/",
			expected: map[string]string{},
		},
		{
			route: &Route{
				ParamsPos: ParamsPos{"id": 0},
			},
			path:     "/123",
			expected: map[string]string{"id": "123"},
		},
		{
			route: &Route{
				ParamsPos: ParamsPos{"uuid": 1},
			},
			path:     "/positions/uuid-here",
			expected: map[string]string{"uuid": "uuid-here"},
		},
		{
			route: &Route{
				ParamsPos: ParamsPos{},
			},
			path:     "/foo/bar",
			expected: map[string]string{},
		},
		{
			route: &Route{
				ParamsPos: ParamsPos{"investimentID": 1, "positionID": 3},
			},
			path:     "/investiment/15/positions/foo",
			expected: map[string]string{"investimentID": "15", "positionID": "foo"},
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.path), func(t *testing.T) {
			gotParams := tc.route.ExtractParams(tc.path)

			if len(gotParams) != len(tc.expected) {
				t.Errorf("Incorrect number of extracted params for path '%s'. Got %v instead of %v", tc.path, len(gotParams), len(tc.expected))
			}
			for param, val := range gotParams {
				if pos, ok := tc.expected[param]; !ok {
					t.Errorf("Got unexpected param for path '%s': %v at position %v", tc.path, param, pos)
				}

				if val != tc.expected[param] {
					t.Errorf("Unexpected value for param %s in path '%s'. Got %s and expected %s", param, tc.path, val, tc.expected[param])
				}
			}
		})
	}
}
