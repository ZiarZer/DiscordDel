package types

type Snowflake string

type User struct {
	Id            Snowflake `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	GlobalName    *string   `json:"global_name"`
	Avatar        *string   `json:"avatar"`
}

type Guild struct {
	Id   Snowflake `json:"id"`
	Name string    `json:"name"`
	Icon *string   `json:"icon"`
}

type Channel struct {
	Id            Snowflake   `json:"id"`
	Name          *string     `json:"name"`
	LastMessageId *Snowflake  `json:"last_message_id"`
	Type          ChannelType `json:"type"`
	ParentId      *Snowflake  `json:"parent_id"`
	GuildId       *Snowflake  `json:"guild_id"`
	MessageCount  *int        `json:"message_count"`
}

type Message struct {
	Id                  Snowflake                   `json:"id"`
	Content             string                      `json:"content"`
	Type                MessageType                 `json:"type"`
	ChannelId           Snowflake                   `json:"channel_id"`
	Author              User                        `json:"author"`
	InteractionMetadata *MessageInteractionMetadata `json:"interaction_metadata"`
	Reactions           []ReactionSummary           `json:"reactions"`
}

type ReactionSummary struct {
	Emoji        Emoji                `json:"emoji"`
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details"`
	Me           int                  `json:"me"`
	MeBurst      int                  `json:"me_burst"`
}

type Emoji struct {
	Id   *Snowflake `json:"id"`
	Name *string    `json:"name"`
}
type ReactionCountDetails struct {
	Burst  int `json:"burst"`
	Normal int `json:"normal"`
}

type Reaction struct {
	MessageId Snowflake `json:"id"`
	UserId    Snowflake `json:"content"`
	Emoji     string    `json:"type"`
	IsBurst   bool      `json:"is_burst"`
}

type ThreadMember struct {
	ThreadId *Snowflake `json:"id"`
	UserId   *Snowflake `json:"user_id"`
	Flags    int        `json:"flags"`
}

type ThreadsResult struct {
	Threads       []Channel      `json:"threads"`
	Members       []ThreadMember `json:"members"`
	FirstMessages []Message      `json:"first_messages"`
	HasMore       bool           `json:"has_more"`
}

type MessageInteractionMetadata struct {
	Id        Snowflake       `json:"id"`
	Type      InteractionType `json:"type"`
	Triggerer User            `json:"user"`
}

type ChannelType int
type MessageType int
type InteractionType int

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

const (
	Default              MessageType = iota // 0
	RecipientAdd                            // 1
	RecipientRemove                         // 2
	Call                                    // 3
	ChannelNameChange                       // 4
	ChannelIconChange                       // 5
	ChannelPinnedMessage                    // 6
	UserJoin                                // 7
	GuildBoost                              // 8
	GuildBoostTier1                         // 9
	GuildBoostTier2                         // 10
	GuildBoostTier3                         // 11
	ChannelFollowAdd                        // 12
)
const (
	GuildDiscoveryDisqualified              MessageType = iota + 14 // 14
	GuildDiscoveryRequalified                                       // 15
	GuildDiscoveryGracePeriodInitialWarning                         // 16
	GuildDiscoveryGracePeriodFinalWarning                           // 17
	ThreadCreated                                                   // 18
	Reply                                                           // 19
	ChatInputCommand                                                // 20
	ThreadStarterMessage                                            // 21
	GuildInviteReminder                                             // 22
	ContextMenuCommand                                              // 23
	AutoModerationAction                                            // 24
	RoleSubscriptionPurchase                                        // 25
	InteractionPremiumUpsell                                        // 26
	StageStart                                                      // 27
	StageEnd                                                        // 28
	StageSpeaker                                                    // 29
)
const (
	StageTopic                          MessageType = 31 // 31
	GuildApplicationPremiumSubscription MessageType = 32 // 32
)
const (
	GuildIncidentAlertModeEnabled  MessageType = iota + 36 // 36
	GuildIncidentAlertModeDisabled                         // 37
	GuildIncidentReportRaid                                // 38
	GuildIncidentReportFalseAlarm                          // 39
	PurchaseNotification           MessageType = 44        // 44
	PollResult                     MessageType = 46        // 46
)

const (
	Ping                           InteractionType = iota + 1 // 1
	ApplicationCommand                                        // 2
	MessageComponent                                          // 3
	ApplicationCommandAutocomplete                            // 4
	ModalSubmit                                               // 5
)
