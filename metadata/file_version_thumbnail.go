package metadata

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"github.com/vmihailenco/msgpack/v5"
)

type FileVersionThumbnail struct {
	ImageType   string
	AspectRatio float64
	CID         *encoding.EncryptedCID
	Thumbhash   []byte
}

func NewFileVersionThumbnail(imageType string, aspectRatio float64, cid *encoding.EncryptedCID, thumbhash []byte) *FileVersionThumbnail {
	return &FileVersionThumbnail{
		ImageType:   imageType,
		AspectRatio: aspectRatio,
		CID:         cid,
		Thumbhash:   thumbhash,
	}
}

func (fvt *FileVersionThumbnail) EncodeMsgpack(enc *msgpack.Encoder) error {
	data := map[int]interface{}{
		2: fvt.AspectRatio,
		3: fvt.CID.ToBytes(),
	}

	if fvt.ImageType != "" {
		data[1] = fvt.ImageType
	}

	if fvt.Thumbhash != nil {
		data[4] = fvt.Thumbhash
	}

	return enc.Encode(data)
}
func (fvt *FileVersionThumbnail) DecodeMsgpack(dec *msgpack.Decoder) error {
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
			err := dec.Decode(&fvt.ImageType)
			if err != nil {
				return err
			}
		case int8(2):
			err := dec.Decode(&fvt.AspectRatio)
			if err != nil {
				return err
			}
		case int8(3):
			val, err := dec.DecodeBytes()
			if err != nil {
				return err
			}
			fvt.CID, err = encoding.EncryptedCIDFromBytes(val)
			if err != nil {
				return err
			}

		case int8(4):
			err := dec.Decode(&fvt.Thumbhash)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (fvt *FileVersionThumbnail) Encode() map[int]interface{} {
	data := map[int]interface{}{
		2: fvt.AspectRatio,
		3: fvt.CID.ToBytes(),
	}

	if fvt.ImageType != "" {
		data[1] = fvt.ImageType
	}

	if fvt.Thumbhash != nil {
		data[4] = fvt.Thumbhash
	}

	return data
}
