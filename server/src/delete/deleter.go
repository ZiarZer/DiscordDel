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
	messagesToDelete := []types.Message{}
	for i := range messages {
		if !messages[i].Pinned && (messages[i].Status == nil || *messages[i].Status != "THREAD_FIRST_MESSAGE") {
			messagesToDelete = append(messagesToDelete, messages[i])
		}
	}
	if len(messagesToDelete) == 0 {
		deleter.Sdk.Log(fmt.Sprintf("No message to delete in channel %s", channelId), utils.INFO)
		return
	}

	channel := deleter.Sdk.GetChannel(channelId, authorizationToken)
	if channel.ThreadMetadata != nil {
		if channel.ThreadMetadata.Locked {
			var messageIds []types.Snowflake
			for i := range messagesToDelete {
				messageIds = append(messageIds, messagesToDelete[i].Id)
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

	for i := range messagesToDelete {
		success := deleter.Sdk.DeleteMessage(authorizationToken, channelId, messagesToDelete[i].Id)
		if success {
			deletedMessageIds = append(deletedMessageIds, messagesToDelete[i].Id)
			if len(messagesToDelete[i].Content) > 0 {
				deleter.Sdk.Log(messagesToDelete[i].Content, nil)
			}
		} else {
			failedToDeleteMessageIds = append(failedToDeleteMessageIds, messagesToDelete[i].Id)
		}
	}

	deleter.Sdk.Repo.UpdateMultipleMessageStatus(deletedMessageIds, "DELETED")
	deleter.Sdk.Repo.UpdateMultipleMessageStatus(failedToDeleteMessageIds, "ERROR")
}
