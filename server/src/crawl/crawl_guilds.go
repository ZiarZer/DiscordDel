package crawl

import (
	"context"
	"fmt"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (crawler *Crawler) CrawlGuild(ctx context.Context, authorIds []types.Snowflake, guildId types.Snowflake) error {
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

func (crawler *Crawler) CrawlAllGuilds(ctx context.Context, authorIds []types.Snowflake) error {
	action := actions.NewMajorAction(utils.CRAWL, utils.ALL, nil, "Crawl all guilds")
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(action, crawler.Sdk.Log, true),
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
