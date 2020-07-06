/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type NotificatorService struct {
	Key
	Attributes NotificatorServiceAttributes `json:"attributes"`
}
type NotificatorServiceResponse struct {
	Data     NotificatorService `json:"data"`
	Included Included           `json:"included"`
}

type NotificatorServiceListResponse struct {
	Data     []NotificatorService `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
}

// MustNotificatorService - returns NotificatorService from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustNotificatorService(key Key) *NotificatorService {
	var notificatorService NotificatorService
	if c.tryFindEntry(key, &notificatorService) {
		return &notificatorService
	}
	return nil
}
