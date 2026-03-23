package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

// LoggerMemberJoinLog builds the member join log entry.
func LoggerMemberJoinLog(member discord.Member) discord.MessageCreate {
	createdAt := member.User.CreatedAt().Format("2006-01-02")
	return loggerMessage(
		"### 📥 メンバー参加",
		fmt.Sprintf("**ユーザー:** <@%d> (%s)\n**アカウント作成:** %s",
			member.User.ID, member.User.Username, createdAt),
	)
}

// LoggerMemberLeaveLog builds the member leave log entry.
func LoggerMemberLeaveLog(user discord.User) discord.MessageCreate {
	return loggerMessage(
		"### 📤 メンバー退出",
		fmt.Sprintf("**ユーザー:** <@%d> (%s)",
			user.ID, user.Username),
	)
}

// LoggerBanLog builds the ban log entry.
func LoggerBanLog(user discord.User) discord.MessageCreate {
	return loggerMessage(
		"### 🔨 BAN",
		fmt.Sprintf("**ユーザー:** <@%d> (ID: %d)",
			user.ID, user.ID),
	)
}

// LoggerUnbanLog builds the unban log entry.
func LoggerUnbanLog(user discord.User) discord.MessageCreate {
	return loggerMessage(
		"### 🔓 BAN解除",
		fmt.Sprintf("**ユーザー:** <@%d> (ID: %d)",
			user.ID, user.ID),
	)
}

func loggerMessage(title, body string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(title),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
