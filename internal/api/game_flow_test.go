package api_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/kairodrad/donkey/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestJoinCreatesSinglePlayerPerUser(t *testing.T) {
	ts := httptest.NewServer(server.New())
	defer ts.Close()
	client := ts.Client()
	reg := func(name string) map[string]string {
		resp, _ := client.Post(ts.URL+"/api/register", "application/json", bytes.NewBufferString(`{"name":"`+name+`"}`))
		var u map[string]string
		json.NewDecoder(resp.Body).Decode(&u)
		return u
	}
	u1 := reg("A")
	u2 := reg("B")
	resp, _ := client.Post(ts.URL+"/api/game/start", "application/json", bytes.NewBufferString(`{"requesterId":"`+u1["id"]+`"}`))
	var g map[string]string
	json.NewDecoder(resp.Body).Decode(&g)
	client.Post(ts.URL+"/api/game/join", "application/json", bytes.NewBufferString(`{"gameId":"`+g["gameId"]+`","userId":"`+u2["id"]+`"}`))
	stateResp, _ := client.Get(ts.URL + "/api/game/" + g["gameId"] + "/state/" + u1["id"])
	var state struct{ Players []interface{} }
	json.NewDecoder(stateResp.Body).Decode(&state)
	assert.Equal(t, 2, len(state.Players))
}

func TestFinalizeDealsAndLogs(t *testing.T) {
	ts := httptest.NewServer(server.New())
	defer ts.Close()
	client := ts.Client()
	reg := func(name string) map[string]string {
		resp, _ := client.Post(ts.URL+"/api/register", "application/json", bytes.NewBufferString(`{"name":"`+name+`"}`))
		var u map[string]string
		json.NewDecoder(resp.Body).Decode(&u)
		return u
	}
	u1 := reg("A")
	u2 := reg("B")
	resp, _ := client.Post(ts.URL+"/api/game/start", "application/json", bytes.NewBufferString(`{"requesterId":"`+u1["id"]+`"}`))
	var g map[string]string
	json.NewDecoder(resp.Body).Decode(&g)
	client.Post(ts.URL+"/api/game/join", "application/json", bytes.NewBufferString(`{"gameId":"`+g["gameId"]+`","userId":"`+u2["id"]+`"}`))
	client.Post(ts.URL+"/api/game/finalize", "application/json", bytes.NewBufferString(`{"gameId":"`+g["gameId"]+`","userId":"`+u1["id"]+`"}`))
	stateResp, _ := client.Get(ts.URL + "/api/game/" + g["gameId"] + "/state/" + u1["id"])
	var state struct {
		Players []struct {
			Cards []string `json:"cards"`
		}
		HasStarted bool `json:"hasStarted"`
	}
	json.NewDecoder(stateResp.Body).Decode(&state)
	assert.True(t, state.HasStarted)
	if assert.Equal(t, 2, len(state.Players)) {
		assert.Equal(t, 26, len(state.Players[0].Cards))
		assert.Equal(t, 0, len(state.Players[1].Cards))
	}
	logsResp, _ := client.Get(ts.URL + "/api/game/" + g["gameId"] + "/logs")
	var logs []struct {
		Message string `json:"message"`
	}
	json.NewDecoder(logsResp.Body).Decode(&logs)
	assert.Equal(t, "Cards dealt. Begin game.", logs[0].Message)
}

func TestAbandonGame(t *testing.T) {
	ts := httptest.NewServer(server.New())
	defer ts.Close()
	client := ts.Client()
	reg := func(name string) map[string]string {
		resp, _ := client.Post(ts.URL+"/api/register", "application/json", bytes.NewBufferString(`{"name":"`+name+`"}`))
		var u map[string]string
		json.NewDecoder(resp.Body).Decode(&u)
		return u
	}
	u1 := reg("A")
	resp, _ := client.Post(ts.URL+"/api/game/start", "application/json", bytes.NewBufferString(`{"requesterId":"`+u1["id"]+`"}`))
	var g map[string]string
	json.NewDecoder(resp.Body).Decode(&g)
	client.Post(ts.URL+"/api/game/abandon", "application/json", bytes.NewBufferString(`{"gameId":"`+g["gameId"]+`","userId":"`+u1["id"]+`"}`))
	stateResp, _ := client.Get(ts.URL + "/api/game/" + g["gameId"] + "/state/" + u1["id"])
	var state struct {
		IsAbandoned bool `json:"isAbandoned"`
	}
	json.NewDecoder(stateResp.Body).Decode(&state)
	assert.True(t, state.IsAbandoned)
}

func TestAdminState(t *testing.T) {
	ts := httptest.NewServer(server.New())
	defer ts.Close()
	client := ts.Client()
	reg := func(name string) map[string]string {
		resp, _ := client.Post(ts.URL+"/api/register", "application/json", bytes.NewBufferString(`{"name":"`+name+`"}`))
		var u map[string]string
		json.NewDecoder(resp.Body).Decode(&u)
		return u
	}
	u1 := reg("A")
	u2 := reg("B")
	resp, _ := client.Post(ts.URL+"/api/game/start", "application/json", bytes.NewBufferString(`{"requesterId":"`+u1["id"]+`"}`))
	var g map[string]string
	json.NewDecoder(resp.Body).Decode(&g)
	client.Post(ts.URL+"/api/game/join", "application/json", bytes.NewBufferString(`{"gameId":"`+g["gameId"]+`","userId":"`+u2["id"]+`"}`))
	client.Post(ts.URL+"/api/game/finalize", "application/json", bytes.NewBufferString(`{"gameId":"`+g["gameId"]+`","userId":"`+u1["id"]+`"}`))
	adminResp, _ := client.Get(ts.URL + "/api/admin/game/" + g["gameId"])
	var state struct {
		Players []struct {
			Cards []string `json:"cards"`
		} `json:"players"`
	}
	json.NewDecoder(adminResp.Body).Decode(&state)
	if assert.Equal(t, 2, len(state.Players)) {
		assert.Equal(t, 26, len(state.Players[0].Cards))
		assert.Equal(t, 26, len(state.Players[1].Cards))
	}
}
