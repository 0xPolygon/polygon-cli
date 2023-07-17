// Package argfuzz implements the randomizers, mutators, and fuzzers
// that can be utilized by gofuzz. It extends on gofuzz default mutation with functionality for manipulating and
// mutating the byte slices of the desired fuzzed inputs
package argfuzz

import (
	"encoding/hex"
	"math/rand"
	"reflect"
	"strconv"

	fuzz "github.com/google/gofuzz"
	"github.com/rs/zerolog/log"
)

type FuzzFunc func(arg []byte) []byte

var fuzzFunctions = []FuzzFunc{
	AddRandomCharactersToStart,
	AddRandomCharactersToMiddle,
	AddRandomCharactersToEnd,
	DeleteRandomCharactersInStart,
	DeleteRandomCharactersInMiddle,
	DeleteRandomCharactersInEnd,
	DuplicateBytes,
	AddSpacesRandomly,
	ReplaceWithRandomCharacters,
	func(arg []byte) []byte {
		return nil
	},
}

var randSrc *rand.Rand

func SetSeed(seed *int64) {
	randSrc = rand.New(rand.NewSource(*seed))
}

func FuzzExecutor(arg []byte) []byte {
	selectedFunc := fuzzFunctions[rand.Intn(len(fuzzFunctions))]
	return selectedFunc(arg)
}

func ByteMutator(arg []byte) []byte {
	arg = FuzzExecutor(arg)
	// fitty-fitty chance of more mutations
	if rand.Intn(2) == 0 {
		return ByteMutator(arg)
	}

	return arg
}

func FuzzRPCArgs(args *[]interface{}, c fuzz.Continue) {
	for i, d := range *args {
		if d == nil {
			d = c.RandString()
		}

		switch dataType := reflect.TypeOf(d).Kind(); dataType {
		case reflect.String:
			(*args)[i] = string(ByteMutator([]byte(d.(string))))
		case reflect.Int:
			numString := strconv.Itoa(d.(int))

			byteArrString := string(ByteMutator([]byte(numString)))
			num, err := strconv.Atoi(byteArrString)
			if err != nil {
				(*args)[i] = byteArrString
			} else {
				(*args)[i] = num
			}
		case reflect.Bool:
			(*args)[i] = c.RandBool()
		case reflect.Float32:
			(*args)[i] = c.Float32()
		case reflect.Float64:
			(*args)[i] = c.Float64()
		case reflect.Ptr:
			c.Fuzz(d)
			(*args)[i] = d
		default:
			(*args)[i] = c.RandString()
		}
	}
}

func RandomByte() byte {
	return byte(rand.Intn(256))
}

func RandomBytesSize(size int) []byte {
	bytes := make([]byte, size)

	_, err := randSrc.Read(bytes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate random bytes from default Source.")
		return []byte{RandomByte()}
	}

	return bytes
}

func RandomBytes() []byte {
	length := rand.Intn(100) + 1 // for the time being, generates 1-100 bytes
	return RandomBytesSize(length)
}

func RandomHexBytes() []byte {
	length := rand.Intn(100) + 1 // for the time being, generates 1-100 bytes
	bytes := RandomBytesSize(length)
	return []byte(hex.EncodeToString(bytes))
}

func RandomBytesOfCharactersOrHexBytes() []byte {
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
	bytes := RandomBytesOfCharactersOrHexBytes()
	return append(bytes, arg...)
}

func AddRandomCharactersToMiddle(arg []byte) []byte {
	n := len(arg)
	bytes := RandomBytesOfCharactersOrHexBytes()

	if n == 0 {
		return bytes
	}

	randomMid := rand.Intn(len(arg))
	return append(arg[:randomMid], append(bytes, arg[randomMid:]...)...)
}

func AddRandomCharactersToEnd(arg []byte) []byte {
	bytes := RandomBytesOfCharactersOrHexBytes()
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

	whitespaceByte := RandomWhitespace()

	if n == 0 {
		return []byte{whitespaceByte}
	}

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

func ReplaceWithRandomCharacters(arg []byte) []byte {
	n := len(arg)
	return RandomBytesSize(n)
}
