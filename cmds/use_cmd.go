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

// use the environment variable to temporarily change the shell prompt
// export PS1='%{$fg_bold[green]%}WXRobot:%{$reset_color%}%{$fg_bold[red]%} (<robot_name>) %{$reset_color%}${ret_status}%{$fg[green]%}%p%{$reset_color%}'

const PS1ChangeFormat = `export PS1='%%{$fg_bold[green]%%}WXRobot:%%{$reset_color%%}%%{$fg_bold[red]%%} (%s) %%{$reset_color%%}${ret_status}%%{$fg[green]%%}%%p%%{$reset_color%%}'
export WXWORK_ROBOT_NAME=%s
`

var UseCmd = cobra.Command{
	Use:     "use",
	Short:   "Use an existing robot to send messages",
	Long:    "切换到一个已存在的机器人来发送消息",
	Example: "$ wxrobot use <name>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(1)
		}

		robotName := args[0]
		// create a temp file with a base_profile to change shell prompt
		wxRobot, getErr := readRobotProfile(robotName)
		if getErr != nil {
			fmt.Println("Err:", getErr.Error())
			os.Exit(1)
		}
		bashProfileTempFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s_bash_profile", robotName))
		bashProfileTempFileContent := fmt.Sprintf(PS1ChangeFormat, fmt.Sprintf("%s-%s", robotName, wxRobot.Description), robotName)
		if writeErr := ioutil.WriteFile(bashProfileTempFilePath, []byte(bashProfileTempFileContent), 0644); writeErr != nil {
			fmt.Println("Err: write robot bash profile error,", writeErr.Error())
			os.Exit(1)
		}
		fmt.Println(fmt.Sprintf("Run command `source %s` to make the robot bash profile effective", bashProfileTempFilePath))
	},
}

func readRobotProfile(robotName string) (wxRobot WxWorkRobot, err error) {
	currentUser, getErr := user.Current()
	if getErr != nil {
		err = fmt.Errorf("get current user error, %s", getErr.Error())
		return
	}
	robotProfileDir := filepath.Join(currentUser.HomeDir, ".wxwork")
	robotProfileFilePath := filepath.Join(robotProfileDir, "robots.json")
	// try to read the robot profile file
	var wxworkRobots map[string]WxWorkRobot
	robotProfileData, readErr := ioutil.ReadFile(robotProfileFilePath)
	if readErr != nil {
		err = fmt.Errorf("read robot profile error, %s", readErr.Error())
		return
	}
	if parseErr := json.Unmarshal(robotProfileData, &wxworkRobots); parseErr != nil {
		err = fmt.Errorf("parse robot profile error, %s", parseErr.Error())
		return
	}
	if v, exists := wxworkRobots[robotName]; !exists {
		err = fmt.Errorf("robot not found by the specified name")
		return
	} else {
		wxRobot = v
	}
	return
}
