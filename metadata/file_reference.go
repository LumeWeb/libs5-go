package metadata

import "github.com/vmihailenco/msgpack/v5"

var _ SerializableMetadata = (*FileReference)(nil)

type FileReference struct {
	Name     string                 `json:"name"`
	Created  int                    `json:"created"`
	Version  int                    `json:"version"`
	File     *FileVersion           `json:"file"`
	Ext      map[string]interface{} `json:"ext"`
	History  map[int]*FileVersion   `json:"history"`
	MimeType string                 `json:"mimeType"`
	URI      string                 `json:"uri"`
	Key      string                 `json:"key"`
}

func NewFileReference(name string, created, version int, file *FileVersion, ext map[string]interface{}, history map[int]*FileVersion, mimeType string) *FileReference {
	return &FileReference{
		Name:     name,
		Created:  created,
		Version:  version,
		File:     file,
		Ext:      ext,
		History:  history,
		MimeType: mimeType,
		URI:      "",
		Key:      "",
	}
}

func (fr *FileReference) Modified() int {
	return fr.File.Ts
}

func (fr *FileReference) EncodeMsgpack(enc *msgpack.Encoder) error {
	data := map[int]interface{}{
		1: fr.Name,
		2: fr.Created,
		4: fr.File,
		5: fr.Version,
	}

	if fr.MimeType != "" {
		data[6] = fr.MimeType
	}

	if fr.Ext != nil {
		data[7] = fr.Ext
	}

	if fr.History != nil {
		historyData := make(map[int]interface{})
		for key, value := range fr.History {
			historyData[key] = value
		}
		data[8] = historyData
	}

	return enc.Encode(data)
}
func (fr *FileReference) DecodeMsgpack(dec *msgpack.Decoder) error {
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
			err := dec.Decode(&fr.Name)
			if err != nil {
				return err
			}
		case int8(2):
			err := dec.Decode(&fr.Created)
			if err != nil {
				return err
			}
		case int8(4):
			err := dec.Decode(&fr.File)
			if err != nil {
				return err
			}
		case int8(5):
			val, err := dec.DecodeInt()
			if err != nil {
				return err
			}

			fr.Version = val
		case int8(6):
			err := dec.Decode(&fr.MimeType)
			if err != nil {
				return err
			}
		case int8(7):
			err := dec.Decode(&fr.Ext)
			if err != nil {
				return err
			}
		case int8(8):
			historyDataLen, err := dec.DecodeMapLen()
			if err != nil {
				return err
			}
			fr.History = make(map[int]*FileVersion, historyDataLen)
			for range fr.History {
				k, err := dec.DecodeInt()
				if err != nil {
					return err
				}

				var fileVersion FileVersion
				err = dec.Decode(&fileVersion)
				if err != nil {
					return err
				}

				fr.History[k] = &fileVersion
			}
		}
	}
	return nil
}
