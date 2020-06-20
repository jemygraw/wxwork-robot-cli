package cmds

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type WxWorkRobotProfile struct {
	Description string `json:"description"`
	HookKey     string `json:"hook_key"`
}

var AddCmd = cobra.Command{
	Use:     "add",
	Short:   "Add a new robot with name, description and hookKey",
	Long:    "添加一个新的机器人或者覆盖已有机器人的信息，需要指定名称，描述和Hook URL地址中的Key",
	Example: "$ wxrobot add <name> <description> <hookKey>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			cmd.Help()
			os.Exit(1)
		}
		robotName := args[0]
		robotDescription := args[1]
		robotHookKey := args[2]
		// try to write the user profile
		if wErr := writeRobotProfile(robotName, robotDescription, robotHookKey); wErr != nil {
			fmt.Println("Err:", wErr.Error())
			os.Exit(1)
		}
	},
}

func writeRobotProfile(robotName, robotDescription, robotHookKey string) (err error) {
	currentUser, getErr := user.Current()
	if getErr != nil {
		err = fmt.Errorf("get current user error, %s", getErr.Error())
		return
	}
	wxRobotProfileDir := filepath.Join(currentUser.HomeDir, ".wxwork")
	if mkdirErr := os.MkdirAll(wxRobotProfileDir, 0755); mkdirErr != nil {
		err = fmt.Errorf("mkdir for robot profile error, %s", mkdirErr.Error())
		return
	}
	wxRobotProfileFilePath := filepath.Join(wxRobotProfileDir, "robots.json")
	// try to read the robot profile file
	var wxworkRobotProfiles map[string]WxWorkRobotProfile
	robotProfileData, readErr := ioutil.ReadFile(wxRobotProfileFilePath)
	if readErr == nil {
		json.Unmarshal(robotProfileData, &wxworkRobotProfiles)
	}
	if wxworkRobotProfiles == nil {
		wxworkRobotProfiles = make(map[string]WxWorkRobotProfile)
	}
	// add the new robot or overwrite the old one
	wxworkRobotProfiles[robotName] = WxWorkRobotProfile{
		Description: robotDescription,
		HookKey:     robotHookKey,
	}
	robotProfileData, _ = json.Marshal(&wxworkRobotProfiles)
	if writeErr := ioutil.WriteFile(wxRobotProfileFilePath, robotProfileData, 0644); writeErr != nil {
		err = fmt.Errorf("write robot profile file error, %s", writeErr.Error())
		return
	}
	return
}
