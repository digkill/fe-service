package store

const (
	MISSING_CHANNEL_ERROR        = "store.sql_channel.get_by_name.missing.app_error"
	MISSING_CHANNEL_MEMBER_ERROR = "store.sql_channel.get_member.missing.app_error"
	CHANNEL_EXISTS_ERROR         = "store.sql_channel.save_channel.exists.app_error"

	MISSING_ACCOUNT_ERROR      = "store.sql_user.missing_account.const"
	MISSING_AUTH_ACCOUNT_ERROR = "store.sql_user.get_by_auth.missing_account.app_error"

	USER_SEARCH_OPTION_NAMES_ONLY              = "names_only"
	USER_SEARCH_OPTION_NAMES_ONLY_NO_FULL_NAME = "names_only_no_full_name"
	USER_SEARCH_OPTION_ALL_NO_FULL_NAME        = "all_no_full_name"
	USER_SEARCH_OPTION_ALLOW_INACTIVE          = "allow_inactive"

	FEATURE_TOGGLE_PREFIX = "feature_enabled_"

	INVITE_TOKEN_NOT_FOUND = "store.sql_token.get_user_invite_token.not_found"
)
