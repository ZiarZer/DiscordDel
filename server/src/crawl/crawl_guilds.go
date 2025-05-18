package crawl

import (
	"context"
	"fmt"

	"github.com/ZiarZer/DiscordDel/types"
)

func (crawler *Crawler) CrawlGuild(ctx context.Context, authorIds []types.Snowflake, guildId types.Snowflake) {
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(fmt.Sprintf("Crawl guild %s", guildId), crawler.Sdk.Log, true, true),
	)
	channels := crawler.Sdk.GetGuildChannels(ctx, guildId)
	for i := range channels {
		crawler.CrawlChannel(ctx, authorIds, channels[i].Id)
	}
}

func (crawler *Crawler) CrawlAllGuilds(ctx context.Context, authorIds []types.Snowflake) {
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction("Crawl all guilds", crawler.Sdk.Log, true, true),
	)
	guilds := crawler.Sdk.GetUserGuilds(ctx)
	for i := range guilds {
		crawler.CrawlGuild(ctx, authorIds, guilds[i].Id)
	}
}
