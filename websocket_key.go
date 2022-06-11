package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func main() {
	clientKeySrc := "the sample nonce"
	key := base64.StdEncoding.EncodeToString([]byte(clientKeySrc))
	fmt.Printf("Sec-WebSocket-Key: %s\n", key)

	salt := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.Sum([]byte(key + salt))
	acceptKey := base64.StdEncoding.EncodeToString(hash[:])
	fmt.Printf("Sec-WebSocket-Accept: %s\n", acceptKey)
}
