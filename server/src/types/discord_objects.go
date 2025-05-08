package types

type ChannelType int

const (
	GuildText         ChannelType = iota // 0
	Dm                                   // 1
	GuildVoice                           // 2
	GroupDm                              // 3
	GuildCategory                        // 4
	GuildAnnouncement                    // 5
)
const (
	AnnouncementThread ChannelType = iota + 10 // 10
	PublicThread                               // 11
	PrivateThread                              // 12
	GuildStageVoice                            // 13
	GuildDirectory                             // 14
	GuildForum                                 // 15
	GuildMedia                                 // 16
)

type User struct {
	Id            string  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	GlobalName    *string `json:"global_name"`
	Avatar        *string `json:"avatar"`
}

type Guild struct {
	Id   string  `json:"id"`
	Name string  `json:"name"`
	Icon *string `json:"icon"`
}

type Channel struct {
	Id            string  `json:"id"`
	Name          *string `json:"name"`
	LastMessageId *string `json:"last_message_id"`
	Type          int     `json:"type"`
	ParentId      *string `json:"parent_id"`
	GuildId       *string `json:"guild_id"`
	MessageCount  *int    `json:"message_count"`
}

type Message struct {
	Id        string `json:"id"`
	Content   string `json:"content"`
	Type      int    `json:"type"`
	ChannelId string `json:"channel_id"`
	Author    User   `json:"author"`
}

type ThreadMember struct {
	Id     *string `json:"id"`
	UserId *string `json:"user_id"`
	Flags  int     `json:"flags"`
}

type ThreadsResult struct {
	Threads       []Channel      `json:"threads"`
	Members       []ThreadMember `json:"members"`
	FirstMessages []Message      `json:"first_messages"`
	HasMore       bool           `json:"has_more"`
}
