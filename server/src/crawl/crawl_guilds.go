package crawl

import "github.com/ZiarZer/DiscordDel/types"

func (crawler *Crawler) CrawlGuild(authorizationToken string, authorIds []types.Snowflake, guildId types.Snowflake) {
	channels := crawler.Sdk.GetGuildChannels(guildId, authorizationToken)
	for i := range channels {
		crawler.CrawlChannel(authorizationToken, authorIds, channels[i].Id)
	}
}

func (crawler *Crawler) CrawlAllGuilds(authorizationToken string, authorIds []types.Snowflake) {
	guilds := crawler.Sdk.GetUserGuilds(authorizationToken)
	for i := range guilds {
		crawler.CrawlGuild(authorizationToken, authorIds, guilds[i].Id)
	}
}
