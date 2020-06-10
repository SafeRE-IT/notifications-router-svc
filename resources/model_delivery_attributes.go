/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type DeliveryAttributes struct {
	Destination     string     `json:"destination"`
	DestinationType string     `json:"destination_type"`
	SentAt          *time.Time `json:"sent_at,omitempty"`
	Status          string     `json:"status"`
}
