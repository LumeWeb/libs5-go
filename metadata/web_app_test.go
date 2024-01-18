package metadata

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getWebappMeta() *WebAppMetadata {
	data := getDirectoryMetaContent()

	var webapp WebAppMetadata

	err := json.Unmarshal(data, &webapp)
	if err != nil {
		panic(err)
	}

	return &webapp
}

func getWebappContent() []byte {
	data := readFile("webapp.json")

	return data
}

func TestWebAppMetadata_EncodeJSON(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Encode",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonDm := getWebappContent()
			dm := &WebAppMetadata{}

			if err := json.Unmarshal(jsonDm, dm); (err != nil) != tt.wantErr {
				t.Errorf("EncodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			jsonData, err := json.MarshalIndent(dm, "", "\t")

			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			buf := bytes.NewBuffer(nil)

			err = json.Indent(buf, jsonData, "", "\t")
			if err != nil {
				t.Errorf("EncodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, buf.Bytes(), jsonData)
		})
	}
}
