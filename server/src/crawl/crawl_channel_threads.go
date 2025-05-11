package crawl

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (crawler *Crawler) crawlChannelThreads(authorizationToken string, mainChannel *types.Channel, crawlingInfo *types.CrawlingInfo) {
	defer actions.StartAction(fmt.Sprintf("Crawl threads in channel %s", mainChannel.Id), crawler.Sdk.Log).EndAction()
	pageSize := 25
	options := discord.SearchChannelThreadsOptions{
		SortBy:    utils.MakePointer("creation_time"),
		SortOrder: utils.MakePointer("asc"),
		Limit:     &pageSize,
	}

	var startOffset int
	if crawlingInfo == nil {
		threads := crawler.fetchChannelThreads(authorizationToken, mainChannel.Id, &options)
		if len(threads) == 0 {
			crawler.Sdk.Log("Channel doesn't have any thread: nothing to do", utils.INFO)
			return
		}
		crawler.Sdk.Repo.InsertChannelCrawlingInfo(mainChannel.Id, threads[0].Id, threads[len(threads)-1].Id, true)
		startOffset = len(threads)
	} else {
		alreadyCrawledCount, _ := crawler.Sdk.Repo.GetChannelChildrenCount(mainChannel.Id)
		startOffset = alreadyCrawledCount
		options.Offset = startOffset - pageSize
		threads := crawler.fetchChannelThreads(authorizationToken, mainChannel.Id, &options)
		for len(threads) == 0 && options.Offset > 0 {
			options.Offset = max(options.Offset-pageSize, 0)
			threads = crawler.fetchChannelThreads(authorizationToken, mainChannel.Id, &options)
		}
		if len(threads) > 0 {
			updatedNewestReadId := threads[len(threads)-1].Id
			for len(threads) > 0 && options.Offset > 0 && utils.GetTimestampFromSnowflake(threads[0].Id) >= utils.GetTimestampFromSnowflake(crawlingInfo.NewestReadId) {
				options.Offset = max(options.Offset-pageSize, 0)
				threads = crawler.fetchChannelThreads(authorizationToken, mainChannel.Id, &options)
			}
			crawler.Sdk.Repo.UpdateChannelCrawlingNewestReadId(mainChannel.Id, updatedNewestReadId)
		}
	}
	options.Offset = startOffset
	threads := crawler.fetchChannelThreads(authorizationToken, mainChannel.Id, &options)
	for len(threads) > 0 {
		options.Offset += len(threads)
		crawler.Sdk.Repo.UpdateChannelCrawlingNewestReadId(mainChannel.Id, threads[len(threads)-1].Id)
		threads = crawler.fetchChannelThreads(authorizationToken, mainChannel.Id, &options)
	}
}

func (crawler *Crawler) fetchChannelThreads(authorizationToken string, mainChannelId types.Snowflake, options *discord.SearchChannelThreadsOptions) []types.Channel {
	threadChannels := crawler.Sdk.SearchChannelThreads(authorizationToken, mainChannelId, options)
	crawler.Sdk.Repo.InsertMultipleChannels(threadChannels)
	return threadChannels
}
