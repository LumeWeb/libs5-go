package cid

import (
	"git.lumeweb.com/LumeWeb/libs5-go/cid/multibase"
	"git.lumeweb.com/LumeWeb/libs5-go/internal/testdata"
	"git.lumeweb.com/LumeWeb/libs5-go/multihash"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"reflect"
	"testing"
)

func TestCID_CopyWith(t *testing.T) {
	type fields struct {
		Multibase multibase.Multibase
		Type      types.CIDType
		Hash      multihash.Multihash
		Size      uint32
	}
	type args struct {
		newType int
		newSize uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cid := &CID{
				Multibase: tt.fields.Multibase,
				Type:      tt.fields.Type,
				Hash:      tt.fields.Hash,
				Size:      tt.fields.Size,
			}
			got, err := cid.CopyWith(tt.args.newType, tt.args.newSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyWith() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CopyWith() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCID_Equals(t *testing.T) {
	type fields struct {
		Multibase multibase.Multibase
		Type      types.CIDType
		Hash      multihash.Multihash
		Size      uint32
	}
	type args struct {
		other *CID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cid := &CID{
				Multibase: tt.fields.Multibase,
				Type:      tt.fields.Type,
				Hash:      tt.fields.Hash,
				Size:      tt.fields.Size,
			}
			if got := cid.Equals(tt.args.other); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCID_HashCode(t *testing.T) {
	type fields struct {
		Multibase multibase.Multibase
		Type      types.CIDType
		Hash      multihash.Multihash
		Size      uint32
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
			cid := &CID{
				Multibase: tt.fields.Multibase,
				Type:      tt.fields.Type,
				Hash:      tt.fields.Hash,
				Size:      tt.fields.Size,
			}
			if got := cid.HashCode(); got != tt.want {
				t.Errorf("HashCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*func TestCID_ToBytes(t *testing.T) {
	FromHash(testdata.RawBase58CID)

	println(len(testdata.RawCIDBytes))
	println(utils.DecodeEndian(testdata.RawCIDBytes[35:]))
	return
	type fields struct {
		Type types.CIDType
		Hash multihash.Multihash
		Size uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
				{
				name: "Bridge CID",
				fields: fields{
					Type: types.CIDTypeBridge,
					Hash: multihash.New(), // Replace with a valid hash value
				},
				want: , // Replace with the expected byte output for Bridge CID
			},
		{
			name: "Raw CID with Non-Zero Size",
			fields: fields{
				Type: types.CIDTypeRaw,
				Hash: *multihash.New(testdata.RawCIDBytes[1:34]),
				Size: utils.DecodeEndian(testdata.RawCIDBytes[34:]),
			},
			want: testdata.RawCIDBytes,
		},
			{
				name: "Raw CID with Zero Size",
				fields: fields{
					Type: types.CIDTypeRaw,
					Hash: yourHashValue, // Replace with a valid hash value
					Size: 0,             // Zero size
				},
				want: yourExpectedBytesForRawCIDWithZeroSize, // Replace with the expected byte output
			},
			{
				name: "Default CID",
				fields: fields{
					Type: types.CIDTypeDefault,
					Hash: yourHashValue, // Replace with a valid hash value
				},
				want: yourExpectedBytesForDefaultCID, // Replace with the expected byte output for Default CID
			},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cid := &CID{
				Type: tt.fields.Type,
				Hash: tt.fields.Hash,
				Size: tt.fields.Size,
			}
			m := multibase.New(cid)
			cid.Multibase = m
			if got := cid.ToBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}*/

func TestCID_ToRegistryCID(t *testing.T) {
	type fields struct {
		Multibase multibase.Multibase
		Type      types.CIDType
		Hash      multihash.Multihash
		Size      uint32
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cid := &CID{
				Multibase: tt.fields.Multibase,
				Type:      tt.fields.Type,
				Hash:      tt.fields.Hash,
				Size:      tt.fields.Size,
			}
			got, err := cid.ToRegistryCID()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToRegistryCID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRegistryCID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCID_ToRegistryEntry(t *testing.T) {
	type fields struct {
		Multibase multibase.Multibase
		Type      types.CIDType
		Hash      multihash.Multihash
		Size      uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cid := &CID{
				Multibase: tt.fields.Multibase,
				Type:      tt.fields.Type,
				Hash:      tt.fields.Hash,
				Size:      tt.fields.Size,
			}
			if got := cid.ToRegistryEntry(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRegistryEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCID_ToString(t *testing.T) {
	type fields struct {
		Multibase multibase.Multibase
		Type      types.CIDType
		Hash      multihash.Multihash
		Size      uint32
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cid := &CID{
				Multibase: tt.fields.Multibase,
				Type:      tt.fields.Type,
				Hash:      tt.fields.Hash,
				Size:      tt.fields.Size,
			}
			got, err := cid.ToString()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCID_getPrefixBytes(t *testing.T) {
	type fields struct {
		Multibase multibase.Multibase
		Type      types.CIDType
		Hash      multihash.Multihash
		Size      uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cid := &CID{
				Multibase: tt.fields.Multibase,
				Type:      tt.fields.Type,
				Hash:      tt.fields.Hash,
				Size:      tt.fields.Size,
			}
			if got := cid.getPrefixBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPrefixBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		cid string
	}
	tests := []struct {
		name    string
		args    args
		want    *CID
		wantErr bool
	}{
		/*        {
		          name:    "Valid Bridge CID",
		          args:args  {cid: ""},
		          want:    nil,
		          wantErr: false,
		      },*/

		{
			name:    "Valid Raw Base 58 CID",
			args:    args{cid: testdata.RawBase58CID},
			want:    New(types.CIDTypeRaw, *multihash.New(testdata.RawCIDBytes[1:34]), testdata.RawCIDSize),
			wantErr: false,
		},
		{
			name:    "Valid Media 58 CID",
			args:    args{cid: testdata.MediaBase58CID},
			want:    New(types.CIDTypeMetadataMedia, *multihash.New(testdata.MediaCIDBytes[1:34]), testdata.MediaCIDSize),
			wantErr: false,
		},
		{
			name:    "Valid Resolver CID",
			args:    args{cid: testdata.ResolverBase58CID},
			want:    New(types.CIDTypeResolver, *multihash.New(testdata.ResolverCIDBytes[1:34]), testdata.ResolverCIDSize),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeCID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeCID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromBytes(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *CID
		wantErr bool
	}{
		{
			name:    "Valid Raw Base 58 CID",
			args:    args{bytes: testdata.RawCIDBytes},
			want:    New(types.CIDTypeRaw, *multihash.New(testdata.RawCIDBytes[1:34]), testdata.RawCIDSize),
			wantErr: false,
		},
		{
			name:    "Valid Media 58 CID",
			args:    args{bytes: testdata.MediaCIDBytes},
			want:    New(types.CIDTypeMetadataMedia, *multihash.New(testdata.MediaCIDBytes[1:34]), testdata.MediaCIDSize),
			wantErr: false,
		},
		{
			name:    "Valid Resolver CID",
			args:    args{bytes: testdata.ResolverCIDBytes},
			want:    New(types.CIDTypeResolver, *multihash.New(testdata.ResolverCIDBytes[1:34]), testdata.ResolverCIDSize),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromBytes(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromHash(t *testing.T) {
	type args struct {
		bytes   interface{}
		size    uint32
		cidType types.CIDType
	}
	tests := []struct {
		name    string
		args    args
		want    *CID
		wantErr bool
	}{
		{
			name:    "Valid Raw Base 58 CID",
			args:    args{bytes: testdata.RawCIDBytes[1:34], size: testdata.RawCIDSize, cidType: types.CIDTypeRaw},
			want:    New(types.CIDTypeRaw, *multihash.New(testdata.RawCIDBytes[1:34]), testdata.RawCIDSize),
			wantErr: false,
		},
		{
			name:    "Valid Media 58 CID",
			args:    args{bytes: testdata.MediaCIDBytes[1:34], size: testdata.MediaCIDSize, cidType: types.CIDTypeMetadataMedia},
			want:    New(types.CIDTypeMetadataMedia, *multihash.New(testdata.MediaCIDBytes[1:34]), testdata.MediaCIDSize),
			wantErr: false,
		},
		{
			name:    "Valid Resolver CID",
			args:    args{bytes: testdata.ResolverCIDBytes[1:34], size: testdata.ResolverCIDSize, cidType: types.CIDTypeResolver},
			want:    New(types.CIDTypeResolver, *multihash.New(testdata.ResolverCIDBytes[1:34]), testdata.ResolverCIDSize),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromHash(tt.args.bytes, tt.args.size, tt.args.cidType)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromRegistry(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *CID
		wantErr bool
	}{
		{
			name:    "Valid Resolver Data",
			args:    args{bytes: testdata.ResolverDataBytes},
			want:    New(types.CIDTypeMetadataWebapp, *multihash.New(testdata.ResolverDataBytes[2:35]), testdata.ResolverDataSize),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromRegistry(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromRegistry() got = %v, want %v", got, tt.want)
			}
		})
	}
}
