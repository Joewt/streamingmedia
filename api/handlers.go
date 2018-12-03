package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yinrenxin/streamingmedia/api/dbops"
	"github.com/yinrenxin/streamingmedia/api/defs"
	"github.com/yinrenxin/streamingmedia/api/session"
	"github.com/yinrenxin/streamingmedia/api/utils"
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
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%v", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	uname := p.ByName("username")
	if ubody.Username != uname {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	pwd, err := dbops.GetUserCredential(ubody.Username)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	id := session.GenerateNewSessionId(ubody.Username)
	ui := &defs.SignedIn{Success: true, SessionId: id}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Printf("Unathorized user\n")
		return
	}
	uname := p.ByName("username")
	u, err := dbops.GetUser(uname)
	if err != nil {
		log.Printf("error get userinfo ", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	ui := &defs.UserInfo{Id: u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Printf("Unathorized user\n")
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal(res, nvbody); err != nil {
		log.Printf("%v", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := dbops.AddNewVideo(nvbody.AuthorId, nvbody.Name)
	if err != nil {
		log.Printf("Error in AddNewVideo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	if resp, err := json.Marshal(vi); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 201)
	}
}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Printf("Unathorized user\n")
		return
	}

	uname := p.ByName("username")
	vs, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ListAllvideos: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	vsi := &defs.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(vsi); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	vid := p.ByName("vid-id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w, "", 204)
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	cbody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, cbody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
	} else {
		sendNormalResponse(w, "ok", 201)
	}
}

func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	cms := &defs.Comments{Comments: cm}
	if resp, err := json.Marshal(cms); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}
