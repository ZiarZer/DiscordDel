package crawl

import (
	"context"
	"fmt"
	"slices"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (crawler *Crawler) crawlMessageReactions(ctx context.Context, message *types.Message, authorIds []types.Snowflake) {
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(fmt.Sprintf("Crawl reactions on message %s", message.Id), crawler.Sdk.TempLog, false, false),
	)
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
			crawler.crawlReactionsOnEmoji(ctx, message.ChannelId, message.Id, emoji, false, authorIds)
		}
		if reaction.CountDetails.Burst > 0 {
			crawler.crawlReactionsOnEmoji(ctx, message.ChannelId, message.Id, emoji, true, authorIds)
		}
	}
}

func (crawler *Crawler) crawlReactionsOnEmoji(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake, emoji string, isBurst bool, authorIds []types.Snowflake) {
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(fmt.Sprintf("Crawl reactions with %s on message %s", emoji, messageId), crawler.Sdk.TempLog, false, false),
	)
	usersReacted := crawler.fetchReactionsOnEmoji(ctx, channelId, messageId, emoji, isBurst, nil, authorIds)
	pageSize := 100
	options := &discord.GetMessageReactionsOptions{Limit: &pageSize}

	for len(usersReacted) == pageSize {
		options.After = &usersReacted[len(usersReacted)-1].Id
		usersReacted = crawler.fetchReactionsOnEmoji(ctx, channelId, messageId, emoji, isBurst, options, authorIds)
	}
}

func (crawler *Crawler) fetchReactionsOnEmoji(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake, emoji string, isBurst bool, options *discord.GetMessageReactionsOptions, authorIds []types.Snowflake) []types.User {
	usersReacted := crawler.Sdk.GetMessageReactions(ctx, channelId, messageId, emoji, isBurst, options)

	userIds := utils.MapWithoutDuplicate(
		utils.Filter(
			usersReacted,
			func(user types.User) bool { return slices.Contains(authorIds, user.Id) },
		),
		func(user types.User) types.Snowflake { return user.Id },
	)
	crawler.Sdk.Repo.InsertMultipleReactions(channelId, messageId, userIds, emoji, isBurst)
	return usersReacted
}
