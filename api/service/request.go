package service

import (
	"bytes"
	"encoding/json"

	"github.com/MarkSmersh/nil-chat/db"
)

type RequestData struct {
	Method Method          `json:"method"`
	ReqID  string          `json:"reqId"`
	Params json.RawMessage `json:"params"`
}

type Request struct {
	UserID int
	DB     *db.DB
}

type Method int

const (
	NewMessage Method = iota + 1
	DeleteMessage
)

func NewRequest(db *db.DB) Request {
	r := Request{
		DB: db,
	}

	return r
}

func (r Request) WithAuth(userId int) Request {
	r.DB = r.DB.WithAuth(userId)
	return r
}

func (r Request) Send(data []byte) Response {
	var req RequestData

	if err := json.Unmarshal(data, &req); err != nil {
		return NewError(req.ReqID, "Nieprawidlowe żądanie")
	}

	if req.ReqID == "" {
		return NewError("", "Brak parametru reqId")
	}

	var f func(data []byte, reqId string) Response

	switch req.Method {
	case NewMessage:
		f = methodWrapper(r.DB.Message().Send)
	case DeleteMessage:
		f = methodWrapper(r.DB.Message().Delete)
	default:
		return NewError(req.ReqID, "Niestający typ żądania")
	}

	return f(req.Params, req.ReqID)
}

func methodWrapper[T, K any, V []byte](f func(T) (K, error)) func(data V, reqId string) Response {
	return func(data V, reqId string) Response {

		read := bytes.NewReader(data)

		d := json.NewDecoder(read)
		d.UseNumber()

		var params T

		if err := d.Decode(&params); err != nil {
			return NewError(reqId, "Nieprawidlowe parametry żądania")
		}

		res, err := f(params)

		if err != nil {
			return NewError(reqId, err.Error())
		}

		return NewResponse(reqId, res)
	}
}
