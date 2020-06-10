/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type RemoveAssetOp struct {
	Key
	Relationships RemoveAssetOpRelationships `json:"relationships"`
}
type RemoveAssetOpResponse struct {
	Data     RemoveAssetOp `json:"data"`
	Included Included      `json:"included"`
}

type RemoveAssetOpListResponse struct {
	Data     []RemoveAssetOp `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustRemoveAssetOp - returns RemoveAssetOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRemoveAssetOp(key Key) *RemoveAssetOp {
	var removeAssetOp RemoveAssetOp
	if c.tryFindEntry(key, &removeAssetOp) {
		return &removeAssetOp
	}
	return nil
}
