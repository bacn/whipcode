//
//  Copyright 2024 whipcode.app (AnnikaV9)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on
//  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific
//  language governing permissions and limitations under the License.
//

package control

import (
	"encoding/hex"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
)

func (ks *KeyStore) CheckKey(userKey string, serverKey []string) bool {
	ks.once.Do(func() {
		hash := argon2.IDKey([]byte(userKey), []byte(serverKey[1]), 1, 4096, 1, 32)
		if hex.EncodeToString(hash) == strings.TrimSpace(serverKey[0]) {
			ks.validKey = userKey
		}
	})

	return userKey == ks.validKey
}

func InitializeKeystore(keyFile string) (*KeyStore, []string) {
	file, errH := os.ReadFile(keyFile)
	if errH != nil {
		log.Fatalf("Fatal: Could not read %s: %s", keyFile, errH)
	}

	keyAndSalt := strings.Split(string(file), "\n")
	if len(keyAndSalt) != 2 || len(keyAndSalt[1]) < 1 {
		log.Fatal("Fatal: Key file format invalid")
	}
	keyStore := KeyStore{}

	return &keyStore, keyAndSalt
}
