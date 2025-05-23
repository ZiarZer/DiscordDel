package delete

import (
	"context"
	"fmt"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

type Deleter struct {
	Sdk          *discord.DiscordSdk
	ActionLogger *actions.ActionLogger
}

type DeleteOptions struct {
	DeletePinned             bool
	DeleteThreadFirstMessage bool
}

func (deleter *Deleter) BulkDeleteCrawledData(ctx context.Context, authorIds []types.Snowflake, guildId *types.Snowflake, options DeleteOptions) {
	var action *actions.Action
	if guildId == nil {
		action = deleter.ActionLogger.StartAction("Delete all crawled data", deleter.Sdk.Log, true, true)
	} else {
		action = deleter.ActionLogger.StartAction(fmt.Sprintf("Delete crawled data of guild %s", *guildId), deleter.Sdk.Log, true, true)
	}
	defer deleter.ActionLogger.EndAction(action)
	channelIds, err := deleter.Sdk.Repo.GetChannelsWithPendingMessages(authorIds, guildId)
	if err != nil {
		if channelIds != nil {
			utils.InternalLog(err.Error(), utils.WARNING)
		} else {
			utils.InternalLog(err.Error(), utils.ERROR)
			return
		}
	}
	for i := range channelIds {
		deleter.DeleteChannelCrawledData(ctx, authorIds, channelIds[i], options)
	}
}

func (deleter *Deleter) DeleteChannelCrawledData(ctx context.Context, authorIds []types.Snowflake, channelId types.Snowflake, options DeleteOptions) {
	defer deleter.ActionLogger.EndAction(
		deleter.ActionLogger.StartAction(fmt.Sprintf("Delete crawled data of channel %s", channelId), deleter.Sdk.Log, true, true),
	)
	deleter.deleteChannelCrawledReactions(ctx, authorIds, channelId)
	deleter.deleteChannelCrawledMessages(ctx, authorIds, channelId, options)
}

func (deleter *Deleter) deleteChannelCrawledMessages(ctx context.Context, authorIds []types.Snowflake, channelId types.Snowflake, options DeleteOptions) {
	messages, err := deleter.Sdk.Repo.GetPendingMessagesByChannelId(channelId, authorIds)
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

	channel := deleter.Sdk.GetChannel(ctx, channelId)
	if channel.ThreadMetadata != nil {
		if channel.ThreadMetadata.Locked {
			messageIds := utils.Map(messagesToDelete, func(message types.Message) types.Snowflake { return message.Id })
			deleter.Sdk.Repo.UpdateMessagesStatus(messageIds, "ERROR")
			deleter.Sdk.Log(fmt.Sprintf("Thread %s is locked, skipping %d messages to delete", channelId, len(messageIds)), utils.ERROR)
			return
		} else if channel.ThreadMetadata.Archived {
			deleter.Sdk.UnarchiveThread(ctx, channelId)
		}
	}

	for i := range messagesToDelete {
		success := deleter.Sdk.DeleteMessage(ctx, channelId, messagesToDelete[i].Id)
		if success {
			deleter.Sdk.Repo.UpdateMessagesStatus([]types.Snowflake{messagesToDelete[i].Id}, "DELETED")
			if len(messagesToDelete[i].Content) > 0 {
				deleter.Sdk.Log(messagesToDelete[i].Content, nil)
			}
		} else {
			deleter.Sdk.Repo.UpdateMessagesStatus([]types.Snowflake{messagesToDelete[i].Id}, "ERROR")
		}
	}
}

func (deleter *Deleter) deleteChannelCrawledReactions(ctx context.Context, authorIds []types.Snowflake, channelId types.Snowflake) {
	reactions, err := deleter.Sdk.Repo.GetPendingReactionsByChannelId(channelId, authorIds)
	if err != nil {
		if reactions != nil {
			utils.InternalLog(err.Error(), utils.WARNING)
		} else {
			utils.InternalLog(err.Error(), utils.ERROR)
			return
		}
	}

	if len(reactions) == 0 {
		deleter.Sdk.Log(fmt.Sprintf("No reaction to delete in channel %s", channelId), utils.INFO)
		return
	}

	channel := deleter.Sdk.GetChannel(ctx, channelId)
	if channel.ThreadMetadata != nil {
		if channel.ThreadMetadata.Locked {
			deleter.Sdk.Repo.UpdateReactionsStatus(reactions, "ERROR")
			deleter.Sdk.Log(fmt.Sprintf("Thread %s is locked, skipping %d reactions to delete", channelId, len(reactions)), utils.ERROR)
			return
		} else if channel.ThreadMetadata.Archived {
			deleter.Sdk.UnarchiveThread(ctx, channelId)
		}
	}

	for i := range reactions {
		success := deleter.Sdk.DeleteOwnReaction(ctx, channelId, reactions[i].MessageId, reactions[i].Emoji)
		if success {
			deleter.Sdk.Repo.UpdateReactionsStatus([]types.Reaction{reactions[i]}, "DELETED")
		} else {
			deleter.Sdk.Repo.UpdateReactionsStatus([]types.Reaction{reactions[i]}, "ERROR")
		}
	}
}
