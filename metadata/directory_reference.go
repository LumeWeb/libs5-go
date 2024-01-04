package metadata

import "github.com/vmihailenco/msgpack/v5"

var _ SerializableMetadata = (*DirectoryReference)(nil)

type DirectoryReference struct {
	Created           uint64                 `json:"created"`
	Name              string                 `json:"name"`
	EncryptedWriteKey Base64UrlBinary        `json:"encryptedWriteKey,string"`
	PublicKey         Base64UrlBinary        `json:"publicKey,string"`
	EncryptionKey     Base64UrlBinary        `json:"encryptionKey,string"`
	Ext               map[string]interface{} `json:"ext"`
	URI               string                 `json:"uri"`
	Key               string                 `json:"key"`
	Size              int64                  `json:"size"`
}

func NewDirectoryReference(created uint64, name string, encryptedWriteKey, publicKey, encryptionKey []byte, ext map[string]interface{}) *DirectoryReference {
	return &DirectoryReference{
		Created:           created,
		Name:              name,
		EncryptedWriteKey: encryptedWriteKey,
		PublicKey:         publicKey,
		EncryptionKey:     encryptionKey,
		Ext:               ext,
		URI:               "",
		Key:               "",
		Size:              0,
	}
}

func (dr *DirectoryReference) EncodeMsgpack(enc *msgpack.Encoder) error {
	data := map[int]interface{}{
		1: dr.Name,
		2: dr.Created,
		3: dr.PublicKey,
		4: dr.EncryptedWriteKey,
	}

	if dr.EncryptionKey != nil {
		data[5] = dr.EncryptionKey
	}

	if dr.Ext != nil {
		data[6] = dr.Ext
	}

	return enc.Encode(data)
}

func (dr *DirectoryReference) DecodeMsgpack(dec *msgpack.Decoder) error {
	var (
		err error
		l   int
	)
	if l, err = dec.DecodeMapLen(); err != nil {
		return err
	}

	for i := 0; i < l; i++ {
		key, err := dec.DecodeInt8()
		if err != nil {
			return err
		}
		value, err := dec.DecodeInterface()
		if err != nil {
			return err
		}
		switch key {
		case int8(1):
			dr.Name = value.(string)
		case int8(2):
			dr.Created = value.(uint64)
		case int8(3):
			dr.PublicKey = value.([]byte)
		case int8(4):
			dr.EncryptedWriteKey = value.([]byte)
		case int8(5):
			dr.EncryptionKey = value.([]byte)
		case int8(6):
			dr.Ext = value.(map[string]interface{})
		}
	}

	return nil
}
