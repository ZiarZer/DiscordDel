package crawl

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

type Crawler struct {
	Sdk *discord.DiscordSdk
}

func (crawler *Crawler) CrawlChannel(authorizationToken string, authorIds []string, channelId string) {
	channel := crawler.Sdk.GetChannel(channelId, authorizationToken)
	if channel == nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get channel %s", channelId), utils.ERROR)
		return
	}

	crawlingInfo, err := crawler.Sdk.Repo.GetChannelCrawlingInfo(channelId)
	if err != nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get crawling info for channel %s", channelId), utils.WARNING)
	}

	if channel.Type == int(types.GuildForum) {
		crawler.crawlChannelThreads(authorizationToken, channel, crawlingInfo)
	} else {
		crawler.crawlChannelMessages(authorizationToken, channel, authorIds, crawlingInfo)
	}
}
