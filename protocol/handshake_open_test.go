package protocol

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
	"testing"
)

func TestHandshakeOpen_EncodeMsgpack(t *testing.T) {
	type fields struct {
		challenge []byte
		networkId string
	}
	type args struct {
		enc *msgpack.Encoder
		buf *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Normal Case",
			fields: fields{
				challenge: []byte("test-challenge"),
				networkId: "test-network",
			},
			args: args{
				buf: new(bytes.Buffer),
			},
			wantErr: false,
		},
		{
			name: "Empty Network ID",
			fields: fields{
				challenge: []byte("test-challenge"),
				networkId: "",
			},
			args: args{
				buf: new(bytes.Buffer),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt.args.enc = msgpack.NewEncoder(tt.args.buf)
		t.Run(tt.name, func(t *testing.T) {
			h := HandshakeOpen{
				challenge: tt.fields.challenge,
				networkId: tt.fields.networkId,
			}
			if err := h.EncodeMsgpack(tt.args.enc); (err != nil) != tt.wantErr {
				t.Errorf("EncodeMsgpack() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check the contents of the buffer to verify encoding
			encodedData := tt.args.buf.Bytes()
			if len(encodedData) == 0 && !tt.wantErr {
				t.Errorf("Expected non-empty encoded data, got empty")
			}

			dec := msgpack.NewDecoder(bytes.NewReader(encodedData))

			protocolMethod, err := dec.DecodeUint()
			if err != nil {
				t.Errorf("DecodeUint() error = %v", err)
			}

			assert.EqualValues(t, types.ProtocolMethodHandshakeOpen, protocolMethod)

			challenge, err := dec.DecodeBytes()
			if err != nil {
				t.Errorf("DecodeBytes() error = %v", err)
			}

			assert.EqualValues(t, tt.fields.challenge, challenge)

			networkId, err := dec.DecodeString()
			if err != nil {
				if err.Error() == "EOF" && tt.fields.networkId != "" {
					t.Logf("DecodeString() networkId missing, got EOF")
				}
				if err.Error() != "EOF" {
					t.Errorf("DecodeString() error = %v", err)
				}
			}

			assert.EqualValues(t, tt.fields.networkId, networkId)
		})
	}
}
func TestHandshakeOpen_DecodeMessage(t *testing.T) {
	type fields struct {
		challenge []byte
		networkId string
		handshake []byte
	}
	type args struct {
		base64EncodedData string // Base64 encoded msgpack data
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "Valid Handshake and NetworkID",
			fields: fields{}, // Fields are now empty
			args: args{
				base64EncodedData: "xBNzYW1wbGVIYW5kc2hha2VEYXRhr3NhbXBsZU5ldHdvcmtJRA==",
			},
			wantErr: assert.NoError,
		},
		{
			name:   "Valid Handshake and Empty NetworkID",
			fields: fields{}, // Fields are now empty
			args: args{
				base64EncodedData: "xBNzYW1wbGVIYW5kc2hha2VEYXRh",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandshakeOpen{}

			decodedData, _ := base64.StdEncoding.DecodeString(tt.args.base64EncodedData)
			reader := bytes.NewReader(decodedData)
			dec := msgpack.NewDecoder(reader)

			tt.wantErr(t, h.DecodeMessage(dec), fmt.Sprintf("DecodeMessage(%v)", tt.args.base64EncodedData))
		})
	}
}
