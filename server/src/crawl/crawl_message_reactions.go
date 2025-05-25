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

func (crawler *Crawler) crawlMessageReactions(ctx context.Context, message *types.Message, authorIds []types.Snowflake) error {
	action := actions.NewAction(fmt.Sprintf("Crawl reactions on message %s", message.Id))
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(action, crawler.Sdk.TempLog, false),
	)
	if message.Reactions == nil {
		return nil
	}

	for i := range message.Reactions {
		reaction := message.Reactions[i]
		emoji := *reaction.Emoji.Name
		if reaction.Emoji.Id != nil {
			emoji = emoji + fmt.Sprintf(":%s", *reaction.Emoji.Id)
		}
		if reaction.CountDetails.Normal > 0 {
			err := crawler.crawlReactionsOnEmoji(ctx, message.ChannelId, message.Id, emoji, false, authorIds)
			if err != nil {
				return err
			}
		}
		if reaction.CountDetails.Burst > 0 {
			err := crawler.crawlReactionsOnEmoji(ctx, message.ChannelId, message.Id, emoji, true, authorIds)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (crawler *Crawler) crawlReactionsOnEmoji(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake, emoji string, isBurst bool, authorIds []types.Snowflake) error {
	action := actions.NewAction(fmt.Sprintf("Crawl reactions with %s on message %s", emoji, messageId))
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(action, crawler.Sdk.TempLog, false),
	)
	usersReacted, err := crawler.fetchReactionsOnEmoji(ctx, channelId, messageId, emoji, isBurst, nil, authorIds)
	if err != nil {
		return err
	}
	pageSize := 100
	options := &discord.GetMessageReactionsOptions{Limit: &pageSize}

	for len(usersReacted) == pageSize {
		options.After = &usersReacted[len(usersReacted)-1].Id
		usersReacted, err = crawler.fetchReactionsOnEmoji(ctx, channelId, messageId, emoji, isBurst, options, authorIds)
		if err != nil {
			return err
		}
	}
	return nil
}

func (crawler *Crawler) fetchReactionsOnEmoji(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake, emoji string, isBurst bool, options *discord.GetMessageReactionsOptions, authorIds []types.Snowflake) ([]types.User, error) {
	usersReacted, err := crawler.Sdk.GetMessageReactions(ctx, channelId, messageId, emoji, isBurst, options)
	if err != nil {
		return nil, err
	}

	userIds := utils.MapWithoutDuplicate(
		utils.Filter(
			usersReacted,
			func(user types.User) bool { return slices.Contains(authorIds, user.Id) },
		),
		func(user types.User) types.Snowflake { return user.Id },
	)
	crawler.Sdk.Repo.InsertMultipleReactions(channelId, messageId, userIds, emoji, isBurst)
	return usersReacted, nil
}
