package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"hammingCode/internal/models"
	"hammingCode/internal/pkg/hamming"
	"log"
	"net/http"
)

type HammingHandler struct {
	uc hamming.Usecase
}

func NewHammingHandler(uc hamming.Usecase) *HammingHandler {
	return &HammingHandler{
		uc: uc,
	}
}

// @Summary код Хэмминга
// @Description	Ручка защищает передаваемую информацию [7,4]-кодом Хэмминга.
// @Tags Channel
// @Accept json
// @Produce json
// @Param data body models.Segment true "Данные для передачи"
// @Success 200 {object} models.Response "Успешный ответ"
// @Success 409 {object} models.Response "Пакет утерян"
// @Failure 400 {object} models.Response "Некорректный запрос"
// @Failure 500 {object} models.Response "Внутренняя ошибка сервера"
// @Router /code [post]
func (h *HammingHandler) DataLink(c *gin.Context) {
	var (
		segment models.Segment
		success bool
		err     error
	)
	responseData := models.Response{
		Segment: segment,
	}

	if err := c.ShouldBindJSON(&segment); err != nil {
		log.Println(errors.Join(err, errors.New("invalid data format")))
		c.JSON(http.StatusBadRequest, responseData)

		return
	}

	responseData.Segment.Payload, success, responseData.Error, err = h.uc.ChannelTransmit(segment.Payload)
	if err != nil {
		log.Println(errors.Join(err, errors.New("invalid len of bit array")))
		c.JSON(http.StatusInternalServerError, responseData)

		return
	}

	if !success {
		log.Println(errors.New("segment is lost"))
		c.JSON(http.StatusConflict, responseData)

		return
	}

	log.Println(responseData)

	responseJson, err := json.Marshal(responseData)
	if err != nil {
		log.Println(errors.Join(err, errors.New("internal error")))
		c.JSON(http.StatusInternalServerError, responseData)

		return
	}

	log.Println(responseJson)

	apiUrl := "127.0.0.1:3000/transfer"
	log.Println("URL API:", apiUrl)
	log.Println(responseJson)
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(responseJson))
	if err != nil {
		log.Println(errors.Join(err, errors.New("error send to transport layer")))
		c.JSON(http.StatusBadRequest, responseData)

		return
	}

	defer resp.Body.Close()

	log.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Println(errors.Join(err, errors.New("bad status code at transport layer")))
		c.JSON(http.StatusBadRequest, responseData)

		return
	}

	c.Status(http.StatusOK)
}
