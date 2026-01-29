package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/MarkSmersh/nil-chat/db/notifier"
	"github.com/MarkSmersh/nil-chat/models"
)

type ResType string

var (
	Res ResType = "res"
	Upd ResType = "upd"
	Err ResType = "err"
)

type Response struct {
	Type     ResType        `json:"type"`
	ReqID    string         `json:"reqId,omitempty"`
	Response any            `json:"response,omitempty"`
	Update   *models.Update `json:"update,omitempty"`
	Error    string         `json:"error,omitempty"`
}

// creates a response object with a type 'res'
func NewResponse(reqId string, res any) Response {
	return Response{
		Type:     Res,
		ReqID:    reqId,
		Response: res,
	}
}

// creates a new update response and returns the Update object for filter purposes
// and the Response object
func NewUpdate(channel notifier.Channel, payload []byte) (models.Update, Response) {
	var u models.Update

	slog.Info(
		fmt.Sprintf("%s", string(payload)),
	)

	bindTo := bindPayload(payload)

	channelToModel := map[notifier.Channel]any{
		notifier.NewMessage:    &u.Message,
		notifier.DeleteMessage: &u.DeleteMessage,
	}

	bindTo(channelToModel[channel])

	return u, Response{
		Type:   Upd,
		Update: &u,
	}
}

func NewError(reqId string, err string) Response {
	return Response{
		Type:  Err,
		ReqID: reqId,
		Error: err,
	}
}

func bindPayload(payload []byte) func(out any) error {
	return func(out any) error {
		r := bytes.NewReader(payload)

		d := json.NewDecoder(r)

		d.UseNumber()

		err := d.Decode(&out)

		return err
	}
}
