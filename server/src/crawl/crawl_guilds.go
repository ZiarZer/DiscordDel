package crawl

import (
	"context"
	"fmt"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (crawler *Crawler) crawlGuildChannels(ctx context.Context, authorIds []types.Snowflake, guildId types.Snowflake) error {
	action := actions.NewMajorAction(utils.CRAWL, utils.GUILD, &guildId, fmt.Sprintf("Crawl guild %s", guildId))
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(action, crawler.Sdk.Log, true),
	)
	channels, err := crawler.Sdk.GetGuildChannels(ctx, guildId)
	if err != nil {
		return err
	}
	for i := range channels {
		err = crawler.CrawlChannel(ctx, authorIds, channels[i].Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (crawler *Crawler) crawlGuildBySearch(ctx context.Context, authorIds []types.Snowflake, guildId types.Snowflake) error {
	action := actions.NewMajorAction(utils.CRAWL, utils.GUILD, &guildId, fmt.Sprintf("Crawl guild %s by search", guildId))
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(action, crawler.Sdk.Log, true),
	)

	var minId *types.Snowflake
	crawlingInfo, err := crawler.Sdk.Repo.GetCrawlingInfo(guildId)
	if err != nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get crawling info for guild %s", guildId), utils.WARNING)
	}
	if crawlingInfo != nil {
		crawler.Sdk.Log(fmt.Sprintf("Found crawling info, starting from %s", crawlingInfo.NewestReadId), utils.INFO)
		minId = &crawlingInfo.NewestReadId
	}
	options := discord.SearchGuildMessagesOptions{
		AuthorIds: authorIds,
		SortBy:    utils.MakePointer("timestamp"),
		SortOrder: utils.MakePointer("asc"),
		MinId:     minId,
	}
	messages, err := crawler.fetchGuildMessagesBySearch(ctx, guildId, &options)
	if err != nil {
		return err
	}
	for len(messages) > 0 {
		if options.Offset == 9975 {
			options.Offset = 0
			options.MinId = &messages[len(messages)-1].Id
		} else {
			options.Offset = min(options.Offset+len(messages), 9975)
		}
		messages, err = crawler.fetchGuildMessagesBySearch(ctx, guildId, &options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (crawler *Crawler) fetchGuildMessagesBySearch(ctx context.Context, guildId types.Snowflake, options *discord.SearchGuildMessagesOptions) ([]types.Message, error) {
	messages, err := crawler.Sdk.SearchGuildMessages(ctx, guildId, options)
	if err != nil {
		return nil, err
	}

	crawler.Sdk.Repo.InsertMultipleMessages(messages, "PENDING")
	crawler.storeGuildMessagesCrawlingInfo(guildId, messages, options)
	return messages, nil
}

func (crawler *Crawler) storeGuildMessagesCrawlingInfo(guildId types.Snowflake, fetchedMessages []types.Message, fetchOptions *discord.SearchGuildMessagesOptions) {
	if fetchOptions == nil || fetchOptions.MinId == nil && fetchOptions.Offset == 0 {
		if len(fetchedMessages) > 0 {
			crawler.Sdk.Repo.InsertCrawlingInfo(
				guildId,
				"GUILD",
				fetchedMessages[0].Id,
				fetchedMessages[len(fetchedMessages)-1].Id,
				true,
			)
		}
	} else {
		var newestReadMessageId *types.Snowflake
		if fetchOptions.MinId != nil && len(fetchedMessages) > 0 {
			newestReadMessageId = &fetchedMessages[len(fetchedMessages)-1].Id
		}
		crawler.Sdk.Repo.UpdateCrawlingInfo(
			guildId,
			nil,
			newestReadMessageId,
			nil,
		)
	}
}
