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

func (uc *HammingUsecase) ChannelTransmit(data string) (string, bool, bool, error) {
	var hasError bool
	channel := models.ChannelLevel{
		ProbabilityError: 8,
		ProbabilityLoss:  1,
	}

	bitArray := utils.StringToBitArray(data)

	log.Println(bitArray)
	encodedData := utils.EncodeData(bitArray)

	if rand.Intn(100) < channel.ProbabilityError {
		log.Println("errorFlag")
		encodedData = utils.SetError(encodedData)
		hasError = true
	}

	log.Println("decodedData")
	decodedData := utils.DecodeData(encodedData)
	log.Println("finish")

	result, err := utils.BitArrayToString(decodedData)
	if err != nil {
		return "", false, false, err
	}

	if rand.Intn(100) < channel.ProbabilityLoss {
		return result, false, hasError, nil
	}

	return result, true, hasError, nil
}
