/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Notification struct {
	Key
	Attributes    NotificationAttributes    `json:"attributes"`
	Relationships NotificationRelationships `json:"relationships"`
}
type NotificationResponse struct {
	Data     Notification `json:"data"`
	Included Included     `json:"included"`
}

type NotificationListResponse struct {
	Data     []Notification `json:"data"`
	Included Included       `json:"included"`
	Links    *Links         `json:"links"`
}

// MustNotification - returns Notification from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustNotification(key Key) *Notification {
	var notification Notification
	if c.tryFindEntry(key, &notification) {
		return &notification
	}
	return nil
}
