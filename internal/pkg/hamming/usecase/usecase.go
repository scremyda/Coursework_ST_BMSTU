package usecase

import (
	"hammingCode/internal/models"
	"hammingCode/internal/pkg/utils"
	"log"
	"math/rand"
)

type HammingUsecase struct {
}

func NewHammingUsecase() *HammingUsecase {
	return &HammingUsecase{}
}

// ChannelTransmit выполняет передачу данных через канал, который может вносить ошибки и потери.
func (uc *HammingUsecase) ChannelTransmit(data string) (string, bool, bool, error) {
	var hasError bool

	// Устанавливаем параметры канала: вероятность ошибки и вероятность потери данных
	channel := models.ChannelLevel{
		ProbabilityError: 8,
		ProbabilityLoss:  1,
	}

	// Преобразуем строку в массив битов
	bitArray := utils.StringToBitArray(data)

	// Логируем массив битов для отладки
	log.Println(bitArray)

	// Кодируем данные с использованием Hamming-кода
	encodedData := utils.EncodeData(bitArray)

	// Вносим ошибку в данные с заданной вероятностью
	if rand.Intn(100) < channel.ProbabilityError {
		log.Println("error is set")
		encodedData = utils.SetError(encodedData)
		hasError = true
	}

	// Декодируем данные
	decodedData := utils.DecodeData(encodedData)

	// Преобразуем массив битов обратно в строку
	result, err := utils.BitArrayToString(decodedData)
	if err != nil {
		return "", false, false, err
	}

	// С заданной вероятностью данные теряются (симулируем потерю данных)
	if rand.Intn(100) < channel.ProbabilityLoss {
		return result, false, hasError, nil
	}

	// Возвращаем результат, флаг успешности передачи и флаг наличия ошибки
	return result, true, hasError, nil
}
