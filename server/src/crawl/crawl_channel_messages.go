package crawl

import (
	"slices"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (crawler *Crawler) crawlChannelMessages(authorizationToken string, channel *types.Channel, authorIds []types.Snowflake, crawlingInfo *types.CrawlingInfo) {
	if channel.Type == types.PublicThread || channel.Type == types.PrivateThread {
		parentChannel := crawler.Sdk.GetChannel(*channel.ParentId, authorizationToken)
		if parentChannel.Type == types.GuildForum {
			threadsData := crawler.Sdk.GetThreadsData(authorizationToken, *channel.ParentId, []types.Snowflake{channel.Id})
			if threadsData != nil {
				firstMessage := threadsData.Threads[channel.Id].FirstMessage
				crawler.Sdk.Repo.InsertMultipleMessages([]types.Message{firstMessage}, "THREAD_FIRST_MESSAGE")
			}
		}
	}

	if channel.LastMessageId == nil {
		crawler.Sdk.Log("Channel doesn't contain messages: nothing to do", utils.INFO)
	} else if crawlingInfo != nil {
		messages := crawler.fetchChannelMessages(authorizationToken, authorIds, channel.Id, &discord.GetChannelMessagesOptions{After: &crawlingInfo.NewestReadId})
		for len(messages) > 0 {
			newestReadMessageId := messages[0].Id
			messages = crawler.fetchChannelMessages(authorizationToken, authorIds, channel.Id, &discord.GetChannelMessagesOptions{After: &newestReadMessageId})
		}
		if !crawlingInfo.ReachedTop {
			messages := crawler.fetchChannelMessages(authorizationToken, authorIds, channel.Id, &discord.GetChannelMessagesOptions{Before: &crawlingInfo.OldestReadId})
			for len(messages) > 0 {
				oldestReadMessageId := messages[len(messages)-1].Id
				messages = crawler.fetchChannelMessages(authorizationToken, authorIds, channel.Id, &discord.GetChannelMessagesOptions{Before: &oldestReadMessageId})
			}
		}
	} else {
		messages := crawler.fetchChannelMessages(authorizationToken, authorIds, channel.Id, nil)
		for len(messages) > 0 {
			oldestReadMessageId := messages[len(messages)-1].Id
			messages = crawler.fetchChannelMessages(authorizationToken, authorIds, channel.Id, &discord.GetChannelMessagesOptions{Before: &oldestReadMessageId})
		}
	}
}

func (crawler *Crawler) fetchChannelMessages(authorizationToken string, authorIds []types.Snowflake, channelId types.Snowflake, options *discord.GetChannelMessagesOptions) []types.Message {
	messages := crawler.Sdk.GetChannelMessages(authorizationToken, channelId, options)

	messagesToStore := utils.Filter(messages, func(message types.Message) bool {
		if message.Type == types.ThreadStarterMessage {
			return false
		}
		return slices.Contains(authorIds, message.Author.Id) ||
			(message.InteractionMetadata != nil && slices.Contains(authorIds, message.InteractionMetadata.Triggerer.Id))
	})
	crawler.Sdk.Repo.InsertMultipleMessages(messagesToStore, "PENDING")
	for i := range messages {
		if messages[i].Reactions != nil {
			crawler.crawlMessageReactions(authorizationToken, &messages[i], authorIds)
		}
	}
	crawler.storeChannelMessagesCrawlingInfo(channelId, messages, options)
	return messages
}

func (crawler *Crawler) storeChannelMessagesCrawlingInfo(channelId types.Snowflake, fetchedMessages []types.Message, fetchOptions *discord.GetChannelMessagesOptions) {
	if fetchOptions == nil {
		if len(fetchedMessages) > 0 {
			crawler.Sdk.Repo.InsertChannelCrawlingInfo(
				channelId,
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
		crawler.Sdk.Repo.UpdateChannelCrawlingInfo(
			channelId,
			oldestReadMessageId,
			newestReadMessageId,
			reachedTop,
		)
	}
}
