package storage

import (
	ed25519p "crypto/ed25519"
	"git.lumeweb.com/LumeWeb/libs5-go/ed25519"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
)

func PrepareProvideMessage(identity *ed25519.KeyPairEd25519, hash *encoding.Multihash, location StorageLocation) []byte {
	// Initialize the list with the record type.
	list := []byte{byte(types.RecordTypeStorageLocation)}

	// Append the full bytes of the hash.
	list = append(list, hash.FullBytes()...)

	// Append the location type.
	list = append(list, byte(location.Type()))

	// Append the expiry time of the location, encoded as 4 bytes.
	list = append(list, utils.EncodeEndian(uint64(location.Expiry()), 4)...)

	// Append the number of parts in the location.
	list = append(list, byte(len(location.Parts())))

	// Iterate over each part in the location.
	for _, part := range location.Parts() {
		// Convert part to bytes.
		bytes := []byte(part)

		// Encode the length of the part as 4 bytes and append.
		list = append(list, utils.EncodeEndian(uint64(len(bytes)), 2)...)

		// Append the actual part bytes.
		list = append(list, bytes...)
	}

	// Append a null byte at the end of the list.
	list = append(list, 0)

	// Sign the list using the node's private key.
	signature := ed25519p.Sign(identity.ExtractBytes(), list)

	// Append the public key and signature to the list.
	finalList := append(list, identity.PublicKey()...)
	finalList = append(finalList, signature...)

	// Return the final byte slice.
	return finalList
}
