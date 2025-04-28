package crawl

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

type Crawler struct {
	Sdk *discord.DiscordSdk
}

func (crawler *Crawler) CrawlChannel(authorizationToken string, channelId string) {
	channel := crawler.Sdk.GetChannel(channelId, authorizationToken)
	if channel == nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get channel %s", channelId), utils.ERROR)
		return
	}
	messages := crawler.fetchChannelMessages(authorizationToken, channelId, &discord.GetChannelMessagesOptions{Around: &channel.LastMessageId})
	for len(messages) > 0 {
		oldestReadMessageId := messages[len(messages)-1].Id
		messages = crawler.fetchChannelMessages(authorizationToken, channelId, &discord.GetChannelMessagesOptions{Before: &oldestReadMessageId})
	}
}

func (crawler *Crawler) fetchChannelMessages(authorizationToken string, channelId string, options *discord.GetChannelMessagesOptions) []types.Message {
	messages := crawler.Sdk.GetChannelMessages(authorizationToken, channelId, options)
	crawler.Sdk.Repo.InsertMultipleMessages(messages)
	return messages
}
