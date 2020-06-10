/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type CreatePaymentRequestOp struct {
	Key
	Attributes    CreatePaymentRequestOpAttributes    `json:"attributes"`
	Relationships CreatePaymentRequestOpRelationships `json:"relationships"`
}
type CreatePaymentRequestOpResponse struct {
	Data     CreatePaymentRequestOp `json:"data"`
	Included Included               `json:"included"`
}

type CreatePaymentRequestOpListResponse struct {
	Data     []CreatePaymentRequestOp `json:"data"`
	Included Included                 `json:"included"`
	Links    *Links                   `json:"links"`
}

// MustCreatePaymentRequestOp - returns CreatePaymentRequestOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreatePaymentRequestOp(key Key) *CreatePaymentRequestOp {
	var createPaymentRequestOp CreatePaymentRequestOp
	if c.tryFindEntry(key, &createPaymentRequestOp) {
		return &createPaymentRequestOp
	}
	return nil
}
