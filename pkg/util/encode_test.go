package util

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "standard",
			want: "[\"foo\",\"bar\"]",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(
				http.MethodGet,
				"/test",
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			Encode(rr, req, http.StatusOK, []string{"foo", "bar"})

			got := strings.TrimSpace(rr.Body.String())
			if got != tc.want {
				t.Errorf("Got: %s, Want: %s", got, tc.want)
			}
		})
	}

}
