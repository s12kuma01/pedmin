package logger

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildMemberJoinLog(member discord.Member) discord.MessageCreate {
	createdAt := member.User.CreatedAt().Format("2006-01-02")
	return logMessage(
		"### 📥 メンバー参加",
		fmt.Sprintf("**ユーザー:** <@%d> (%s)\n**アカウント作成:** %s",
			member.User.ID, member.User.Username, createdAt),
	)
}

func BuildMemberLeaveLog(user discord.User) discord.MessageCreate {
	return logMessage(
		"### 📤 メンバー退出",
		fmt.Sprintf("**ユーザー:** <@%d> (%s)",
			user.ID, user.Username),
	)
}

func BuildBanLog(user discord.User) discord.MessageCreate {
	return logMessage(
		"### 🔨 BAN",
		fmt.Sprintf("**ユーザー:** <@%d> (ID: %d)",
			user.ID, user.ID),
	)
}

func BuildUnbanLog(user discord.User) discord.MessageCreate {
	return logMessage(
		"### 🔓 BAN解除",
		fmt.Sprintf("**ユーザー:** <@%d> (ID: %d)",
			user.ID, user.ID),
	)
}

func logMessage(title, body string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(title),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
