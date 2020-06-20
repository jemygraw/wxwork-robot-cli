package cmds

import (
	"fmt"
	"github.com/duoke/base/chatops/wechat"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var sendCmdRobotName string
var sendCmdRobotTextMessage string
var sendCmdRobotMarkdownFile string
var sendCmdRobotImageFile string
var sendCmdRobotMediaFile string
var sendCmdRobotNewsFile string

func init() {
	SendCmd.PersistentFlags().StringVarP(&sendCmdRobotName, "robot", "r", "", "robot name")
	SendCmd.PersistentFlags().StringVarP(&sendCmdRobotTextMessage, "text", "t", "", "text message")
	SendCmd.PersistentFlags().StringVarP(&sendCmdRobotMarkdownFile, "markdown", "m", "", "markdown file")
	SendCmd.PersistentFlags().StringVarP(&sendCmdRobotImageFile, "image", "i", "", "image file")
	SendCmd.PersistentFlags().StringVarP(&sendCmdRobotMediaFile, "file", "f", "", "media file")
	SendCmd.PersistentFlags().StringVarP(&sendCmdRobotNewsFile, "news", "n", "", "news file")
}

// SendCmd first check whether the robot name specified in the command args, if not, it will
// check the environment variable WXWORK_ROBOT_NAME to find the robot name.
// If none of the robot name found, it emits error.
var SendCmd = cobra.Command{
	Use:   "send",
	Short: "Send various kinds of messages supported by the robot",
	Long:  "发送企业微信机器人支持的各种类型的消息，支持文本，Markdown，图片，资讯，文件",
	Run: func(cmd *cobra.Command, args []string) {
		if sendCmdRobotName == "" {
			// check whether there are robot name in environment variable
			sendCmdRobotName = os.Getenv(WXWorkRobotNameEnv)
			if sendCmdRobotName == "" {
				fmt.Println("Err: robot name not specified")
				os.Exit(1)
			}
		}
		var err error
		// check various kinds of messages
		if sendCmdRobotTextMessage != "" {
			err = sendMessage(sendCmdRobotName, wechat.WxMessageTypeText, sendCmdRobotTextMessage)
		} else if sendCmdRobotMarkdownFile != "" {
			err = sendMessage(sendCmdRobotName, wechat.WxMessageTypeMarkdown, sendCmdRobotMarkdownFile)
		} else if sendCmdRobotImageFile != "" {
			err = sendMessage(sendCmdRobotName, wechat.WxMessageTypeImage, sendCmdRobotImageFile)
		} else if sendCmdRobotMediaFile != "" {
			err = sendMessage(sendCmdRobotName, wechat.WxMessageTypeFile, sendCmdRobotMediaFile)
		} else if sendCmdRobotNewsFile != "" {
			err = sendMessage(sendCmdRobotName, wechat.WxMessageTypeNews, sendCmdRobotNewsFile)
		} else {
			err = fmt.Errorf("must specify a kind of message to send")
		}

		if err != nil {
			fmt.Println("Err:", err.Error())
			os.Exit(1)
		}
	},
}

func sendMessage(robotName, messageType string, messageBody string) (err error) {
	wxworkRobotProfiles, getErr := readRobotProfiles()
	if getErr != nil {
		err = getErr
		return
	}
	var wxworkRobotProfile WxWorkRobotProfile
	if v, exists := wxworkRobotProfiles[robotName]; !exists {
		err = fmt.Errorf("robot profile not found, please add it first")
		return
	} else {
		wxworkRobotProfile = v
	}
	wxworkRobot := wechat.NewWxWorkRobot(wxworkRobotProfile.HookKey)
	switch messageType {
	case wechat.WxMessageTypeText:
		if sendErr := wxworkRobot.SendTextMessage(messageBody); sendErr != nil {
			err = fmt.Errorf("send text message error, %s", sendErr.Error())
			return
		}
	case wechat.WxMessageTypeMarkdown, wechat.WxMessageTypeImage, wechat.WxMessageTypeFile, wechat.WxMessageTypeNews:
		messageContent, readErr := ioutil.ReadFile(messageBody)
		if readErr != nil {
			err = fmt.Errorf("read message content error, %s", readErr.Error())
			return
		}
		switch messageType {
		case wechat.WxMessageTypeMarkdown:
			if sendErr := wxworkRobot.SendMarkdownMessage(string(messageContent)); sendErr != nil {
				err = fmt.Errorf("send markdown message error, %s", sendErr.Error())
				return
			}
		case wechat.WxMessageTypeImage:
			if sendErr := wxworkRobot.SendImageMessage(messageContent); sendErr != nil {
				err = fmt.Errorf("send image message error, %s", sendErr.Error())
				return
			}
		case wechat.WxMessageTypeFile:
			// upload file first
			fileName := filepath.Base(messageBody)
			mediaID, _, uploadErr := wxworkRobot.UploadFile(messageContent, fileName)
			if uploadErr != nil {
				err = fmt.Errorf("upload media file error, %s", uploadErr.Error())
				return
			}
			// send file message
			if sendErr := wxworkRobot.SendFileMessage(mediaID); sendErr != nil {
				err = fmt.Errorf("send file message error, %s", sendErr.Error())
				return
			}
		case wechat.WxMessageTypeNews:
			var newsArticles []wechat.WxNewsMessageArticle
			if decodeErr := yaml.Unmarshal(messageContent, &newsArticles); decodeErr != nil {
				err = fmt.Errorf("parse news content error, %s", decodeErr.Error())
				return
			}
			// send news message
			if sendErr := wxworkRobot.SendNewsMessage(newsArticles); sendErr != nil {
				err = fmt.Errorf("send news message error, %s", sendErr.Error())
				return
			}
		}
	}
	return
}
