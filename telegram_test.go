package main

import "testing"

func TestRenderMessage(t *testing.T) {
	tests := []struct {
		name string
		req  SendRequest
		want string
	}{
		{
			name: "level and title",
			req:  SendRequest{Message: "Build finished", Title: "CI", Level: "success"},
			want: "✅ CI\nBuild finished",
		},
		{
			name: "level only",
			req:  SendRequest{Message: "Build finished", Level: "success"},
			want: "✅ Build finished",
		},
		{
			name: "title only",
			req:  SendRequest{Message: "Build finished", Title: "CI"},
			want: "CI\nBuild finished",
		},
		{
			name: "plain",
			req:  SendRequest{Message: "Build finished"},
			want: "Build finished",
		},
		{
			name: "info level",
			req:  SendRequest{Message: "Running", Level: "info"},
			want: "ℹ️ Running",
		},
		{
			name: "warning level",
			req:  SendRequest{Message: "Slow", Level: "warning"},
			want: "⚠️ Slow",
		},
		{
			name: "error level",
			req:  SendRequest{Message: "Failed", Level: "error"},
			want: "❌ Failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RenderMessage(tc.req)
			if got != tc.want {
				t.Errorf("RenderMessage() = %q, want %q", got, tc.want)
			}
		})
	}
}
