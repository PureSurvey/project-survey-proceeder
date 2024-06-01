package trackers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"project-survey-proceeder/internal/configuration"
	"project-survey-proceeder/internal/enums"
	"project-survey-proceeder/internal/events/model"
	"strconv"
	"strings"
)

type Decryptor struct {
	encryptionKey []byte
}

func NewDecryptor(appConfig *configuration.AppConfiguration) *Decryptor {
	key, _ := hex.DecodeString(appConfig.EncryptionSecret)

	return &Decryptor{encryptionKey: key}
}

func (e *Decryptor) DecryptEvent(eventString string) (*model.Event, error) {
	decrypted, _ := decrypt(eventString, e.encryptionKey)

	event, err := e.getEventFromString(decrypted)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func decrypt(encryptedString string, key []byte) (string, error) {
	encryptedBytes, err := base64.URLEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext), nil
}

func (e *Decryptor) getEventFromString(eventString string) (*model.Event, error) {
	params := strings.Split(eventString, ",")
	if len(params) < 8 {
		return nil, errors.New("")
	}

	event := &model.Event{
		ValidQuestionsWithAnswers: map[int][]int{},
	}

	et, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}
	event.EventType = enums.EventType(et)

	unitId, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, err
	}
	event.UnitId = unitId

	validTo, err := strconv.ParseInt(params[2], 10, 64)
	if err != nil {
		return nil, err
	}
	event.ValidTo = validTo

	if params[3] != "" {
		validSurveys := strings.Split(params[3], ";")
		ids, err := convert(validSurveys)
		if err != nil {
			return nil, err
		}

		event.ValidSurveys = ids
	}

	if params[4] != "" {
		validQuestions := strings.Split(params[4], ";")
		ids, err := convert(validQuestions)
		if err != nil {
			return nil, err
		}

		event.ValidQuestions = ids
	}

	if params[5] != "" {
		validQuestionsWithAnswers := strings.Split(params[5], "/")
		for _, answer := range validQuestionsWithAnswers {
			idx := strings.Index(answer, ":")
			qi, err := strconv.Atoi(answer[:idx])
			if err != nil {
				return nil, err
			}

			validOptions := strings.Split(answer[idx+1:], ";")
			ids, err := convert(validOptions)
			if err != nil {
				return nil, err
			}

			event.ValidQuestionsWithAnswers[qi] = ids
		}
	}

	return event, nil
}

func convert(arr []string) ([]int, error) {
	var result []int
	for _, el := range arr {
		val, err := strconv.Atoi(el)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}

	return result, nil
}
