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
		segment models.Segment // Сегмент данных, получаемый из запроса
		success bool           // Успешность передачи данных
		err     error          // Переменная для ошибок
		//responseData models.Response // Ответ на запрос
	)

	// Привязываем JSON из запроса к структуре segment
	if err := c.ShouldBindJSON(&segment); err != nil {
		// Логируем ошибку при неверном формате данных
		log.Println(errors.Join(err, errors.New("invalid data format")))
		// Отправляем ответ с ошибкой 400 (Bad Request)
		c.JSON(http.StatusBadRequest, segment)
		return
	}

	// Инициализируем ответ с полученным сегментом
	//responseData = models.Response{
	//	Segment: segment,
	//}

	// Передаем данные через канал Хэмминга
	segment.Payload, success, segment.Error, err = h.uc.ChannelTransmit(segment.Payload)
	log.Println(segment.Payload)

	log.Println(segment)
	if err != nil {
		// Логируем ошибку при неверной длине массива бит
		log.Println(errors.Join(err, errors.New("invalid len of bit array")))
		// Отправляем ответ с ошибкой 500 (Internal Server Error)
		c.JSON(http.StatusInternalServerError, segment)
		return
	}

	// Проверяем, не потерян ли сегмент
	if !success {
		// Логируем сообщение о потере сегмента
		log.Println(errors.New("segment is lost"))
		// Отправляем ответ с ошибкой 409 (Conflict)
		c.JSON(http.StatusConflict, segment)
		return
	}

	// Логируем ответные данные
	log.Println(segment)

	// Преобразуем ответные данные в JSON
	responseJson, err := json.Marshal(segment)
	if err != nil {
		// Логируем ошибку при преобразовании в JSON
		log.Println(errors.Join(err, errors.New("internal error")))
		// Отправляем ответ с ошибкой 500 (Internal Server Error)
		c.JSON(http.StatusInternalServerError, segment)
		return
	}

	// Логируем JSON-данные и URL API
	log.Println(responseJson)

	apiUrl := "http://192.168.227.33:8000/transfer/"
	log.Println("URL API:", apiUrl)
	log.Println(responseJson)

	// Отправляем POST-запрос на транспортный уровень
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(responseJson))
	if err != nil {
		// Логируем ошибку при отправке данных
		log.Println(errors.Join(err, errors.New("error send to transport layer")))
		// Отправляем ответ с ошибкой 400 (Bad Request)
		c.JSON(http.StatusBadRequest, segment)
		return
	}

	defer resp.Body.Close()

	// Логируем код состояния ответа
	log.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusNoContent {
		// Логируем ошибку при неверном коде состояния
		log.Println(errors.Join(err, errors.New("bad status code at transport layer")))
		// Отправляем ответ с ошибкой 400 (Bad Request)
		c.JSON(http.StatusBadRequest, segment)
		return
	}

	// Устанавливаем код состояния 200 (OK)
	c.Status(http.StatusOK)
}
