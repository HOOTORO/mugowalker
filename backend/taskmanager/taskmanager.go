package taskmanager

import (
	"encoding/json"
	"fmt"
	"mugowalker/backend"
	"mugowalker/backend/afk"

	"mugowalker/backend/bot"
	"mugowalker/backend/image"

	st "mugowalker/backend/settings"
	"strings"
)

type TaskManager struct {
	*backend.Config
	*afk.Daywalker
}

func NewTaskManager(bc *backend.Config, out func(string, string)) *TaskManager {
	f := func(a1, a2 string) { out("message", fmt.Sprintf("%v |> %v", a1, a2)) }
	ocr := image.NewEngine(bc)
	bb := afk.Nightstalker(bot.New(f, ocr), bc.Pilot)

	if bc.DevicePath != "" {
		bc.Online = bb.Connect(bc.DevicePath)
	}

	bc.Log.Debug(fmt.Sprintf("[TM]: %v", bb))
	return &TaskManager{Config: bc, Daywalker: bb}
}

func (tm *TaskManager) InitDevice(str string) bool {
	if !tm.Connect(str) {
		return false
	}
	gamestatus := tm.PS(tm.GameId)
	tm.Log.Debug(gamestatus)
	if gamestatus != "" {
		tm.KillApp(strings.Split(tm.GameId, "/")[0])
	}
	err := tm.StartApp(tm.GameId)
	tm.Log.Debug(fmt.Sprintf("[REMGaMe]: %v", err))

	return true
}

func (tm *TaskManager) UpdateConfig(data ...interface{}) {
	dummy := &st.Settings{}
	var jsonStruct map[string]interface{} = data[0].(map[string]interface{})
	tob, err := json.Marshal(jsonStruct)
	if err != nil {
		tm.Log.Error("err during update settings")
	}
	json.Unmarshal(tob, dummy)
	tm.Settings = dummy
	tm.Store(tob, "conf.json", true)

}
func (tm *TaskManager) UpdateAnyConfig(file string, cfgobj interface{}, data ...interface{}) {
	var jsonStruct map[string]interface{} = data[0].(map[string]interface{})
	tob, err := json.Marshal(jsonStruct)
	if err != nil {
		tm.Log.Error("Err during update settings")
	}
	json.Unmarshal(tob, &cfgobj)
	tm.Store(tob, file, true)
	// tm.Settings = dummy
}
func (tm *TaskManager) UpdatePilot(data ...interface{}) {
	dummy := &st.Pilot{}
	var jsonStruct map[string]interface{} = data[0].(map[string]interface{})
	tob, err := json.Marshal(jsonStruct)
	if err != nil {
		tm.Log.Error("Err during update settings")
	}
	json.Unmarshal(tob, dummy)
	tm.Store(tob, "acc.json", true)
}

func (tm *TaskManager) RunTask(t string) {
	if strings.Contains(t, "Daily") {
		afk.Daily(tm.Daywalker)
	}
}
