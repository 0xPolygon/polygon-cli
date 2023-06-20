package argfuzz

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/google/gofuzz"
)

type Hex64String string

func MutateRPCArgs(args *[]interface{}, c fuzz.Continue) {
	// TODO: think of better mutators here
	// some examples would be a valid hex, invalid hex, passing integers where
	// expecting string, etc.
	for i, d := range *args {
		switch d.(type) {
		case string:
			// random string vals
			(*args)[i] = c.RandString()
		case int:
			// todo: need to think what to randmoize
			(*args)[i] = c.Intn(100)
		case bool:
			// random bool val
			(*args)[i] = c.RandBool()
		case Hex64String:
			// generate a valid mutation or for invalid fuzz
			val := (*args)[i].(Hex64String)
			(*args)[i] = val.Mutate()
		default:
			// no mutation for unknown types
			fmt.Println("no match")
		}
	}
}

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (val Hex64String) Mutate() string {
	randHex, err := RandomHex(32)
	if err != nil {
		fmt.Println("err")
	}
	return "0x" + randHex
}
