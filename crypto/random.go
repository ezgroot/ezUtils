package crypto

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var (
	VisualCharAll       = 0
	VisualCharNumber    = 1
	VisualCharCapital   = 2
	VisualCharLowercase = 3
	VisualCharSpecial   = 4
)

var NumberCharMap = make(map[byte]uint64)
var NumberCharList = make([]uint64, 10)

var CapitalCharMap = make(map[byte]uint64)
var CapitalCharList = make([]uint64, 26)

var LowercaseCharMap = make(map[byte]uint64)
var LowercaseCharList = make([]uint64, 26)

var SpecialCharMap = make(map[byte]uint64)
var SpecialCharList = make([]uint64, 128)

var VisualCharList = make([]uint64, 256)

func init() {
	// ASCII number : 48<= char <= 57
	NumberCharMap[48] = 48 // 0
	NumberCharMap[49] = 49 // 1
	NumberCharMap[50] = 50 // 2
	NumberCharMap[51] = 51 // 3
	NumberCharMap[52] = 52 // 4
	NumberCharMap[53] = 53 // 5
	NumberCharMap[54] = 54 // 6
	NumberCharMap[55] = 55 // 7
	NumberCharMap[56] = 56 // 8
	NumberCharMap[57] = 57 // 9

	// ASCII uppercase letter : 65<= char <= 90
	CapitalCharMap[65] = 65 // A
	CapitalCharMap[66] = 66 // B
	CapitalCharMap[67] = 67 // C
	CapitalCharMap[68] = 68 // D
	CapitalCharMap[69] = 69 // E
	CapitalCharMap[70] = 70 // F
	CapitalCharMap[71] = 71 // G
	CapitalCharMap[72] = 72 // H
	CapitalCharMap[73] = 73 // I
	CapitalCharMap[74] = 74 // J
	CapitalCharMap[75] = 75 // K
	CapitalCharMap[76] = 76 // L
	CapitalCharMap[77] = 77 // M
	CapitalCharMap[78] = 78 // N
	CapitalCharMap[79] = 79 // O
	CapitalCharMap[80] = 80 // P
	CapitalCharMap[81] = 81 // Q
	CapitalCharMap[82] = 82 // R
	CapitalCharMap[83] = 83 // S
	CapitalCharMap[84] = 84 // T
	CapitalCharMap[85] = 85 // U
	CapitalCharMap[86] = 86 // V
	CapitalCharMap[87] = 87 // W
	CapitalCharMap[88] = 88 // X
	CapitalCharMap[89] = 89 // Y
	CapitalCharMap[90] = 90 // Z

	// ASCII Lower case letters : 97<= char <= 122
	LowercaseCharMap[97] = 97   // a
	LowercaseCharMap[98] = 98   // b
	LowercaseCharMap[99] = 99   // c
	LowercaseCharMap[100] = 100 // d
	LowercaseCharMap[101] = 101 // e
	LowercaseCharMap[102] = 102 // f
	LowercaseCharMap[103] = 103 // g
	LowercaseCharMap[104] = 104 // h
	LowercaseCharMap[105] = 105 // i
	LowercaseCharMap[106] = 106 // j
	LowercaseCharMap[107] = 107 // k
	LowercaseCharMap[108] = 108 // l
	LowercaseCharMap[109] = 109 // m
	LowercaseCharMap[110] = 110 // n
	LowercaseCharMap[111] = 111 // o
	LowercaseCharMap[112] = 112 // p
	LowercaseCharMap[113] = 113 // q
	LowercaseCharMap[114] = 114 // r
	LowercaseCharMap[115] = 115 // s
	LowercaseCharMap[116] = 116 // t
	LowercaseCharMap[117] = 117 // u
	LowercaseCharMap[118] = 118 // v
	LowercaseCharMap[119] = 119 // w
	LowercaseCharMap[120] = 120 // x
	LowercaseCharMap[121] = 121 // y
	LowercaseCharMap[122] = 122 // z

	// ASCII Special characters :
	SpecialCharMap[35] = 35 // #
	SpecialCharMap[36] = 36 // $
	SpecialCharMap[37] = 37 // %
	SpecialCharMap[38] = 38 // &

	for _, v := range SpecialCharMap {
		VisualCharList = append(VisualCharList, v)
		SpecialCharList = append(SpecialCharList, v)
	}

	for _, v := range NumberCharMap {
		VisualCharList = append(VisualCharList, v)
		NumberCharList = append(NumberCharList, v)
	}

	for _, v := range CapitalCharMap {
		VisualCharList = append(VisualCharList, v)
		CapitalCharList = append(CapitalCharList, v)
	}

	for _, v := range LowercaseCharMap {
		VisualCharList = append(VisualCharList, v)
		LowercaseCharList = append(LowercaseCharList, v)
	}
}

func getVisualChar(charType int) (byte, error) {
	var b byte

	switch charType {
	case VisualCharNumber:
		{
			len := len(NumberCharList)
			index, err := Number(uint64(len))
			if err != nil {
				return b, err
			}

			b = byte(NumberCharList[index])
		}
	case VisualCharCapital:
		{
			len := len(CapitalCharList)
			index, err := Number(uint64(len))
			if err != nil {
				return b, err
			}

			b = byte(CapitalCharList[index])
		}
	case VisualCharLowercase:
		{
			len := len(LowercaseCharList)
			index, err := Number(uint64(len))
			if err != nil {
				return b, err
			}

			b = byte(LowercaseCharList[index])
		}
	case VisualCharSpecial:
		{
			len := len(SpecialCharList)
			index, err := Number(uint64(len))
			if err != nil {
				return b, err
			}

			b = byte(SpecialCharList[index])
		}
	default:
		{
			len := len(VisualCharList)
			index, err := Number(uint64(len))
			if err != nil {
				return b, err
			}

			b = byte(VisualCharList[index])
		}
	}

	return b, nil
}

// Number get random number，0<= value < maxNum.
func Number(maxNum uint64) (uint64, error) {
	max := new(big.Int).SetUint64(maxNum)
	number, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return number.Uint64(), nil
}

// String get random string that length is defined
func String(length int) (string, error) {
	byteList := make([]byte, length)

	_, err := rand.Read(byteList)
	if err != nil {
		return "", err
	}

	for i, v := range byteList {
		_, ok := CapitalCharMap[v]
		if !ok {
			for {
				byteList[i], err = getVisualChar(0)
				if err != nil {
					fmt.Printf("getVisualChar error = %s", err)
				} else {
					break
				}
			}
		}
	}

	return string(byteList), nil
}

// Password generate random password of specified length
func Password(length uint64) (string, error) {
	if length < 6 {
		return "", errors.New("password is too short, min length at 6")
	}

	numberCharPosition, err := Number(length)
	if err != nil {
		return "", err
	}

	var capitalCharPosition uint64
	for {
		capitalCharPosition, err = Number(length)
		if err != nil {
			return "", err
		}

		if capitalCharPosition != numberCharPosition {
			break
		}
	}

	var lowercaseCharPosition uint64
	for {
		lowercaseCharPosition, err = Number(length)
		if err != nil {
			fmt.Printf("%s", err)
			return "", err
		}

		if (lowercaseCharPosition != numberCharPosition) &&
			(lowercaseCharPosition != capitalCharPosition) {
			break
		}
	}

	var specialCharPosition uint64
	for {
		specialCharPosition, err = Number(length)
		if err != nil {
			fmt.Printf("%s", err)
			return "", err
		}

		if (specialCharPosition != numberCharPosition) &&
			(specialCharPosition != capitalCharPosition) &&
			(specialCharPosition != lowercaseCharPosition) {
			break
		}
	}

	pwdList := make([]byte, length)

	_, err = rand.Read(pwdList)
	if err != nil {
		return "", err
	}

	pwdList[numberCharPosition], err = getVisualChar(1)
	if err != nil {
		return "", err
	}

	pwdList[capitalCharPosition], err = getVisualChar(2)
	if err != nil {
		return "", err
	}

	pwdList[lowercaseCharPosition], err = getVisualChar(3)
	if err != nil {
		return "", err
	}

	pwdList[specialCharPosition], err = getVisualChar(4)
	if err != nil {
		return "", err
	}

	for i, v := range pwdList {
		_, ok := CapitalCharMap[v]
		if !ok {
			for {
				pwdList[i], err = getVisualChar(0)
				if err != nil {
					fmt.Printf("getVisualChar error = %s", err)
				} else {
					break
				}
			}
		}
	}

	return string(pwdList), nil
}
