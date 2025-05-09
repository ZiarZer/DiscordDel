package crawl

import (
	"fmt"
	"slices"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
)

func (crawler *Crawler) crawlMessageReactions(authorizationToken string, message *types.Message, authorIds []types.Snowflake) {
	if message.Reactions == nil {
		return
	}

	for i := range message.Reactions {
		reaction := message.Reactions[i]
		emoji := *reaction.Emoji.Name
		if reaction.Emoji.Id != nil {
			emoji = emoji + fmt.Sprintf(":%s", *reaction.Emoji.Id)
		}
		if reaction.CountDetails.Normal > 0 {
			crawler.crawlReactionsOnEmoji(authorizationToken, message.ChannelId, message.Id, emoji, false, authorIds)
		}
		if reaction.CountDetails.Burst > 0 {
			crawler.crawlReactionsOnEmoji(authorizationToken, message.ChannelId, message.Id, emoji, true, authorIds)
		}
	}
}

func (crawler *Crawler) crawlReactionsOnEmoji(authorizationToken string, channelId types.Snowflake, messageId types.Snowflake, emoji string, isBurst bool, authorIds []types.Snowflake) {
	usersReacted := crawler.fetchReactionsOnEmoji(authorizationToken, channelId, messageId, emoji, isBurst, nil, authorIds)
	pageSize := 100
	options := &discord.GetMessageReactionsOptions{Limit: &pageSize}

	for len(usersReacted) == pageSize {
		options.After = &usersReacted[len(usersReacted)-1].Id
		usersReacted = crawler.fetchReactionsOnEmoji(authorizationToken, channelId, messageId, emoji, isBurst, options, authorIds)
	}
}

func (crawler *Crawler) fetchReactionsOnEmoji(authorizationToken string, channelId types.Snowflake, messageId types.Snowflake, emoji string, isBurst bool, options *discord.GetMessageReactionsOptions, authorIds []types.Snowflake) []types.User {
	usersReacted := crawler.Sdk.GetMessageReactions(authorizationToken, channelId, messageId, emoji, isBurst, options)
	var userIds []types.Snowflake
	for i := range usersReacted {
		if slices.Contains(authorIds, usersReacted[i].Id) && !slices.Contains(userIds, usersReacted[i].Id) {
			userIds = append(userIds, usersReacted[i].Id)
		}
	}
	crawler.Sdk.Repo.InsertMultipleReactions(messageId, userIds, emoji, isBurst)
	return usersReacted
}
