package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yinrenxin/streamingmedia/api/dbops"
	"github.com/yinrenxin/streamingmedia/api/defs"
	"github.com/yinrenxin/streamingmedia/api/session"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resp, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(resp, ubody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success: true, SessionId: id}
	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 201)
	}

}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	io.WriteString(w, uname)
}
