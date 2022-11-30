package values

import (
	"testing"

	"github.com/divolgin/helmsplain/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestGetFromData(t *testing.T) {
	log.SetDebug(false)

	tests := []struct {
		name     string
		template string
		want     []Value
	}{
		{
			name:     "basic template with pipe",
			template: `<h1>{{.Values.images.tag | quote }} {{ .age }}</h1>`,
			want: []Value{
				{
					Key: ".Values.images.tag",
					Pos: 6,
				},
			},
		},
		{
			name:     "basic template with function",
			template: `<h1>{{ (quote .Values.images.tag) }}</h1>`,
			want: []Value{
				{
					Key: ".Values.images.tag",
					Pos: 21,
				},
			},
		},
		{
			name:     "using conditionals",
			template: `<h1>{{ if true }} some text {{ (quote .Values.images.tag) }} some more text {{else}} {{ .Values.registry.host }} {{end}}`,
			want: []Value{
				{
					Key: ".Values.images.tag",
					Pos: 45,
				},
				{
					Key: ".Values.registry.host",
					Pos: 88,
				},
			},
		},
		{
			name:     "using 'with' in template with function",
			template: `<h1>{{with .Values}} some text {{ (quote .images.tag) }} some more text {{end}}`,
			want: []Value{
				{
					Key: ".Values.images.tag",
					Pos: 48,
				},
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
