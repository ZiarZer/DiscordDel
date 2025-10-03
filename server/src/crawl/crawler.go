package crawl

import (
	"context"
	"fmt"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

type Crawler struct {
	Sdk          *discord.DiscordSdk
	ActionLogger *actions.ActionLogger
}

func (crawler *Crawler) CrawlAllGuilds(ctx context.Context, authorIds []types.Snowflake, shouldCrawlReactions bool) error {
	action := actions.NewMajorAction(utils.CRAWL, utils.ALL, nil, "Crawl all guilds")
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(action, crawler.Sdk.Log, true),
	)
	guilds, err := crawler.Sdk.GetUserGuilds(ctx)
	if err != nil {
		return err
	}
	for i := range guilds {
		err = crawler.CrawlGuild(ctx, authorIds, guilds[i].Id, shouldCrawlReactions)
		if err != nil {
			return err
		}
	}
	return nil
}

func (crawler *Crawler) CrawlGuild(ctx context.Context, authorIds []types.Snowflake, guildId types.Snowflake, shouldCrawlReactions bool) error {
	if shouldCrawlReactions {
		return crawler.crawlGuildChannels(ctx, authorIds, guildId)
	}
	return crawler.crawlGuildBySearch(ctx, authorIds, guildId)
}

func (crawler *Crawler) CrawlChannel(ctx context.Context, authorIds []types.Snowflake, channelId types.Snowflake) error {
	channel, err := crawler.Sdk.GetChannel(ctx, channelId)
	if err != nil {
		return err
	}
	if channel == nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get channel %s", channelId), utils.ERROR)
		return nil
	}

	crawlingInfo, err := crawler.Sdk.Repo.GetCrawlingInfo(channelId)
	if err != nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get crawling info for channel %s", channelId), utils.WARNING)
	}

	action := actions.NewMajorAction(utils.CRAWL, utils.CHANNEL, &channelId, fmt.Sprintf("Crawl channel %s", channelId))
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(action, crawler.Sdk.Log, true),
	)
	if channel.Type == types.GuildForum {
		err = crawler.crawlChannelThreads(ctx, channel, crawlingInfo)
	} else {
		err = crawler.crawlChannelMessages(ctx, channel, authorIds, crawlingInfo)
	}
	return err
}
