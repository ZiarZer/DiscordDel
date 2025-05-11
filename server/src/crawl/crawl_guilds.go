package crawl

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/types"
)

func (crawler *Crawler) CrawlGuild(authorizationToken string, authorIds []types.Snowflake, guildId types.Snowflake) {
	defer actions.StartAction(fmt.Sprintf("Crawl guild %s", guildId), crawler.Sdk.Log).EndAction()
	channels := crawler.Sdk.GetGuildChannels(guildId, authorizationToken)
	for i := range channels {
		crawler.CrawlChannel(authorizationToken, authorIds, channels[i].Id)
	}
}

func (crawler *Crawler) CrawlAllGuilds(authorizationToken string, authorIds []types.Snowflake) {
	defer actions.StartAction("Crawl all guilds", crawler.Sdk.Log).EndAction()
	guilds := crawler.Sdk.GetUserGuilds(authorizationToken)
	for i := range guilds {
		crawler.CrawlGuild(authorizationToken, authorIds, guilds[i].Id)
	}
}
