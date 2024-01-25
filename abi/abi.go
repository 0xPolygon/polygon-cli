// Package abi provides utilities for interacting with contracts.
// There are 2 primary functionalities:
// - encoding
// - decoding
// All specifications of encoding and decoding adheres to Solidity's ABI specs:
// https://docs.soliditylang.org/en/latest/abi-spec.html
package abi

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	// "github.com/alecthomas/repr"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EncodedItem is a data structured used to represent the differences between static and dynamic types.
// Static types are encoded in-place (as a head) and
// dynamic types are encoded at a separately allocated location after the current block (i.e., tail).
// See more here: https://docs.soliditylang.org/en/latest/abi-spec.html#formal-specification-of-the-encoding
type EncodedItem struct {
	Head string
	Tail string
}

// This is the place holder length of a dynamic type, which is 64 string characters (32 bytes).
// This place holder just points to the location of the actual value of the dynamic type.
const PlaceholderPointerLength = 64

// FunctionSignature Parser
// FunctionSignature represents the overall structure of a function signature.
type FunctionSignature struct {
	FunctionName string             `parser:"@Ident"`
	FunctionArgs []*FunctionArgType `parser:"'(' @@* ( ',' @@ )* ')'"`
}

// FunctionArgType represents a single argument which can be a base type, a tuple, or an array.
// At a high level, this lexer parses all types as (@Ident | @@)@Brackets, where:
// - @Ident is a base type
// - @@ is a recursive call to represent each item of the tuple
// - @Bracket indicates the nested level of the array
type FunctionArgType struct {
	Type  string             `parser:"(@Ident"`
	Tuple *FunctionTupleType `parser:"| @@)"`
	Array []string           `parser:"@Bracket*"`
}

// FunctionTupleType represents a tuple in the signature.
type FunctionTupleType struct {
	Elements []*FunctionArgType `parser:"('(' @@ ( ',' @@ )* ')')"`
}

var (
	funcSigLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[a-zA-Z_][a-zA-Z_\d]*`},
		{Name: `Bracket`, Pattern: `\[\d*\]*`},
		{Name: `Punct`, Pattern: `[(),\[\]]`},
		{Name: `whitespace`, Pattern: `\s+`},
	})
	FunctionSignatureParser = participle.MustBuild[FunctionSignature](
		participle.Lexer(funcSigLexer),
	)
)

// Function Argument Input Object Parser
// Object represents one Function argument object
// This adheres to the available inputs that Solidity functions can accept such as:
type Object struct {
	Val       string `parser:"@Ident"`
	Stringval string `parser:"| @String"` // NOTE: Strictly string values must be around quotations: "
	Tuple     Tuple  `parser:"| @@"`
	Array     Array  `parser:"| '[' @@* ']'"`
}
type Tuple struct {
	Elements []Object `parser:"('(' @@ (',' @@)* ')')"`
}
type Array struct {
	Elements []Object `parser:"@@ (',' @@)*"`
}

var (
	objectLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[a-zA-Z_\d-]+`},
		{Name: `String`, Pattern: `"([^"\\]|\\.)*"`},
		{Name: `Punct`, Pattern: `[(),\[\]]`},
		{Name: `whitespace`, Pattern: `\s+`},
	})
	ObjectParser = participle.MustBuild[Object](
		participle.Lexer(objectLexer),
	)
)

// Encode returns the string encoding of an Object
// Example:
// functionSignature := `getMultipliedAndAddNumber((uint256,bool,string[]))`
// functionSignatureInputs := []string{`(200, true, ["a","b"])`}
//
// This is the representation of the function signature:
//
//	FunctionSignature{
//	  FunctionName: "getMultipliedAndAddNumber",
//	  FunctionArgs: []*main.FunctionArgType{
//	    {
//	      Tuple: &main.FunctionTupleType{
//	        Elements: []*main.FunctionArgType{
//	          {
//	            Type: "uint256",
//	          },
//	          {
//	            Type: "bool",
//	          },
//	          {
//	            Type: "string",
//	            Array: []string{
//	              "[]",
//	            },
//	          },
//	        },
//	      },
//	    },
//	  },
//	}
//
// This is the object input type representation of a function argument input
//
//	Object{
//	  Tuple: main.Tuple{
//	    Elements: []main.Object{
//	      {
//	        Val: "200",
//	      },
//	      {
//	        Val: "true",
//	      },
//	      {
//	        Array: main.Array{
//	          Elements: []main.Object{
//	            {
//	              Stringval: "\"a\"",
//	            },
//	            {
//	              Stringval: "\"b\"",
//	            },
//	          },
//	        },
//	      },
//	    },
//	  },
//	}
//
// From there, we traverse the function signature AST and the parsed object to encode each item, respectively.
func (fs *FunctionSignature) Encode(functionArguments []string) (string, error) {
	if len(fs.FunctionArgs) != len(functionArguments) {
		return "", fmt.Errorf("# of arguments doesn't match")
	}

	vals := make([]EncodedItem, 0)
	tailLoc := 0 // idx ptr to beginning of tail location

	for idx, functionArg := range fs.FunctionArgs {
		if functionArg.Type == "string" && len(functionArg.Array) == 0 {
			// this is a little hacky but when it's a pure string type (not a string array), make sure the function argument is surrounded
			// by quotations as that's how the lexer identifies a string.
			// a string in an array or tuple would already be enclosed by quotation marks.
			functionArguments[idx] = ValidateStringIsQuoted(functionArguments[idx])
		}
		object, err := ObjectParser.ParseString("", functionArguments[idx])
		if err != nil {
			return "", err
		}

		encodedString, err := functionArg.EncodeInput(*object)
		if err != nil {
			return "", err
		}
		if encodedString.Tail != "" {
			// is dynamic so head is none
			tailLoc += PlaceholderPointerLength
		} else {
			tailLoc += len(encodedString.Head)
		}
		vals = append(vals, encodedString)
	}

	// backfill dynamic types...
	head := ""
	tail := ""
	for _, val := range vals {
		currHead := val.Head
		if val.Tail != "" {
			// i.e., dynamic
			pointerLoc := (tailLoc + len(tail)) / 2 // convert the hex length to bytes
			pointerLocHex, err := ConvertInt(fmt.Sprintf("%d", pointerLoc))
			if err != nil {
				return "", err
			}
			currHead = pointerLocHex
			tail += val.Tail
		}
		head += currHead
	}

	return head + tail, nil
}

// EncodeInput encodes an object with the mapped function arg type.
func (fat FunctionArgType) EncodeInput(object Object) (EncodedItem, error) {
	vals := make([]EncodedItem, 0)
	tailLoc := 0 // idx ptr to beginning of tail location
	var convertedVal string
	var conversionErr error

	switch {
	case len(fat.Array) > 0:
		// Array must go first since a representation of it can evaluate to true on more than one of these instances.
		// An array representation of `string[][3]`` would be something like:
		//	{
		//	  Type: "string",
		//	  Array: []string{
		//	    "[]",
		//	    "[3]",
		//	  },
		//	},
		// numOfNestedLevel = count the nested level, i.e. len of fat.Array
		// while numOfNestedLevel > 0: numOfNestedLevel-- and recurse
		// when numOfNestedLevel == 0: iterate through each Array.Elements and do regular conversion when reach base type
		subFat := fat
		subFat.Array = subFat.Array[1:] // this is the scenario of numOfNestedLevel > 0: numOfNestedLevel-- and recurse

		lenOfArray := len(object.Array.Elements)
		lenOfArrayHex, err := ConvertInt(fmt.Sprintf("%d", lenOfArray))
		if err != nil {
			return EncodedItem{}, err
		}
		vals = append(vals, EncodedItem{Head: lenOfArrayHex})

		// TODO: For a fixed size array, type[M], validate the number of len(object.Array.Elements) == "M"
		for _, itemObject := range object.Array.Elements {
			encodedInput, err := subFat.EncodeInput(itemObject)
			if err != nil {
				return EncodedItem{}, err
			}

			if encodedInput.Tail != "" {
				// a dynamic input, so the head is none
				tailLoc += PlaceholderPointerLength
			} else {
				tailLoc += len(encodedInput.Head)
			}

			vals = append(vals, encodedInput)
		}
	case fat.Tuple != nil:
		// Tuple is less complicated than arrays since we just iterate through each items
		if len(object.Tuple.Elements) != len(fat.Tuple.Elements) {
			return EncodedItem{}, fmt.Errorf("Mismatched length of tuple elements. Expected: %d elements, received %d", len(fat.Tuple.Elements), len(object.Tuple.Elements))
		}

		for idx, tupleItemObject := range object.Tuple.Elements {
			encodedInput, err := fat.Tuple.Elements[idx].EncodeInput(tupleItemObject)
			if err != nil {
				return EncodedItem{}, err
			}

			if encodedInput.Tail != "" {
				// a dynamic input, so the head is none
				tailLoc += PlaceholderPointerLength
			} else {
				tailLoc += len(encodedInput.Head)
			}

			vals = append(vals, encodedInput)
		}
	case strings.Contains(fat.Type, "string"):
		stringVal := object.Val
		if object.Stringval != "" {
			stringVal = strings.Trim(object.Stringval, `"`) // TODO: maybe move this logic (the grabbed surrounding quotations) into the lexer?
		}
		convertedVal, conversionErr = ConvertString(stringVal)
		if conversionErr != nil {
			return EncodedItem{}, fmt.Errorf("Failed to convert %s to an %s type. %v", object.Val, fat.Type, conversionErr)
		}
		return EncodedItem{Tail: convertedVal}, nil
	case fat.Type == "bytes":
		convertedVal, conversionErr = ConvertBytes(object.Val)
		if conversionErr != nil {
			return EncodedItem{}, fmt.Errorf("Failed to convert %s to an %s type. %v", object.Val, fat.Type, conversionErr)
		}
		return EncodedItem{Tail: convertedVal}, nil
	case strings.HasPrefix(fat.Type, "int"):
		// TODO: validate input is within int<size> limit
		convertedVal, conversionErr = ConvertInt(object.Val)
		if conversionErr != nil {
			return EncodedItem{}, fmt.Errorf("Failed to convert %s to an %s type. %v", object.Val, fat.Type, conversionErr)
		}
		return EncodedItem{Head: convertedVal}, nil
	case strings.HasPrefix(fat.Type, "uint"):
		// TODO: validate input is within uint<size> limit
		convertedVal, conversionErr = ConvertUint(object.Val)
		if conversionErr != nil {
			return EncodedItem{}, fmt.Errorf("Failed to convert %s to an %s type. %v", object.Val, fat.Type, conversionErr)
		}
		return EncodedItem{Head: convertedVal}, nil
	case strings.Contains(fat.Type, "bool"):
		convertedVal, conversionErr = ConvertBool(object.Val)
		if conversionErr != nil {
			return EncodedItem{}, fmt.Errorf("Failed to convert %s to an %s type. %v", object.Val, fat.Type, conversionErr)
		}
		return EncodedItem{Head: convertedVal}, nil
	case strings.HasPrefix(fat.Type, "bytes"):
		convertedVal, conversionErr = ConvertByteSize(object.Val, fat.Type)
		if conversionErr != nil {
			return EncodedItem{}, fmt.Errorf("Failed to convert %s to an %s type. %v", object.Val, fat.Type, conversionErr)
		}
		return EncodedItem{Head: convertedVal}, nil
	case fat.Type == "address":
		convertedVal, conversionErr = ConvertAddress(object.Val)
		if conversionErr != nil {
			return EncodedItem{}, fmt.Errorf("Failed to convert %s to an %s type. %v", object.Val, fat.Type, conversionErr)
		}
		return EncodedItem{Head: convertedVal}, nil
	default:
		return EncodedItem{}, fmt.Errorf("Invalid type %s", fat.Type)
	}

	// backfill dynamic types...
	head := ""
	tail := ""
	if len(fat.Array) > 0 || fat.Tuple != nil {
		for _, val := range vals {
			currHead := val.Head
			if val.Tail != "" {
				// is dynamic
				pointerLoc := (tailLoc + len(tail)) / 2 // converrt the hex length to bytes
				pointerLocHex, err := ConvertInt(fmt.Sprintf("%d", pointerLoc))
				if err != nil {
					return EncodedItem{}, err
				}
				currHead = pointerLocHex
				tail += val.Tail
			}
			head += currHead
		}
	}

	// made it here so, it's either tuple or array. so we must take into account the following non-dynamic (static) edge cases for tuple or array:
	// - (T1,...,Tk) if Ti is static for all 1 <= i <= k
	// - T[k] for any static T and any k >= 0
	if fat.IsStaticType() {
		return EncodedItem{Head: head + tail}, nil
	}

	return EncodedItem{Tail: head + tail}, nil
}

// IsStaticType returns true if the function arg type is static.
// According to [solidity docs](https://docs.soliditylang.org/en/latest/abi-spec.html#formal-specification-of-the-encoding),
// ```
// Definition: The following types are called “dynamic”:
// - bytes
// - string
// - T[] for any T
// - T[k] for any dynamic T and any k >= 0
// - (T1,...,Tk) if Ti is dynamic for some 1 <= i <= k
//
// All other types are called “static”.
// ```
//
// so these are static:
// - uint<M>: unsigned integer type of M bits, 0 < M <= 256, M % 8 == 0. e.g. uint32, uint8, uint256.
// - int<M>: two’s complement signed integer type of M bits, 0 < M <= 256, M % 8 == 0
// - bool: equivalent to uint8 restricted to the values 0 and 1
// - address
// - bytes<M>: binary type of M bytes, 0 < M <= 32
// - (T1,...,Tk) if Ti is static for all 1 <= i <= k
// - T[k] for any static T and any k >= 0
func (fat FunctionArgType) IsStaticType() bool {
	switch {
	case len(fat.Array) > 0:
		// - T[k] for any static T and any k >= 0
		arrayDimensions := len(fat.Array)
		if fat.Array[arrayDimensions-1] == "[]" {
			return false
		}
		subFat := fat
		subFat.Array = subFat.Array[:arrayDimensions-1]
		return subFat.IsStaticType()
	case fat.Tuple != nil:
		// - (T1,...,Tk) if Ti is static for all 1 <= i <= k
		isStatic := true
		for _, tupleTypeObject := range fat.Tuple.Elements {
			isStatic = isStatic && tupleTypeObject.IsStaticType()
		}
		return isStatic
	case strings.Contains(fat.Type, "string"):
		return false
	case fat.Type == "bytes":
		return false
	case strings.HasPrefix(fat.Type, "int"):
		return true
	case strings.HasPrefix(fat.Type, "uint"):
		return true
	case strings.Contains(fat.Type, "bool"):
		return true
	case strings.HasPrefix(fat.Type, "bytes"):
		return true
	case fat.Type == "address":
		return true
	}

	return false
}

// leftPadWithZeros fills a string to the left with 0s until targetLen
func leftPadWithZeros(value string, targetLen int) string {
	strLen := len(value)
	return strings.Repeat("0", targetLen-strLen) + value
}

// rightPadWithZeros fills a string to the right with 0s until targetLen
func rightPadWithZeros(value string, targetLen int) string {
	strLen := len(value)
	return value + strings.Repeat("0", targetLen-strLen)
}

// ConvertInt converts an int input to the 32 byte hex encoding, left padded with 0s
func ConvertInt(value string) (string, error) {
	if len(value) < 1 {
		return "", fmt.Errorf("Error: expected at least one digit")
	}
	bigInt := new(big.Int)

	_, ok := bigInt.SetString(value, 10)
	if !ok {
		return "", fmt.Errorf("Error: Invalid integer string. Failed to convert %s to big int", value)
	}

	var hexString string
	bytes := bigInt.Bytes()
	if bigInt.Sign() < 0 {
		// convert to two's complement bytes
		twosComplement := make([]byte, len(bytes))
		for i, b := range bytes {
			twosComplement[i] = ^b
		}

		bigInt.SetBytes(twosComplement)
		bigInt.Add(bigInt, big.NewInt(1))

		bytes = bigInt.Bytes()
		hexValue := fmt.Sprintf("%x", bytes)
		hexString = strings.Repeat("f", 64-len(hexValue)) + hexValue // left pad with "f" since it's the complement
	} else {
		// Convert to hexadecimal representation
		hexValue := bigInt.Text(16)
		hexString = leftPadWithZeros(hexValue, 64)
	}

	return hexString, nil
}

// ConvertUint converts a uint input to the 32 byte hex encoding, left padded with 0s
func ConvertUint(value string) (string, error) {
	if len(value) < 1 {
		return "", fmt.Errorf("Error: expected at least one digit")
	}
	if value[0] == '-' {
		return "", fmt.Errorf("Error: Invalid integer string. %s can't be negative", value)
	}

	bigInt := new(big.Int)

	_, ok := bigInt.SetString(value, 10)
	if !ok {
		return "", fmt.Errorf("Error: Invalid integer string. Failed to convert %s to big int", value)
	}

	// Convert to hexadecimal representation
	hexValue := bigInt.Text(16)
	hexString := leftPadWithZeros(hexValue, 64)

	return hexString, nil
}

// ConvertBool converts a string representation of "true" or "false" into a 32 byte hex encoding (0 or 1), left padded with 0s
func ConvertBool(value string) (string, error) {
	switch value {
	case "true":
		return "0000000000000000000000000000000000000000000000000000000000000001", nil
	case "false":
		return "0000000000000000000000000000000000000000000000000000000000000000", nil
	default:
		return "", fmt.Errorf("bool must be either 'true' or 'false', found %s", value)
	}
}

// ConvertString returns hex string length + string hex encoding
func ConvertString(value string) (string, error) {
	// TODO: Maybe trim excess quotations?
	valueSize := len(value)
	valueSizeHex, err := ConvertInt(fmt.Sprintf("%d", valueSize))
	if err != nil {
		return "", err
	}

	if value == "" {
		return valueSizeHex, nil
	}

	utf8Bytes := []byte(value)                  // convert string to UTF8 encoded bytes
	hexEncoded := hex.EncodeToString(utf8Bytes) // encode bytes as hex string

	// spaceRequired is pretty much the number of 32 bytes we need to store the string
	spaceRequired := ((len(hexEncoded) + 63) / 64) * 64 // get ceiling

	hexString := rightPadWithZeros(hexEncoded, spaceRequired)

	return valueSizeHex + hexString, nil
}

// ConvertBytes returns bytes length + bytes hex encoding
// NOTE: this is similar to ConvertString but `value` is expected to be a hex input already
func ConvertBytes(value string) (string, error) {
	if len(value)%2 != 0 {
		return "", fmt.Errorf("Odd number of digits")
	}

	valueSize := len(value) / 2 // it's hex so / 2 for actual length in bytes
	valueSizeHex, err := ConvertInt(fmt.Sprintf("%d", valueSize))
	if err != nil {
		return "", err
	}

	if value == "" {
		return valueSizeHex, nil
	}

	// spaceRequired is pretty much the number of 32 bytes we need to store the string
	spaceRequired := ((len(value) + 63) / 64) * 64 // get ceiling

	hexString := rightPadWithZeros(value, spaceRequired)

	return valueSizeHex + hexString, nil
}

// ConvertByteSize converts a hex value to bytes<size>
// Eg.
// "0x123456", bytes3
// returns 0x0123450000000000000000000000000000000000000000000000000000000000
func ConvertByteSize(value string, byteType string) (string, error) {
	byteSizeString := strings.TrimPrefix(byteType, "bytes")
	byteSize, err := strconv.Atoi(byteSizeString)
	if err != nil {
		return "", err
	}

	if !(byteSize > 0 && byteSize <= 32) {
		return "", fmt.Errorf("Invalid size for type %s", byteType)
	}

	value = strings.TrimPrefix(value, "0x")

	// Check if the string is a valid hexadecimal
	_, err = strconv.ParseUint(value, 16, 64)
	if err != nil {
		return "", err
	}

	if len(value) != byteSize*2 {
		return "", fmt.Errorf("Invalid string length %s", value)
	}

	paddedHex := rightPadWithZeros(value, 64)

	return paddedHex, nil
}

// ConvertAddress converts an address input to the address encoding which is left padded with 0s until 32 bytes (or 64 characters for the string hex).
// Example: 0x00...00<address>
func ConvertAddress(value string) (string, error) {
	address := ethcommon.FromHex(value)

	hexEncoded := hex.EncodeToString(address) // Encode bytes as hex string
	hexString := leftPadWithZeros(hexEncoded, 64)
	return hexString, nil
}

// GetFunctionSignatureObject returns the FunctionSignature representation of a given function signature string
func GetFunctionSignatureObject(functionSig string) (FunctionSignature, error) {
	functionSig, err := ExtractFunctionNameAndFunctionArgs(functionSig)
	if err != nil {
		return FunctionSignature{}, err
	}

	functionSigObject, err := FunctionSignatureParser.ParseString("", functionSig)
	if err != nil {
		return FunctionSignature{}, fmt.Errorf("Failed to parse function sig %s. Error: %v", functionSig, err)
	}

	return *functionSigObject, nil
}

// ExtractFunctionNameAndFunctionArgs takes any form of function input and return the first parenthesis encountered.
// Example:
// Input: "someFuncName(uint256,(string,string)[][],(uint256,(bool,string)))(uint,string)"
// Returns: "someFuncName(uint256,(string,string)[][],(uint256,(bool,string)))"
func ExtractFunctionNameAndFunctionArgs(input string) (string, error) {
	input = strings.Replace(input, " ", "", -1) // remove spaces per solidity doc: https://docs.soliditylang.org/en/latest/abi-spec.html#function-selector

	var depth int
	var endIndex int
	var found bool

loop:
	for i, char := range input {
		switch char {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				endIndex = i
				found = true
				break loop
			}
		}
	}

	if !found || endIndex == 0 {
		return "", fmt.Errorf("ExtractFunctionNameAndFunctionArgs: invalid parenthesis %s", input)
	}

	return input[:endIndex+1], nil
}

// AbiEncode Takes a function signature and function inputs and returns the encoded byte calldata
func AbiEncode(functionSig string, functionInputs []string) (string, error) {
	functionSigObject, err := GetFunctionSignatureObject(functionSig)
	if err != nil {
		return "", err
	}

	encoding, err := functionSigObject.Encode(functionInputs)
	if err != nil {
		return "", err
	}

	hashedFunctionSig, err := HashFunctionSelector(functionSig)
	if err != nil {
		return "", err
	}

	return "0x" + hashedFunctionSig + encoding, nil
}

// HashFunctionSelector returns a function selectors keccak256 hashed encoding
func HashFunctionSelector(functionSig string) (string, error) {
	// TODO: There are some things that can be improved here, such as:
	// - Error handling for types that don't exist
	// - Converting exact "int" and "uint" types to "int256" and "uint256" automatically
	shortenedFunctionSig, err := ExtractFunctionNameAndFunctionArgs(functionSig)
	if err != nil {
		return "", err
	}
	hashedFunctionSig := fmt.Sprintf("%x", crypto.Keccak256([]byte(shortenedFunctionSig))[:4])

	return hashedFunctionSig, nil
}

// ValidateStringIsQuoted takes a string and validates that it's surrounded by quotations.
// Example:
// - `abc` returns `"abc"`
// - `"abc"` returns `"abc"`
func ValidateStringIsQuoted(s string) string {
	if len(s) == 0 {
		return `""`
	}
	if len(s) > 0 && s[0] != '"' {
		s = `"` + s
	}
	if len(s) > 1 && s[len(s)-1] != '"' {
		s = s + `"`
	}
	return s
}
