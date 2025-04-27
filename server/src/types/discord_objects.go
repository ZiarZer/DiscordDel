package types

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	GlobalName    string `json:"global_name"`
	Avatar        string `json:"avatar"`
}

type Guild struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Channel struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	LastMessageId string `json:"last_message_id"`
	Type          int    `json:"type"`
	ParentId      string `json:"parent_id"`
	GuildId       string `json:"guild_id"`
	MessageCount  int    `json:"message_count"`
}

type Message struct {
	Id        string `json:"id"`
	Content   string `json:"content"`
	Type      int    `json:"type"`
	ChannelId string `json:"channel_id"`
	Author    User   `json:"author"`
}
