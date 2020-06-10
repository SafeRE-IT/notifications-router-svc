/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type CreatePaymentRequest struct {
	Key
	Attributes    CreatePaymentRequestAttributes    `json:"attributes"`
	Relationships CreatePaymentRequestRelationships `json:"relationships"`
}
type CreatePaymentRequestResponse struct {
	Data     CreatePaymentRequest `json:"data"`
	Included Included             `json:"included"`
}

type CreatePaymentRequestListResponse struct {
	Data     []CreatePaymentRequest `json:"data"`
	Included Included               `json:"included"`
	Links    *Links                 `json:"links"`
}

// MustCreatePaymentRequest - returns CreatePaymentRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreatePaymentRequest(key Key) *CreatePaymentRequest {
	var createPaymentRequest CreatePaymentRequest
	if c.tryFindEntry(key, &createPaymentRequest) {
		return &createPaymentRequest
	}
	return nil
}
