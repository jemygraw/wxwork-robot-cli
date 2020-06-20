package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const ListOutputFormat = "%-20s%-20s"

var ListCmd = cobra.Command{
	Use:   "list",
	Short: "List all the robot profiles",
	Long:  "列出当前所有的企微机器人信息",
	Run: func(cmd *cobra.Command, args []string) {
		wxworkRobotProfiles, getErr := readRobotProfiles()
		if getErr != nil {
			fmt.Println("Err:", getErr.Error())
			os.Exit(1)
		}
		fmt.Println(fmt.Sprintf(ListOutputFormat, "Name", "Description"))
		fmt.Println(fmt.Sprintf(ListOutputFormat, "----", "-----------"))
		for robotName, wxWorkRobot := range wxworkRobotProfiles {
			fmt.Println(fmt.Sprintf(ListOutputFormat, robotName, wxWorkRobot.Description))
		}
	},
}
