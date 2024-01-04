package serialize

import (
	"encoding/base64"
	"github.com/multiformats/go-multibase"
)

func UnmarshalBase64UrlJSON(data []byte) ([]byte, error) {
	strData := string(data)
	if len(strData) >= 2 && strData[0] == '"' && strData[len(strData)-1] == '"' {
		strData = strData[1 : len(strData)-1]
	}

	if strData == "null" {
		return nil, nil
	}

	if strData[0] == 'u' {
		_, decoded, err := multibase.Decode(strData)
		if err != nil {
			return nil, err
		}
		return decoded, nil
	}

	decodedData, err := base64.RawURLEncoding.DecodeString(strData)
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}
