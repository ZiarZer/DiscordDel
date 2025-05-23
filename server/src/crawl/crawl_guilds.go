package crawl

import (
	"context"
	"fmt"

	"github.com/ZiarZer/DiscordDel/types"
)

func (crawler *Crawler) CrawlGuild(ctx context.Context, authorIds []types.Snowflake, guildId types.Snowflake) error {
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(fmt.Sprintf("Crawl guild %s", guildId), crawler.Sdk.Log, true, true),
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

func (crawler *Crawler) CrawlAllGuilds(ctx context.Context, authorIds []types.Snowflake) error {
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction("Crawl all guilds", crawler.Sdk.Log, true, true),
	)
	guilds, err := crawler.Sdk.GetUserGuilds(ctx)
	if err != nil {
		return err
	}
	for i := range guilds {
		err = crawler.CrawlGuild(ctx, authorIds, guilds[i].Id)
		if err != nil {
			return err
		}
	}
	return nil
}
