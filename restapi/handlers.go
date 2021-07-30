package restapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/devskiller-tasks/messaging-app-golang/smsproxy"
)

func sendSmsHandler(smsProxy smsproxy.SmsProxy) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			handleError(writer, http.StatusBadRequest, err)
			return
		}

		smsRequest := &SendSmsRequest{}

		err = json.Unmarshal(body, smsRequest)

		if err != nil {
			handleError(writer, http.StatusBadRequest, err)
			return
		}

		sendingResult, err := smsProxy.Send(smsproxy.SendMessage{Message: smsRequest.Content, PhoneNumber: smsRequest.PhoneNumber})

		if err != nil {
			switch err.(type) {
			case *smsproxy.ValidationError:
				handleError(writer, http.StatusBadRequest, err)
				return
			default:
				handleError(writer, http.StatusInternalServerError, err)
				return
			}

		}

		responseBody, err := json.Marshal(sendingResult)
		writer.WriteHeader(http.StatusAccepted)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte("Error serializing response"))
			return
		}

		_, err = writer.Write(responseBody)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte("Error writing HTTP response"))
			return
		}

		// HINT: you can use `handleError()` function when handling any error
		// 1. read SendSmsRequest from request. If error occurs, return HTTP Status 400
		// 2. try sending an SMS using `smsProxy.Send(...)`
		// if `smsProxy.Send(...)` returns error which is of type *smsproxy.ValidationError -> return HTTP Status 400
		// if it's a different error -> return HTTP Status 500
		// 3. if everything went OK, return HTTP Status 202 and serialize `SendingResult` from `smsproxy/api.go`, sending it as Response Body
	}
}

func getSmsStatusHandler(smsProxy smsproxy.SmsProxy) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		messageID, err := getMessageID(request.URL.RequestURI())
		if err != nil {
			handleError(writer, http.StatusInternalServerError, err)
			return
		}
		result, err := smsProxy.GetStatus(messageID.String())
		if err != nil {
			handleError(writer, http.StatusInternalServerError, err)
			return
		}

		responseBody, err := json.Marshal(SmsStatusResponse{Status: result})
		if err != nil {
			handleError(writer, http.StatusInternalServerError, err)
			return
		}

		if _, err = writer.Write(responseBody); err != nil {
			log.Println(errors.Wrapf(err, "cannot write http response").Error())
		}
	}
}

func getMessageID(uri string) (uuid.UUID, error) {
	uriParts := strings.Split(uri, "/")
	parse, err := uuid.Parse(uriParts[2])
	return parse, err
}

func handleError(writer http.ResponseWriter, status int, err error) {
	response := HttpErrorResponse{Error: err.Error()}
	jsonBody, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Error serializing response"))
		log.Println(errors.Wrapf(err, "error serializing json response").Error())
	}
	writer.WriteHeader(status)
	_, err = writer.Write(jsonBody)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(errors.Wrapf(err, "error writing HTTP response").Error())
	}
}
