package delete

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

type Deleter struct {
	Sdk *discord.DiscordSdk
}

type DeleteOptions struct {
	DeletePinned             bool
	DeleteThreadFirstMessage bool
}

func (deleter *Deleter) DeleteChannelCrawledData(authorizationToken string, authorIds []types.Snowflake, channelId types.Snowflake, options DeleteOptions) {
	defer actions.StartAction(fmt.Sprintf("Delete crawled data of channel %s", channelId), deleter.Sdk.Log, true).EndAction()
	deleter.deleteChannelCrawledMessages(authorizationToken, authorIds, channelId, options)
}

func (deleter *Deleter) deleteChannelCrawledMessages(authorizationToken string, authorIds []types.Snowflake, channelId types.Snowflake, options DeleteOptions) {
	messages, err := deleter.Sdk.Repo.GetMessagesByChannelId(channelId, authorIds)
	if err != nil {
		if messages != nil {
			utils.InternalLog(err.Error(), utils.WARNING)
		} else {
			utils.InternalLog(err.Error(), utils.ERROR)
			return
		}
	}

	messagesToDelete := utils.Filter(messages, func(message types.Message) bool { return message.Type != types.ThreadStarterMessage })
	if !options.DeletePinned {
		messagesToDelete = utils.Filter(messagesToDelete, func(message types.Message) bool { return !message.Pinned })
	}
	if !options.DeleteThreadFirstMessage {
		messagesToDelete = utils.Filter(messagesToDelete, func(message types.Message) bool {
			return message.Status == nil || *message.Status != "THREAD_FIRST_MESSAGE"
		})
	}
	if len(messagesToDelete) == 0 {
		deleter.Sdk.Log(fmt.Sprintf("No message to delete in channel %s", channelId), utils.INFO)
		return
	}

	channel := deleter.Sdk.GetChannel(channelId, authorizationToken)
	if channel.ThreadMetadata != nil {
		if channel.ThreadMetadata.Locked {
			messageIds := utils.Map(messagesToDelete, func(message types.Message) types.Snowflake { return message.Id })
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
