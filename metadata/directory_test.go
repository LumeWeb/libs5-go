package metadata

import (
	"encoding/json"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	cmp "github.com/LumeWeb/go-cmp"
	"github.com/vmihailenco/msgpack/v5"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func readFile(filename string) []byte {
	filePath := filepath.Join("testdata", filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return data
}

func getDirectoryMeta() *DirectoryMetadata {
	data := readFile("directory.json")

	var dir DirectoryMetadata

	err := json.Unmarshal(data, &dir)
	if err != nil {
		panic(err)
	}

	return &dir
}

func TestDirectoryMetadata_DecodeMsgpack(t *testing.T) {
	type fields struct {
		Details       DirectoryMetadataDetails
		Directories   map[string]DirectoryReference
		Files         map[string]FileReference
		ExtraMetadata ExtraMetadata
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Decode",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonDm := getDirectoryMeta()
			dm := &DirectoryMetadata{}

			if err := msgpack.Unmarshal(readFile("directory.bin"), dm); (err != nil) != tt.wantErr {
				t.Errorf("DecodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}

			fmt.Println(cmp.Diff(jsonDm, dm))
			if !cmp.Equal(jsonDm, dm) {
				t.Errorf("DecodeMsgpack() error = %v, wantErr %v", "msgpack does not match json", tt.wantErr)
			}
		})
	}
}

func TestDirectoryMetadata_EncodeMsgpack(t *testing.T) {
	type fields struct {
		Details       DirectoryMetadataDetails
		Directories   map[string]DirectoryReference
		Files         map[string]FileReference
		ExtraMetadata ExtraMetadata
	}
	type args struct {
		enc *msgpack.Encoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := &DirectoryMetadata{
				Details:       tt.fields.Details,
				Directories:   tt.fields.Directories,
				Files:         tt.fields.Files,
				ExtraMetadata: tt.fields.ExtraMetadata,
			}
			if err := dm.EncodeMsgpack(tt.args.enc); (err != nil) != tt.wantErr {
				t.Errorf("EncodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirectoryReference_DecodeMsgpack(t *testing.T) {
	type fields struct {
		Created           uint64
		Name              string
		EncryptedWriteKey []byte
		PublicKey         []byte
		EncryptionKey     []byte
		Ext               map[string]interface{}
		URI               string
		Key               string
		Size              int64
	}
	type args struct {
		dec *msgpack.Decoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := &DirectoryReference{
				Created:           tt.fields.Created,
				Name:              tt.fields.Name,
				EncryptedWriteKey: tt.fields.EncryptedWriteKey,
				PublicKey:         tt.fields.PublicKey,
				EncryptionKey:     tt.fields.EncryptionKey,
				Ext:               tt.fields.Ext,
				URI:               tt.fields.URI,
				Key:               tt.fields.Key,
				Size:              tt.fields.Size,
			}
			if err := dr.DecodeMsgpack(tt.args.dec); (err != nil) != tt.wantErr {
				t.Errorf("DecodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirectoryReference_EncodeMsgpack(t *testing.T) {
	type fields struct {
		Created           uint64
		Name              string
		EncryptedWriteKey []byte
		PublicKey         []byte
		EncryptionKey     []byte
		Ext               map[string]interface{}
		URI               string
		Key               string
		Size              int64
	}
	type args struct {
		enc *msgpack.Encoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := &DirectoryReference{
				Created:           tt.fields.Created,
				Name:              tt.fields.Name,
				EncryptedWriteKey: tt.fields.EncryptedWriteKey,
				PublicKey:         tt.fields.PublicKey,
				EncryptionKey:     tt.fields.EncryptionKey,
				Ext:               tt.fields.Ext,
				URI:               tt.fields.URI,
				Key:               tt.fields.Key,
				Size:              tt.fields.Size,
			}
			if err := dr.EncodeMsgpack(tt.args.enc); (err != nil) != tt.wantErr {
				t.Errorf("EncodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileReference_DecodeMsgpack(t *testing.T) {
	type fields struct {
		Name     string
		Created  int
		Version  int
		File     *FileVersion
		Ext      map[string]interface{}
		History  map[int]*FileVersion
		MimeType string
		URI      string
		Key      string
	}
	type args struct {
		dec *msgpack.Decoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &FileReference{
				Name:     tt.fields.Name,
				Created:  tt.fields.Created,
				Version:  tt.fields.Version,
				File:     tt.fields.File,
				Ext:      tt.fields.Ext,
				History:  tt.fields.History,
				MimeType: tt.fields.MimeType,
				URI:      tt.fields.URI,
				Key:      tt.fields.Key,
			}
			if err := fr.DecodeMsgpack(tt.args.dec); (err != nil) != tt.wantErr {
				t.Errorf("DecodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileReference_EncodeMsgpack(t *testing.T) {
	type fields struct {
		Name     string
		Created  int
		Version  int
		File     *FileVersion
		Ext      map[string]interface{}
		History  map[int]*FileVersion
		MimeType string
		URI      string
		Key      string
	}
	type args struct {
		enc *msgpack.Encoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &FileReference{
				Name:     tt.fields.Name,
				Created:  tt.fields.Created,
				Version:  tt.fields.Version,
				File:     tt.fields.File,
				Ext:      tt.fields.Ext,
				History:  tt.fields.History,
				MimeType: tt.fields.MimeType,
				URI:      tt.fields.URI,
				Key:      tt.fields.Key,
			}
			if err := fr.EncodeMsgpack(tt.args.enc); (err != nil) != tt.wantErr {
				t.Errorf("EncodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileReference_Modified(t *testing.T) {
	type fields struct {
		Name     string
		Created  int
		Version  int
		File     *FileVersion
		Ext      map[string]interface{}
		History  map[int]*FileVersion
		MimeType string
		URI      string
		Key      string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &FileReference{
				Name:     tt.fields.Name,
				Created:  tt.fields.Created,
				Version:  tt.fields.Version,
				File:     tt.fields.File,
				Ext:      tt.fields.Ext,
				History:  tt.fields.History,
				MimeType: tt.fields.MimeType,
				URI:      tt.fields.URI,
				Key:      tt.fields.Key,
			}
			if got := fr.Modified(); got != tt.want {
				t.Errorf("Modified() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileVersionThumbnail_DecodeMsgpack(t *testing.T) {
	type fields struct {
		ImageType   string
		AspectRatio float64
		CID         *encoding.EncryptedCID
		Thumbhash   []byte
	}
	type args struct {
		dec *msgpack.Decoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fvt := &FileVersionThumbnail{
				ImageType:   tt.fields.ImageType,
				AspectRatio: tt.fields.AspectRatio,
				CID:         tt.fields.CID,
				Thumbhash:   tt.fields.Thumbhash,
			}
			if err := fvt.DecodeMsgpack(tt.args.dec); (err != nil) != tt.wantErr {
				t.Errorf("DecodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileVersionThumbnail_Encode(t *testing.T) {
	type fields struct {
		ImageType   string
		AspectRatio float64
		CID         *encoding.EncryptedCID
		Thumbhash   []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   map[int]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fvt := &FileVersionThumbnail{
				ImageType:   tt.fields.ImageType,
				AspectRatio: tt.fields.AspectRatio,
				CID:         tt.fields.CID,
				Thumbhash:   tt.fields.Thumbhash,
			}
			if got := fvt.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileVersionThumbnail_EncodeMsgpack(t *testing.T) {
	type fields struct {
		ImageType   string
		AspectRatio float64
		CID         *encoding.EncryptedCID
		Thumbhash   []byte
	}
	type args struct {
		enc *msgpack.Encoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fvt := &FileVersionThumbnail{
				ImageType:   tt.fields.ImageType,
				AspectRatio: tt.fields.AspectRatio,
				CID:         tt.fields.CID,
				Thumbhash:   tt.fields.Thumbhash,
			}
			if err := fvt.EncodeMsgpack(tt.args.enc); (err != nil) != tt.wantErr {
				t.Errorf("EncodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileVersion_CID(t *testing.T) {
	type fields struct {
		Ts           int
		EncryptedCID *encoding.EncryptedCID
		PlaintextCID *encoding.CID
		Thumbnail    *FileVersionThumbnail
		Hashes       []*encoding.Multihash
		Ext          map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *encoding.CID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := &FileVersion{
				Ts:           tt.fields.Ts,
				EncryptedCID: tt.fields.EncryptedCID,
				PlaintextCID: tt.fields.PlaintextCID,
				Thumbnail:    tt.fields.Thumbnail,
				Hashes:       tt.fields.Hashes,
				Ext:          tt.fields.Ext,
			}
			if got := fv.CID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileVersion_DecodeMsgpack(t *testing.T) {
	type fields struct {
		Ts           int
		EncryptedCID *encoding.EncryptedCID
		PlaintextCID *encoding.CID
		Thumbnail    *FileVersionThumbnail
		Hashes       []*encoding.Multihash
		Ext          map[string]interface{}
	}
	type args struct {
		dec *msgpack.Decoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := &FileVersion{
				Ts:           tt.fields.Ts,
				EncryptedCID: tt.fields.EncryptedCID,
				PlaintextCID: tt.fields.PlaintextCID,
				Thumbnail:    tt.fields.Thumbnail,
				Hashes:       tt.fields.Hashes,
				Ext:          tt.fields.Ext,
			}
			if err := fv.DecodeMsgpack(tt.args.dec); (err != nil) != tt.wantErr {
				t.Errorf("DecodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileVersion_EncodeMsgpack(t *testing.T) {
	type fields struct {
		Ts           int
		EncryptedCID *encoding.EncryptedCID
		PlaintextCID *encoding.CID
		Thumbnail    *FileVersionThumbnail
		Hashes       []*encoding.Multihash
		Ext          map[string]interface{}
	}
	type args struct {
		enc *msgpack.Encoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := &FileVersion{
				Ts:           tt.fields.Ts,
				EncryptedCID: tt.fields.EncryptedCID,
				PlaintextCID: tt.fields.PlaintextCID,
				Thumbnail:    tt.fields.Thumbnail,
				Hashes:       tt.fields.Hashes,
				Ext:          tt.fields.Ext,
			}
			if err := fv.EncodeMsgpack(tt.args.enc); (err != nil) != tt.wantErr {
				t.Errorf("EncodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
