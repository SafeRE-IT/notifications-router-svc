/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type CreateManageOfferRequestOp struct {
	Key
	Attributes    CreateManageOfferRequestOpAttributes    `json:"attributes"`
	Relationships CreateManageOfferRequestOpRelationships `json:"relationships"`
}
type CreateManageOfferRequestOpResponse struct {
	Data     CreateManageOfferRequestOp `json:"data"`
	Included Included                   `json:"included"`
}

type CreateManageOfferRequestOpListResponse struct {
	Data     []CreateManageOfferRequestOp `json:"data"`
	Included Included                     `json:"included"`
	Links    *Links                       `json:"links"`
}

// MustCreateManageOfferRequestOp - returns CreateManageOfferRequestOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateManageOfferRequestOp(key Key) *CreateManageOfferRequestOp {
	var createManageOfferRequestOp CreateManageOfferRequestOp
	if c.tryFindEntry(key, &createManageOfferRequestOp) {
		return &createManageOfferRequestOp
	}
	return nil
}
