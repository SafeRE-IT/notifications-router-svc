/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type ManageOfferRequest struct {
	Key
	Attributes ManageOfferRequestAttributes `json:"attributes"`
}
type ManageOfferRequestResponse struct {
	Data     ManageOfferRequest `json:"data"`
	Included Included           `json:"included"`
}

type ManageOfferRequestListResponse struct {
	Data     []ManageOfferRequest `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
}

// MustManageOfferRequest - returns ManageOfferRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustManageOfferRequest(key Key) *ManageOfferRequest {
	var manageOfferRequest ManageOfferRequest
	if c.tryFindEntry(key, &manageOfferRequest) {
		return &manageOfferRequest
	}
	return nil
}
