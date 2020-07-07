package settings

type SettingsProvider interface {
	IsTopicEnabled(accountId, topic string) (bool, error)
	GetChannels(accountId string) ([]string, error)
	GetLocale(accountId string) (*string, error)
}
