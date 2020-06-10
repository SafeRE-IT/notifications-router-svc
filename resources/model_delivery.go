/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Delivery struct {
	Key
	Attributes DeliveryAttributes `json:"attributes"`
}
type DeliveryResponse struct {
	Data     Delivery `json:"data"`
	Included Included `json:"included"`
}

type DeliveryListResponse struct {
	Data     []Delivery `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustDelivery - returns Delivery from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustDelivery(key Key) *Delivery {
	var delivery Delivery
	if c.tryFindEntry(key, &delivery) {
		return &delivery
	}
	return nil
}
