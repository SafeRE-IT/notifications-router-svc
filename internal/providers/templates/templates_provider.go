package templates

type TemplatesProvider interface {
	GetTemplate(topic, channel, locale string) ([]byte, error)
}
