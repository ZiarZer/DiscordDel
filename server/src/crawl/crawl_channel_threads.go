package crawl

import (
	"context"
	"fmt"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (crawler *Crawler) crawlChannelThreads(ctx context.Context, mainChannel *types.Channel, crawlingInfo *types.CrawlingInfo) error {
	defer crawler.ActionLogger.EndAction(
		crawler.ActionLogger.StartAction(fmt.Sprintf("Crawl threads in channel %s", mainChannel.Id), crawler.Sdk.Log, true, false),
	)
	pageSize := 25
	options := discord.SearchChannelThreadsOptions{
		SortBy:    utils.MakePointer("creation_time"),
		SortOrder: utils.MakePointer("asc"),
		Limit:     &pageSize,
	}

	var startOffset int
	if crawlingInfo == nil {
		threads, err := crawler.fetchChannelThreads(ctx, mainChannel.Id, &options)
		if err != nil {
			return err
		} else if len(threads) == 0 {
			crawler.Sdk.Log("Channel doesn't have any thread: nothing to do", utils.INFO)
			return nil
		}
		crawler.Sdk.Repo.InsertChannelCrawlingInfo(mainChannel.Id, threads[0].Id, threads[len(threads)-1].Id, true)
		startOffset = len(threads)
	} else {
		alreadyCrawledCount, _ := crawler.Sdk.Repo.GetChannelChildrenCount(mainChannel.Id)
		startOffset = alreadyCrawledCount
		options.Offset = startOffset - pageSize
		threads, err := crawler.fetchChannelThreads(ctx, mainChannel.Id, &options)
		if err != nil {
			return err
		}
		for len(threads) == 0 && options.Offset > 0 {
			options.Offset = max(options.Offset-pageSize, 0)
			threads, err = crawler.fetchChannelThreads(ctx, mainChannel.Id, &options)
			if err != nil {
				return err
			}
		}
		if len(threads) > 0 {
			updatedNewestReadId := threads[len(threads)-1].Id
			for len(threads) > 0 && options.Offset > 0 && utils.GetTimestampFromSnowflake(threads[0].Id) >= utils.GetTimestampFromSnowflake(crawlingInfo.NewestReadId) {
				options.Offset = max(options.Offset-pageSize, 0)
				threads, err = crawler.fetchChannelThreads(ctx, mainChannel.Id, &options)
				if err != nil {
					return err
				}
			}
			crawler.Sdk.Repo.UpdateChannelCrawlingNewestReadId(mainChannel.Id, updatedNewestReadId)
		}
	}
	options.Offset = startOffset
	threads, err := crawler.fetchChannelThreads(ctx, mainChannel.Id, &options)
	if err != nil {
		return err
	}
	for len(threads) > 0 {
		options.Offset += len(threads)
		crawler.Sdk.Repo.UpdateChannelCrawlingNewestReadId(mainChannel.Id, threads[len(threads)-1].Id)
		threads, err = crawler.fetchChannelThreads(ctx, mainChannel.Id, &options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (crawler *Crawler) fetchChannelThreads(ctx context.Context, mainChannelId types.Snowflake, options *discord.SearchChannelThreadsOptions) ([]types.Channel, error) {
	threadChannels, err := crawler.Sdk.SearchChannelThreads(ctx, mainChannelId, options)
	if err != nil {
		return nil, err
	}
	crawler.Sdk.Repo.InsertMultipleChannels(threadChannels)
	return threadChannels, nil
}
