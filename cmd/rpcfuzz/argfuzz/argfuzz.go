package argfuzz

import (
	"encoding/hex"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/google/gofuzz"
)

func MutateExecutor(arg []byte, c fuzz.Continue) []byte {
	switch rand.Intn(10) {
	case 0:
		return AddRandomCharactersToStart(arg)
	case 1:
		return AddRandomCharactersToMiddle(arg)
	case 2:
		return AddRandomCharactersToEnd(arg)
	case 3:
		return DeleteRandomCharactersInStart(arg)
	case 4:
		return DeleteRandomCharactersInMiddle(arg)
	case 5:
		return DeleteRandomCharactersInEnd(arg)
	case 6:
		return DuplicateBytes(arg)
	case 7:
		return AddSpacesRandomly(arg)
	case 8:
		return []byte(c.RandString())
	case 9:
		return nil
	default:
		// Do nothing
	}

	return arg
}

func ByteMutator(arg []byte, c fuzz.Continue) []byte {
	arg = MutateExecutor(arg, c)
	// fitty-fitty chance of more mutations
	if rand.Intn(2) == 0 && arg != nil {
		return ByteMutator(arg, c)
	}

	return arg
}

func MutateRPCArgs(args *[]interface{}, c fuzz.Continue) {
	for i, d := range *args {
		switch d.(type) {
		case string:
			(*args)[i] = string(ByteMutator([]byte(d.(string)), c))
		case int:
			numString := strconv.Itoa(d.(int))

			byteArrString := string(ByteMutator([]byte(numString), c))
			num, err := strconv.Atoi(byteArrString)
			if err != nil {
				(*args)[i] = byteArrString
			} else {
				(*args)[i] = num
			}
		case bool:
			(*args)[i] = c.RandBool()
		default:
			if reflect.TypeOf(d).Kind() == reflect.Ptr {
				c.Fuzz(d)
				(*args)[i] = d
			} else {
				(*args)[i] = c.RandString()
			}
		}
	}
}

func RandomByte() byte {
	return byte(rand.Intn(256))
}

func RandomHexByte() byte {
	hex := []byte("0123456789abcdefABCDEF")
	return byte(hex[rand.Intn(len(hex))])
}

func RandomBytes() ([]byte, error) {
	length := rand.Intn(100) + 1 // for the time being, generates 1-100 characters
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return []byte{RandomByte()}, err
	}

	return bytes, nil
}

func RandomHexBytes() ([]byte, error) {
	length := rand.Intn(100) + 1 // for the time being, generates 1-100 characters
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return []byte{RandomHexByte()}, err
	}
	resBytes := []byte(hex.EncodeToString(bytes)[:length])

	return resBytes, nil
}

func RandomBytesOfCharactersOrHexBytes() ([]byte, error) {
	if rand.Intn(2) == 0 {
		return RandomBytes()
	}

	return RandomHexBytes()
}

func RandomWhitespace() byte {
	whitespaces := []byte{' ', '\n', '\t'}

	whitespaceIdx := rand.Intn(len(whitespaces))

	return whitespaces[whitespaceIdx]
}

func AddRandomCharactersToStart(arg []byte) []byte {
	bytes, err := RandomBytesOfCharactersOrHexBytes()
	if err != nil {
		return arg
	}
	return append(bytes, arg...)
}

func AddRandomCharactersToMiddle(arg []byte) []byte {
	n := len(arg)
	if n == 0 {
		return arg
	}

	bytes, err := RandomBytesOfCharactersOrHexBytes()
	if err != nil {
		return arg
	}
	randomMid := rand.Intn(len(arg))
	return append(arg[:randomMid], append(bytes, arg[randomMid:]...)...)
}

func AddRandomCharactersToEnd(arg []byte) []byte {
	bytes, err := RandomBytesOfCharactersOrHexBytes()
	if err != nil {
		return arg
	}
	return append(arg, bytes...)
}

func DeleteRandomCharactersInStart(arg []byte) []byte {
	n := len(arg)
	if n == 0 {
		return arg
	}

	numberBytesToDelete := rand.Intn(n)
	return arg[numberBytesToDelete:]
}

func DeleteRandomCharactersInMiddle(arg []byte) []byte {
	n := len(arg)
	if n == 0 {
		return arg
	}

	numberBytesStart := rand.Intn(n)
	numberBytesEnd := numberBytesStart + rand.Intn(n-numberBytesStart)
	return append(arg[:numberBytesStart], arg[numberBytesEnd:]...)
}

func DeleteRandomCharactersInEnd(arg []byte) []byte {
	n := len(arg)
	if n == 0 {
		return arg
	}

	numberBytesToDelete := rand.Intn(n)
	return arg[:numberBytesToDelete]
}

func DuplicateBytes(arg []byte) []byte {
	numOfDuplicates := rand.Intn(10) + 1 // for now duplicate at most 10 times
	duplicatedArgs := arg
	for i := 0; i < numOfDuplicates; i++ {
		duplicatedArgs = append(duplicatedArgs, arg...)
	}

	return duplicatedArgs
}

func AddSpacesRandomly(arg []byte) []byte {
	n := len(arg)
	if n == 0 {
		return arg
	}

	whitespaceByte := RandomWhitespace()
	randomMid := rand.Intn(len(arg))
	funcs := []func() []byte{
		func() []byte { return append([]byte{whitespaceByte}, arg...) },
		func() []byte {
			return append(arg[:randomMid], append([]byte{whitespaceByte}, arg[randomMid:]...)...)
		},
		func() []byte { return append(arg, whitespaceByte) },
	}

	return funcs[rand.Intn(len(funcs))]()
}
