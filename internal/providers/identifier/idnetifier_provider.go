package identifier

type IdentifierProvider interface {
	GetIdentifierByChannel(channel string, accountId string) (*Identifier, error)
}
