/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type NotificationAttributes struct {
	Channel      *string   `json:"channel,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	Message      Message   `json:"message"`
	Priority     int32     `json:"priority"`
	ScheduledFor time.Time `json:"scheduled_for"`
	Token        *string   `json:"token,omitempty"`
	Topic        string    `json:"topic"`
}
