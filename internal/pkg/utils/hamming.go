package utils

import (
	"errors"
	"math"
	"math/rand"
)

// computeParityBit вычисляет бит четности для заданных позиций в массиве.
// encoded - массив данных, positions - позиции для вычисления бита четности.
func computeParityBit(encoded []int, positions []int) int {
	sum := 0
	for _, pos := range positions {
		sum += encoded[pos] // Суммируем значения в заданных позициях
	}
	return sum % 2 // Возвращаем остаток от деления на 2 (бит четности)
}

// EncodeHamming74 кодирует 4-битные данные в 7-битный код Хэмминга (7,4).
// data - массив данных длиной 4 бита.
func EncodeHamming74(data []int) []int {
	encoded := make([]int, 7)
	encoded[0] = data[0] // Первый информационный бит
	encoded[1] = data[1] // Второй информационный бит
	encoded[2] = data[2] // Третий информационный бит
	encoded[4] = data[3] // Четвертый информационный бит

	// Вычисляем и устанавливаем биты четности
	encoded[3] = computeParityBit(encoded, []int{0, 1, 2}) // Первый бит четности
	encoded[5] = computeParityBit(encoded, []int{0, 1, 4}) // Второй бит четности
	encoded[6] = computeParityBit(encoded, []int{0, 2, 4}) // Третий бит четности

	return encoded // Возвращаем закодированные данные
}

// EncodeData кодирует весь массив данных, разбивая его на блоки по 4 бита
// и кодируя каждый блок с помощью Hamming-кода (7,4).
func EncodeData(data []int) []int {
	encodedMessage := make([]int, 0)

	for i := 0; i < len(data); i += 4 {
		block := data[i:min(i+4, len(data))]                     // Извлекаем блок данных длиной 4 бита
		encodedBlock := EncodeHamming74(block)                   // Кодируем блок
		encodedMessage = append(encodedMessage, encodedBlock...) // Добавляем закодированный блок в сообщение
	}

	return encodedMessage // Возвращаем закодированное сообщение
}

// SetError вносит случайную ошибку в данные, инвертируя один случайный бит.
// data - массив данных.
func SetError(data []int) []int {
	errorPos := rand.Intn(len(data))          // Выбираем случайную позицию для ошибки
	data[errorPos] = (data[errorPos] + 1) % 2 // Инвертируем бит (0 становится 1, 1 становится 0)

	return data // Возвращаем данные с ошибкой
}

// DecodeHamming74 декодирует 7-битный код Хэмминга (7,4), исправляя одну ошибку, если она есть.
// data - массив данных длиной 7 бит.
func DecodeHamming74(data []int) []int {
	decodedData := make([]int, 4)
	syndrome := make([]int, 3)
	// Вычисляем синдромы для проверки ошибок
	syndrome[0] = (data[0] + data[1] + data[2] + data[3]) % 2
	syndrome[1] = (data[0] + data[1] + data[4] + data[5]) % 2
	syndrome[2] = (data[0] + data[2] + data[4] + data[6]) % 2

	if (syndrome[0] + syndrome[1] + syndrome[2]) > 0 { // Если сумма синдромов больше 0, есть ошибка
		pos := syndrome[0]*4 + syndrome[1]*2 + syndrome[2]                                  // Вычисляем позицию ошибки
		data[int(math.Abs(float64(pos-7)))] = (data[int(math.Abs(float64(pos-7)))] + 1) % 2 // Исправляем ошибку
	}

	// Извлекаем информационные биты из закодированных данных
	decodedData[0] = data[0]
	decodedData[1] = data[1]
	decodedData[2] = data[2]
	decodedData[3] = data[4]

	return decodedData // Возвращаем декодированные данные
}

// DecodeData декодирует весь массив данных, разбивая его на блоки по 7 бит
// и декодируя каждый блок с помощью Hamming-кода (7,4).
func DecodeData(data []int) []int {
	decodedMessage := make([]int, 0)

	for i := 0; i < len(data); i += 7 {
		block := data[i:min(i+7, len(data))]                     // Извлекаем блок данных длиной 7 бит
		decodedBlock := DecodeHamming74(block)                   // Декодируем блок
		decodedMessage = append(decodedMessage, decodedBlock...) // Добавляем декодированный блок в сообщение
	}

	return decodedMessage // Возвращаем декодированное сообщение
}

// StringToBitArray переводит строку символов в массив битов (0 и 1).
// s - входная строка.
func StringToBitArray(s string) []int {
	var bitArray []int
	for _, c := range s {
		// Получаем 8-битное представление каждого символа
		for i := 7; i >= 0; i-- {
			bit := (c >> i) & 1                   // Извлекаем i-й бит символа
			bitArray = append(bitArray, int(bit)) // Добавляем бит в массив
		}
	}

	return bitArray // Возвращаем массив битов
}

// BitArrayToString переводит массив битов в строку символов.
// bitArray - входной массив битов.
func BitArrayToString(bitArray []int) (string, error) {
	if len(bitArray)%8 != 0 { // Проверяем, что длина массива битов кратна 8
		return "", errors.New("the length of bit array must be a multiple of 8")
	}

	bytes := make([]byte, len(bitArray)/8)
	for i := 0; i < len(bytes); i++ {
		for j := 0; j < 8; j++ {
			bytes[i] = bytes[i]<<1 | byte(bitArray[i*8+j]) // Собираем байт из 8 битов
		}
	}

	return string(bytes), nil // Преобразуем массив байтов в строку и возвращаем
}
