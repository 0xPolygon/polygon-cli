package argfuzz

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"math/rand"
	"strings"

	"github.com/google/gofuzz"
)

type (
	// probably don't need this for fuzzing
	BaseType interface {
		Mutate() string
	}
	Hex64String string
	Hex42String string
	Hex32String string
)

func RandomStringMutator(arg string, c fuzz.Continue) string {
	switch rand.Intn(10) {
	case 0:
		return EmptyString()
	case 1:
		return c.RandString()
	case 2:
		return AddRandomCharacter(arg)
	case 3:
		return AddCharacters(arg, c.RandString())
	case 4:
		return RemoveRandomCharacter(arg)
	case 5:
		return GenerateRandomHexStringFixedSize(len(arg))
	case 6:
		return GenerateRandomHexStringRandomSize()
	case 7:
		return GenerateInvalidHexString(len(arg))
	case 8:
		return TestWhitespaceControlChars(arg)
	case 9:
		return TestCaseSensitivity(arg)
	default:
		return c.RandString()
	}
}

func MutateStringArgs(args *string, c fuzz.Continue) {
	*args = RandomStringMutator(*args, c)
}

func MutateRPCArgs(args *[]interface{}, c fuzz.Continue) {
	for i, d := range *args {
		switch d.(type) {
		case string:
			(*args)[i] = RandomStringMutator((*args)[i].(string), c)
		case int:
			// use large number
			// use float?
			(*args)[i] = c.Intn(100)
		case bool:
			// random bool val
			// use string bool value
			// use 0/1
			(*args)[i] = c.RandBool()
		case BaseType:
			// Mutate to valid type
			val := (*args)[i].(BaseType)
			(*args)[i] = val.Mutate()
		default:
			// No mutation for unknown types
			log.Error().Msg("Unable to retreive type for fuzzing")
		}
	}
}

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := cryptoRand.Read(bytes); err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(bytes), nil
}

func (val Hex64String) Mutate() string {
	return GenerateRandomHexStringFixedSize(32)
}

func (val Hex42String) Mutate() string {
	return GenerateRandomHexStringFixedSize(21)
}

func (val Hex32String) Mutate() string {
	return GenerateRandomHexStringFixedSize(16)
}

func EmptyString() string {
	return ""
}

func GenerateRandomString(input string) string {
	f := fuzz.New()
	var fuzzedString string
	f.Fuzz(&fuzzedString)
	return fuzzedString
}

func AddRandomCharacter(input string) string {
	index := rand.Intn(len(input) + 1)
	char := string(rune(rand.Intn(128))) // a random ASCII val added
	return input[:index] + char + input[index:]
}

func AddCharacters(input string, characters string) string {
	index := rand.Intn(len(input) + 1)
	return input[:index] + characters + input[index:]
}

func RemoveRandomCharacter(input string) string {
	if len(input) == 0 {
		return input
	}
	index := rand.Intn(len(input))
	return input[:index] + input[index+1:]
}

func IsHexString(input string) bool {
	_, err := ParseHexString(input)
	return err == nil
}
func ParseHexString(input string) ([]byte, error) {
	return hex.DecodeString(input)
}

func GenerateRandomHexStringFixedSize(length int) string {
	buf := make([]byte, length)
	rand.Read(buf)
	return "0x" + hex.EncodeToString(buf)
}

func GenerateRandomHexStringRandomSize() string {
	buf := make([]byte, rand.Intn(100)+1) // length 1 - 100
	rand.Read(buf)
	return "0x" + hex.EncodeToString(buf)
}

func GenerateInvalidHexString(length int) string {
	invalidChars := []rune{'G', 'Z', '#', '$', '%', '@'}
	result := make([]rune, length)
	for i := 0; i < length; i++ {
		if rand.Intn(5) == 0 {
			result[i] = invalidChars[rand.Intn(len(invalidChars))]
		} else {
			result[i] = rune(rand.Intn(16)) + '0'
		}
	}
	return string(result)
}

func TestWhitespaceControlChars(input string) string {
	whitespaceChars := []rune{' ', '\t', '\n'}
	controlChars := []rune{'\x00', '\x01', '\x02'}

	// Add leading whitespace
	if len(input) > 0 && rand.Intn(2) == 0 {
		input = string(whitespaceChars[rand.Intn(len(whitespaceChars))]) + input
	}

	// Add trailing whitespace
	if len(input) > 0 && rand.Intn(2) == 0 {
		input += string(whitespaceChars[rand.Intn(len(whitespaceChars))])
	}

	// Add internal whitespace
	if len(input) > 1 && rand.Intn(2) == 0 {
		index := rand.Intn(len(input)-1) + 1
		input = input[:index] + string(whitespaceChars[rand.Intn(len(whitespaceChars))]) + input[index:]
	}

	// Add control characters
	numControlChars := rand.Intn(3) // Random number of control characters (0-2)
	for i := 0; i < numControlChars; i++ {
		index := rand.Intn(len(input) + 1)
		input = input[:index] + string(controlChars[rand.Intn(len(controlChars))]) + input[index:]
	}

	return input
}

func TestCaseSensitivity(input string) string {
	if rand.Intn(2) == 0 {
		return strings.ToUpper(input)
	}
	return strings.ToLower(input)
}
