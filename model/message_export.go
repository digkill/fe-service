package model

type MessageExport struct {
	TeamId          *string
	TeamName        *string
	TeamDisplayName *string

	ChannelId          *string
	ChannelName        *string
	ChannelDisplayName *string
	ChannelType        *string

	UserId    *string
	UserEmail *string
	Username  *string

	PostId         *string
	PostCreateAt   *int64
	PostMessage    *string
	PostType       *string
	PostRootId     *string
	PostOriginalId *string
	PostFileIds    StringArray
}
