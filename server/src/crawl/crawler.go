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

func (crawler *Crawler) CrawlChannel(ctx context.Context, authorIds []types.Snowflake, channelId types.Snowflake) error {
	channel, err := crawler.Sdk.GetChannel(ctx, channelId)
	if err != nil {
		return err
	}
	if channel == nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get channel %s", channelId), utils.ERROR)
		return nil
	}

	crawlingInfo, err := crawler.Sdk.Repo.GetChannelCrawlingInfo(channelId)
	if err != nil {
		crawler.Sdk.Log(fmt.Sprintf("Failed to get crawling info for channel %s", channelId), utils.WARNING)
	}

	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(fmt.Sprintf("Crawl channel %s", channelId), crawler.Sdk.Log, true, true),
	)
	if channel.Type == types.GuildForum {
		err = crawler.crawlChannelThreads(ctx, channel, crawlingInfo)
	} else {
		err = crawler.crawlChannelMessages(ctx, channel, authorIds, crawlingInfo)
	}
	return err
}
