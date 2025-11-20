package crawl

import (
	"context"
	"fmt"
	"slices"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (crawler *Crawler) crawlChannelMessages(ctx context.Context, channel *types.Channel, authorIds []types.Snowflake, crawlingInfo *types.CrawlingInfo) error {
	action := actions.NewAction(fmt.Sprintf("Crawl messages in channel %s", channel.Id))
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(action, crawler.Sdk.TempLog, false),
	)
	if channel.Type == types.PublicThread || channel.Type == types.PrivateThread {
		parentChannel, err := crawler.Sdk.GetChannel(ctx, *channel.ParentId)
		if err != nil {
			return err
		}
		if parentChannel.Type == types.GuildForum {
			threadsData, err := crawler.Sdk.GetThreadsData(ctx, *channel.ParentId, []types.Snowflake{channel.Id})
			if err != nil {
				return err
			}
			if threadsData != nil {
				firstMessage := threadsData.Threads[channel.Id].FirstMessage
				status := types.SKIPPED
				if slices.Contains(authorIds, firstMessage.Author.Id) {
					status = types.PENDING
				}
				crawler.Sdk.Repo.InsertMultipleMessages([]types.Message{firstMessage}, status)
			}
		}
	}

	if channel.LastMessageId == nil {
		crawler.Sdk.Log("Channel doesn't contain messages: nothing to do", utils.INFO)
	} else if crawlingInfo != nil {
		messages, err := crawler.fetchChannelMessages(ctx, authorIds, channel.Id, &discord.GetChannelMessagesOptions{After: &crawlingInfo.NewestReadId})
		if err != nil {
			return err
		}
		for len(messages) > 0 {
			newestReadMessageId := messages[0].Id
			messages, err = crawler.fetchChannelMessages(ctx, authorIds, channel.Id, &discord.GetChannelMessagesOptions{After: &newestReadMessageId})
			if err != nil {
				return err
			}
		}
		if !crawlingInfo.ReachedTop {
			messages, err := crawler.fetchChannelMessages(ctx, authorIds, channel.Id, &discord.GetChannelMessagesOptions{Before: &crawlingInfo.OldestReadId})
			if err != nil {
				return err
			}
			for len(messages) > 0 {
				oldestReadMessageId := messages[len(messages)-1].Id
				messages, err = crawler.fetchChannelMessages(ctx, authorIds, channel.Id, &discord.GetChannelMessagesOptions{Before: &oldestReadMessageId})
				if err != nil {
					return err
				}
			}
		}
	} else {
		messages, err := crawler.fetchChannelMessages(ctx, authorIds, channel.Id, nil)
		if err != nil {
			return err
		}
		for len(messages) > 0 {
			oldestReadMessageId := messages[len(messages)-1].Id
			messages, err = crawler.fetchChannelMessages(ctx, authorIds, channel.Id, &discord.GetChannelMessagesOptions{Before: &oldestReadMessageId})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (crawler *Crawler) fetchChannelMessages(ctx context.Context, authorIds []types.Snowflake, channelId types.Snowflake, options *discord.GetChannelMessagesOptions) ([]types.Message, error) {
	messages, err := crawler.Sdk.GetChannelMessages(ctx, channelId, options)
	if err != nil {
		return nil, err
	}

	messagesToStore := utils.Filter(messages, func(message types.Message) bool {
		if message.Type == types.ThreadStarterMessage {
			return false
		}
		return slices.Contains(authorIds, message.Author.Id) ||
			(message.InteractionMetadata != nil && slices.Contains(authorIds, message.InteractionMetadata.Triggerer.Id))
	})
	crawler.Sdk.Repo.InsertMultipleMessages(messagesToStore, types.PENDING)
	for i := range messages {
		if messages[i].Reactions != nil {
			crawler.crawlMessageReactions(ctx, &messages[i], authorIds)
		}
	}
	crawler.storeChannelMessagesCrawlingInfo(channelId, messages, options)
	return messages, nil
}

func (crawler *Crawler) storeChannelMessagesCrawlingInfo(channelId types.Snowflake, fetchedMessages []types.Message, fetchOptions *discord.GetChannelMessagesOptions) {
	if fetchOptions == nil {
		if len(fetchedMessages) > 0 {
			crawler.Sdk.Repo.InsertCrawlingInfo(
				channelId,
				"CHANNEL",
				fetchedMessages[len(fetchedMessages)-1].Id,
				fetchedMessages[0].Id,
				false,
			)
		}
	} else {
		var oldestReadMessageId *types.Snowflake
		if fetchOptions.Before != nil && len(fetchedMessages) > 0 {
			oldestReadMessageId = &fetchedMessages[len(fetchedMessages)-1].Id
		}
		var newestReadMessageId *types.Snowflake
		if fetchOptions.After != nil && len(fetchedMessages) > 0 {
			newestReadMessageId = &fetchedMessages[0].Id
		}
		var reachedTop *bool
		if fetchOptions.Before != nil && len(fetchedMessages) == 0 {
			reachedTop = utils.MakePointer(true)
		}
		crawler.Sdk.Repo.UpdateCrawlingInfo(
			channelId,
			oldestReadMessageId,
			newestReadMessageId,
			reachedTop,
		)
	}
}
