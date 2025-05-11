package delete

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

type Deleter struct {
	Sdk *discord.DiscordSdk
}

func (deleter *Deleter) DeleteChannelCrawledData(authorizationToken string, authorIds []types.Snowflake, channelId types.Snowflake) {
	deleter.deleteChannelCrawledMessages(authorizationToken, authorIds, channelId)
}

func (deleter *Deleter) deleteChannelCrawledMessages(authorizationToken string, authorIds []types.Snowflake, channelId types.Snowflake) {
	messages, err := deleter.Sdk.Repo.GetMessagesByChannelId(channelId, authorIds)
	if err != nil {
		if messages != nil {
			utils.InternalLog(err.Error(), utils.WARNING)
		} else {
			utils.InternalLog(err.Error(), utils.ERROR)
			return
		}
	}
	if len(messages) == 0 {
		deleter.Sdk.Log(fmt.Sprintf("No message to delete in channel %s", channelId), utils.INFO)
		return
	}

	channel := deleter.Sdk.GetChannel(channelId, authorizationToken)
	if channel.ThreadMetadata != nil {
		if channel.ThreadMetadata.Locked {
			var messageIds []types.Snowflake
			for i := range messages {
				messageIds = append(messageIds, messages[i].Id)
			}
			deleter.Sdk.Repo.UpdateMultipleMessageStatus(messageIds, "ERROR")
			deleter.Sdk.Log(fmt.Sprintf("Thread %s is locked, skipping %d messages to delete", channelId, len(messageIds)), utils.ERROR)
			return
		} else if channel.ThreadMetadata.Archived {
			deleter.Sdk.UnarchiveThread(authorizationToken, channelId)
		}
	}

	var deletedMessageIds []types.Snowflake
	var failedToDeleteMessageIds []types.Snowflake

	for i := range messages {
		success := deleter.Sdk.DeleteMessage(authorizationToken, channelId, messages[i].Id)
		if success {
			deletedMessageIds = append(deletedMessageIds, messages[i].Id)
			if len(messages[i].Content) > 0 {
				deleter.Sdk.Log(messages[i].Content, nil)
			}
		} else {
			failedToDeleteMessageIds = append(failedToDeleteMessageIds, messages[i].Id)
		}
	}

	deleter.Sdk.Repo.UpdateMultipleMessageStatus(deletedMessageIds, "DELETED")
	deleter.Sdk.Repo.UpdateMultipleMessageStatus(failedToDeleteMessageIds, "ERROR")
}
