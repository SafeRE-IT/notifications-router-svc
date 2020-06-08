/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type CreateNotificationAttributes struct {
	DeliveryType *string    `json:"delivery_type,omitempty"`
	Locale       *string    `json:"locale,omitempty"`
	Message      Message    `json:"message"`
	Priority     *int32     `json:"priority,omitempty"`
	SendTime     *time.Time `json:"send_time,omitempty"`
	Token        *string    `json:"token,omitempty"`
	Topic        string     `json:"topic"`
}
