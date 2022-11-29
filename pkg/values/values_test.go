package values

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFromData(t *testing.T) {
	tests := []struct {
		name     string
		template string
		want     []string
	}{
		{
			name:     "basic template",
			template: `<h1>{{.Values.images.tag | quote }} {{ .age }}</h1>`,
			want: []string{
				".Values.images.tag",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := require.New(t)
			got := GetFromData(tt.template)
			req.Equal(tt.want, got)
		})
	}
}
