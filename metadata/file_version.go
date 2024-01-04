package metadata

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"github.com/vmihailenco/msgpack/v5"
)

type FileVersion struct {
	Ts           int                    `json:"ts"`
	EncryptedCID *encoding.EncryptedCID `json:"encryptedCID,string"`
	PlaintextCID *encoding.CID          `json:"plaintextCID,string"`
	Thumbnail    *FileVersionThumbnail  `json:"thumbnail"`
	Hashes       []*encoding.Multihash  `json:"hashes"`
	Ext          map[string]interface{} `json:"ext"`
}

func NewFileVersion(ts int, encryptedCID *encoding.EncryptedCID, plaintextCID *encoding.CID, thumbnail *FileVersionThumbnail, hashes []*encoding.Multihash, ext map[string]interface{}) *FileVersion {
	return &FileVersion{
		Ts:           ts,
		EncryptedCID: encryptedCID,
		PlaintextCID: plaintextCID,
		Thumbnail:    thumbnail,
		Hashes:       hashes,
		Ext:          ext,
	}
}

func (fv *FileVersion) EncodeMsgpack(enc *msgpack.Encoder) error {
	data := map[int]interface{}{
		8: fv.Ts,
	}

	if fv.EncryptedCID != nil {
		data[1] = fv.EncryptedCID.ToBytes()
	}

	if fv.PlaintextCID != nil {
		data[2] = fv.PlaintextCID.ToBytes()
	}

	if len(fv.Hashes) > 0 {
		hashesData := make([][]byte, len(fv.Hashes))
		for i, hash := range fv.Hashes {
			hashesData[i] = hash.FullBytes
		}
		data[9] = hashesData
	}

	if fv.Thumbnail != nil {
		data[10] = fv.Thumbnail.Encode()
	}

	return enc.Encode(data)
}

func (fv *FileVersion) DecodeMsgpack(dec *msgpack.Decoder) error {
	mapLen, err := dec.DecodeMapLen()

	if err != nil {
		return err
	}

	for i := 0; i < mapLen; i++ {
		key, err := dec.DecodeInt8()
		if err != nil {
			return err
		}
		switch key {

		case int8(1):
			err := dec.Decode(&fv.EncryptedCID)
			if err != nil {
				return err
			}

		case int8(2):
			err := dec.Decode(&fv.PlaintextCID)
			if err != nil {
				return err
			}
		case int8(8):
			err := dec.Decode(&fv.Ts)
			if err != nil {
				return err
			}
		case int8(9):
			hashesData, err := dec.DecodeSlice()
			if err != nil {
				return err
			}

			fv.Hashes = make([]*encoding.Multihash, len(hashesData))
			for i, hashData := range hashesData {
				hashBytes := hashData.([]byte)
				fv.Hashes[i] = encoding.NewMultihash(hashBytes)
			}

		case int8(10):
			err := dec.Decode(&fv.Thumbnail)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fv *FileVersion) CID() *encoding.CID {
	if fv.PlaintextCID != nil {
		return fv.PlaintextCID
	}
	return &fv.EncryptedCID.OriginalCID
}
