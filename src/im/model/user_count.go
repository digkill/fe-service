package model

// Options for counting users
type UserCountOptions struct {
	// Should include users that are bots
	IncludeBotAccounts bool
	// Should include deleted users (of any type)
	IncludeDeleted bool
	// Exclude regular users
	ExcludeRegularUsers bool
	// Only include users on a specific team. "" for any team.
	TeamId string
}
