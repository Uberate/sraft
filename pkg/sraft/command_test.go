package sraft

import (
	"github.com/stretchr/testify/assert"
	"github.io/uberate/sraft/pkg/plugins/storage"
	"strings"
	"testing"
)

func TestCommand_Exec(t *testing.T) {

	const storageType = storage.MemoryV1EngineName
	oldSets := map[string]string{
		"test1": "test1",
		"test2": "test2",
	}

	type fields struct {
		FullPath string
		Value    string
		Action   string
	}
	type args struct {
		storage storage.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    map[string]string
	}{
		{
			name: "Test common set function",
			fields: fields{
				FullPath: "test",
				Value:    "test",
				Action:   CommandActionSet,
			},
			want: map[string]string{
				"test1": "test1",
				"test2": "test2",
				"test":  "test",
			},
		},
		{
			name: "Test cover set",
			fields: fields{
				FullPath: "test1",
				Value:    "test",
				Action:   CommandActionSet,
			},
			want: map[string]string{
				"test1": "test",
				"test2": "test2",
			},
		},
		{
			name: "Test common put function",
			fields: fields{
				FullPath: "test",
				Value:    "test",
				Action:   CommandActionPut,
			},
			want: map[string]string{
				"test1": "test1",
				"test2": "test2",
				"test":  "test",
			},
		},
		{
			name: "Test common put function with upper action value",
			fields: fields{
				FullPath: "test",
				Value:    "test",
				Action:   strings.ToUpper(CommandActionPut),
			},
			want: map[string]string{
				"test1": "test1",
				"test2": "test2",
				"test":  "test",
			},
		},
		{
			name: "Test common put cover",
			fields: fields{
				FullPath: "test1",
				Value:    "test",
				Action:   CommandActionPut,
			},
			want: map[string]string{
				"test1": "test1",
				"test2": "test2",
			},
		},
		{
			name: "Test common delete function",
			fields: fields{
				FullPath: "test1",
				Value:    "test",
				Action:   CommandActionDelete,
			},
			want: map[string]string{
				"test2": "test2",
			},
		},
		{
			name: "Test delete not value",
			fields: fields{
				FullPath: "test",
				Value:    "test",
				Action:   CommandActionDelete,
			},
			want: map[string]string{
				"test1": "test1",
				"test2": "test2",
			},
		},
		{
			name: "Test common delete function",
			fields: fields{
				FullPath: "test1",
				Value:    "test",
				Action:   "Want error action",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := storage.GetStorageEngine(storageType)
			_ = s.Clean()
			for key, value := range oldSets {
				s.Set(key, value)
			}
			cmd := &Command{
				FullPath: tt.fields.FullPath,
				Value:    tt.fields.Value,
				Action:   tt.fields.Action,
			}
			err := cmd.Exec(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			resMap := map[string]string{}
			for _, path := range s.Paths("") {
				v, _ := s.Get(path)
				resMap[path] = v
			}
			for key, value := range tt.want {
				if sValue, ok := s.Get(key); !ok || sValue != value {
					t.Errorf("Want: [%v], but got: [%v]", tt.want, resMap)
				}
			}
			assert.Equal(t, tt.want, resMap)
		})
	}
}
