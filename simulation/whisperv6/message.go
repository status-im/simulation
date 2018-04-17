package whisperv6

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/whisper/whisperv6"
)

// const from github.com/ethereum/go-ethereum/whisper/whisperv5/doc.go
const (
	aesKeyLength = 32
)

func generateMessage(ttl int, symkeyID string) *whisperv6.NewMessage {
	// set all the parameters except p.Dst and p.Padding
	buf := make([]byte, 4)
	rand.Read(buf)
	sz := rand.Intn(400)

	msg := &whisperv6.NewMessage{
		PowTarget: 0.01,
		PowTime:   1,
		Payload:   make([]byte, sz),
		SymKeyID:  symkeyID,
		Topic:     whisperv6.BytesToTopic(buf),
		TTL:       uint32(ttl),
	}
	rand.Read(msg.Payload)

	return msg
}
