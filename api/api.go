package api

import (
	"adil/vulnerableDB"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"adil/helpers"
)

type Login struct {
	Username string
	Password string
}

type Response struct {
	Data []vulnerableDB.User
}

type ErrResponse struct {
	Message string
}

func login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	var formattedBody Login
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	login := vulnerableDB.VulnerableLogin(formattedBody.Username, formattedBody.Password)

	if len(login) > 0 {
		resp := Response{Data: login}
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(resp)
	}
}
