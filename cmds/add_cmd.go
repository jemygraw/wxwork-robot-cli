package cmds

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type WxWorkRobot struct {
	Description string `json:"description"`
	HookURL     string `json:"hook_url"`
}

var AddCmd = cobra.Command{
	Use:     "add",
	Short:   "Add a new robot with name, description and hookURL",
	Long:    "添加一个新的机器人或者覆盖已有机器人的信息，需要指定名称，描述和Hook URL",
	Example: "$ wxrobot add <name> <description> <hookURL>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			cmd.Help()
			os.Exit(1)
		}
		robotName := args[0]
		robotDescription := args[1]
		robotHookURL := args[2]

		if !strings.HasPrefix(robotHookURL, "https") {
			fmt.Println("Err: invalid robot hook url")
			os.Exit(1)
		}
		// try to write the user profile
		if wErr := writeRobotProfile(robotName, robotDescription, robotHookURL); wErr != nil {
			fmt.Println("Err:", wErr.Error())
			os.Exit(1)
		}
	},
}

func writeRobotProfile(robotName, robotDescription, robotHookURL string) (err error) {
	currentUser, getErr := user.Current()
	if getErr != nil {
		err = fmt.Errorf("get current user error, %s", getErr.Error())
		return
	}
	robotProfileDir := filepath.Join(currentUser.HomeDir, ".wxwork")
	if mkdirErr := os.MkdirAll(robotProfileDir, 0755); mkdirErr != nil {
		err = fmt.Errorf("mkdir for robot profile error, %s", mkdirErr.Error())
		return
	}
	robotProfileFilePath := filepath.Join(robotProfileDir, "robots.json")
	// try to read the robot profile file
	var wxworkRobots map[string]WxWorkRobot
	robotProfileData, readErr := ioutil.ReadFile(robotProfileFilePath)
	if readErr == nil {
		json.Unmarshal(robotProfileData, &wxworkRobots)
	}
	if wxworkRobots == nil {
		wxworkRobots = make(map[string]WxWorkRobot)
	}
	// add the new robot or overwrite the old one
	wxworkRobots[robotName] = WxWorkRobot{
		Description: robotDescription,
		HookURL:     base64.StdEncoding.EncodeToString([]byte(robotHookURL)),
	}
	robotProfileData, _ = json.Marshal(&wxworkRobots)
	if writeErr := ioutil.WriteFile(robotProfileFilePath, robotProfileData, 0644); writeErr != nil {
		err = fmt.Errorf("write robot profile file error, %s", writeErr.Error())
		return
	}
	return
}
