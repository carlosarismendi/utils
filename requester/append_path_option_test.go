package requester

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAppendPath(t *testing.T) {
	type testCase struct {
		name     string
		url      string
		path     string
		expected string
	}
	testsCases := []testCase{
		{
			name:     "appendPathWithStartingSlashToURLWithSlash",
			url:      "localhost:8080/",
			path:     "/test",
			expected: "localhost:8080/test",
		},
		{
			name:     "appendPathWithDoubleStartingSlashToURLWithSlash",
			url:      "localhost:8080/",
			path:     "//test",
			expected: "localhost:8080/test",
		},
		{
			name:     "appendPathWithStartingSlashToURLWithoutSlash",
			url:      "localhost:8080",
			path:     "/test",
			expected: "localhost:8080/test",
		},
		{
			name:     "appendPathWithDoubleStartingSlashToURLWithoutSlash",
			url:      "localhost:8080",
			path:     "//test",
			expected: "localhost:8080/test",
		},

		{
			name:     "appendPathWithoutStartingSlashToURLWithSlash",
			url:      "localhost:8080/",
			path:     "test",
			expected: "localhost:8080/test",
		},
		{
			name:     "appendPathWithoutStartingSlashToURLWithSlash",
			url:      "localhost:8080/",
			path:     "test",
			expected: "localhost:8080/test",
		},
		{
			name:     "appendPathWithoutStartingSlashToURLWithoutSlash",
			url:      "localhost:8080",
			path:     "test",
			expected: "localhost:8080/test",
		},
		{
			name:     "appendPathWithoutStartingSlashToURLWithoutSlash",
			url:      "localhost:8080",
			path:     "test",
			expected: "localhost:8080/test",
		},
	}

	for _, test := range testsCases {
		t.Run(test.name, func(t *testing.T) {
			// ARRANGE
			req := NewRequester(URL(test.url), AppendPath(test.path))

			// ACT
			httpReq, err := req.prepareRequest()

			// ASSERT
			require.NoError(t, err)
			require.Equal(t, test.expected, httpReq.URL.String())
		})
	}
}
