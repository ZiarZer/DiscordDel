package crawl

import (
	"fmt"
	"slices"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

type Crawler struct {
	Sdk *discord.DiscordSdk
}

func (crawler *Crawler) CrawlChannel(authorizationToken string, authorIds []string, channelId string) {
	channel := crawler.Sdk.GetChannel(channelId, authorizationToken)
	if channel == nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get channel %s", channelId), utils.ERROR)
		return
	}

	crawlingInfo, err := crawler.Sdk.Repo.GetChannelCrawlingInfo(channelId)
	if err != nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get crawling info for channel %s", channelId), utils.WARNING)
	}

	if channel.LastMessageId == nil {
		crawler.Sdk.Log("Channel doesn't contain messages: nothing to do", utils.INFO)
	} else if crawlingInfo != nil {
		messages := crawler.fetchChannelMessages(authorizationToken, authorIds, channelId, &discord.GetChannelMessagesOptions{After: &crawlingInfo.NewestReadMessageId})
		for len(messages) > 0 {
			newestReadMessageId := messages[0].Id
			messages = crawler.fetchChannelMessages(authorizationToken, authorIds, channelId, &discord.GetChannelMessagesOptions{After: &newestReadMessageId})
		}
		if !crawlingInfo.ReachedTop {
			for len(messages) > 0 {
				oldestReadMessageId := messages[len(messages)-1].Id
				messages = crawler.fetchChannelMessages(authorizationToken, authorIds, channelId, &discord.GetChannelMessagesOptions{Before: &oldestReadMessageId})
			}
		}
	} else {
		messages := crawler.fetchChannelMessages(authorizationToken, authorIds, channelId, nil)
		for len(messages) > 0 {
			oldestReadMessageId := messages[len(messages)-1].Id
			messages = crawler.fetchChannelMessages(authorizationToken, authorIds, channelId, &discord.GetChannelMessagesOptions{Before: &oldestReadMessageId})
		}
	}
}

func (crawler *Crawler) fetchChannelMessages(authorizationToken string, authorIds []string, channelId string, options *discord.GetChannelMessagesOptions) []types.Message {
	messages := crawler.Sdk.GetChannelMessages(authorizationToken, channelId, options)

	messagesToStore := []types.Message{}
	for i := range messages {
		if slices.Contains(authorIds, messages[i].Author.Id) {
			messagesToStore = append(messagesToStore, messages[i])
		}
	}
	crawler.Sdk.Repo.InsertMultipleMessages(messagesToStore)
	crawler.storeChannelCrawlingInfo(channelId, messages, options)
	return messages
}

func (crawler *Crawler) storeChannelCrawlingInfo(channelId string, fetchedMessages []types.Message, fetchOptions *discord.GetChannelMessagesOptions) {
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
		var oldestReadMessageId *string
		if fetchOptions.Before != nil && len(fetchedMessages) > 0 {
			oldestReadMessageId = &fetchedMessages[len(fetchedMessages)-1].Id
		}
		var newestReadMessageId *string
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
