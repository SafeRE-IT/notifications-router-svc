/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateNotification struct {
	Key
	Attributes    CreateNotificationAttributes    `json:"attributes"`
	Relationships CreateNotificationRelationships `json:"relationships"`
}
type CreateNotificationRequest struct {
	Data     CreateNotification `json:"data"`
	Included Included           `json:"included"`
}

type CreateNotificationListRequest struct {
	Data     []CreateNotification `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
}

// MustCreateNotification - returns CreateNotification from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateNotification(key Key) *CreateNotification {
	var createNotification CreateNotification
	if c.tryFindEntry(key, &createNotification) {
		return &createNotification
	}
	return nil
}
