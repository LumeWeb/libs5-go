package protocol

import "crypto/rand"

func GenerateChallenge() []byte {
	challenge := make([]byte, 32)
	_, err := rand.Read(challenge)
	if err != nil {
		panic(err)
	}

	return challenge
}
