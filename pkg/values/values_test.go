package values

import (
	"testing"

	"github.com/divolgin/helmsplain/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestGetFromData(t *testing.T) {
	log.SetDebug(true)

	tests := []struct {
		name     string
		template string
		want     []string
	}{
		{
			name:     "basic template with pipe",
			template: `<h1>{{.Values.images.tag | quote }} {{ .age }}</h1>`,
			want: []string{
				".Values.images.tag",
			},
		},
		{
			name:     "basic template with function",
			template: `<h1>{{ (quote .Values.images.tag) }}</h1>`,
			want: []string{
				".Values.images.tag",
			},
		},
		{
			name:     "using conditionals",
			template: `<h1>{{ if true }} some text {{ (quote .Values.images.tag) }} some more text {{else}} {{ .Values.registry.host }} {{end}}`,
			want: []string{
				".Values.images.tag",
				".Values.registry.host",
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
