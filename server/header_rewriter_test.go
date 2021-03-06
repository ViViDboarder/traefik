package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderRewriter_Rewrite(t *testing.T) {
	testCases := []struct {
		desc       string
		remoteAddr string
		trustedIPs []string
		insecure   bool
		expected   map[string]string
	}{
		{
			desc:       "Secure & authorized",
			remoteAddr: "10.10.10.10:80",
			trustedIPs: []string{"10.10.10.10"},
			insecure:   false,
			expected: map[string]string{
				"X-Foo":           "bar",
				"X-Forwarded-For": "30.30.30.30",
			},
		},
		{
			desc:       "Secure & unauthorized",
			remoteAddr: "50.50.50.50:80",
			trustedIPs: []string{"10.10.10.10"},
			insecure:   false,
			expected: map[string]string{
				"X-Foo":           "bar",
				"X-Forwarded-For": "",
			},
		},
		{
			desc:       "Secure & authorized error",
			remoteAddr: "10.10.10.10",
			trustedIPs: []string{"10.10.10.10"},
			insecure:   false,
			expected: map[string]string{
				"X-Foo":           "bar",
				"X-Forwarded-For": "",
			},
		},
		{
			desc:       "insecure & authorized",
			remoteAddr: "10.10.10.10:80",
			trustedIPs: []string{"10.10.10.10"},
			insecure:   true,
			expected: map[string]string{
				"X-Foo":           "bar",
				"X-Forwarded-For": "30.30.30.30",
			},
		},
		{
			desc:       "insecure & unauthorized",
			remoteAddr: "50.50.50.50:80",
			trustedIPs: []string{"10.10.10.10"},
			insecure:   true,
			expected: map[string]string{
				"X-Foo":           "bar",
				"X-Forwarded-For": "30.30.30.30",
			},
		},
		{
			desc:       "insecure & authorized error",
			remoteAddr: "10.10.10.10",
			trustedIPs: []string{"10.10.10.10"},
			insecure:   true,
			expected: map[string]string{
				"X-Foo":           "bar",
				"X-Forwarded-For": "30.30.30.30",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			rewriter, err := NewHeaderRewriter(test.trustedIPs, test.insecure)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "http://20.20.20.20/foo", nil)
			require.NoError(t, err)
			req.RemoteAddr = test.remoteAddr

			req.Header.Set("X-Foo", "bar")
			req.Header.Set("X-Forwarded-For", "30.30.30.30")

			rewriter.Rewrite(req)

			for key, value := range test.expected {
				assert.Equal(t, value, req.Header.Get(key))
			}
		})
	}
}
