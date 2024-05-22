package utils

import (
	"errors"
	"math"
	"math/rand"
)

func computeParityBit(encoded []int, positions []int) int {
	sum := 0
	for _, pos := range positions {
		sum += encoded[pos]
	}
	return sum % 2
}

func EncodeHamming74(data []int) []int {
	encoded := make([]int, 7)
	encoded[0] = data[0]
	encoded[1] = data[1]
	encoded[2] = data[2]
	encoded[4] = data[3]

	encoded[3] = computeParityBit(encoded, []int{0, 1, 2})
	encoded[5] = computeParityBit(encoded, []int{0, 1, 4})
	encoded[6] = computeParityBit(encoded, []int{0, 2, 4})

	return encoded
}

func EncodeData(data []int) []int {
	encodedMessage := make([]int, 0)

	for i := 0; i < len(data); i += 4 {
		block := data[i:min(i+4, len(data))]   // Блок размером 4 байта
		encodedBlock := EncodeHamming74(block) // Кодируем блок
		encodedMessage = append(encodedMessage, encodedBlock...)
	}

	return encodedMessage
}

func SetError(data []int) []int {
	errorPos := rand.Intn(len(data))
	data[errorPos] = (data[errorPos] + 1) % 2

	return data
}

func DecodeHamming74(data []int) []int {
	decodedData := make([]int, 4)
	syndrome := make([]int, 3)
	syndrome[0] = (data[0] + data[1] + data[2] + data[3]) % 2
	syndrome[1] = (data[0] + data[1] + data[4] + data[5]) % 2
	syndrome[2] = (data[0] + data[2] + data[4] + data[6]) % 2

	if (syndrome[0] + syndrome[1] + syndrome[2]) > 0 {
		pos := syndrome[0]*4 + syndrome[1]*2 + syndrome[2]
		data[int(math.Abs(float64(pos-7)))] = (data[int(math.Abs(float64(pos-7)))] + 1) % 2
	}

	decodedData[0] = data[0]
	decodedData[1] = data[1]
	decodedData[2] = data[2]
	decodedData[3] = data[4]

	return decodedData
}

func DecodeData(data []int) []int {
	decodedMessage := make([]int, 0)

	for i := 0; i < len(data); i += 7 {
		block := data[i:min(i+7, len(data))]
		decodedBlock := DecodeHamming74(block)
		decodedMessage = append(decodedMessage, decodedBlock...)
	}

	return decodedMessage
}

// StringToBitArray переводит строку символов в массив битов (0 и 1)
func StringToBitArray(s string) []int {
	var bitArray []int
	for _, c := range s {
		// Получаем 8-битное представление каждого символа
		for i := 7; i >= 0; i-- {
			bit := (c >> i) & 1
			bitArray = append(bitArray, int(bit))
		}
	}

	return bitArray
}

// BitArrayToString переводит массив битов в строку символов
func BitArrayToString(bitArray []int) (string, error) {
	if len(bitArray)%8 != 0 {
		return "", errors.New("the length of bit array must be a multiple of 8")
	}

	bytes := make([]byte, len(bitArray)/8)
	for i := 0; i < len(bytes); i++ {
		for j := 0; j < 8; j++ {
			bytes[i] = bytes[i]<<1 | byte(bitArray[i*8+j])
		}
	}

	return string(bytes), nil
}
