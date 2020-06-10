/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type CreateNotificationAttributes struct {
	Channel      *string    `json:"channel,omitempty"`
	Message      Message    `json:"message"`
	Priority     *int32     `json:"priority,omitempty"`
	ScheduledFor *time.Time `json:"scheduled_for,omitempty"`
	Token        *string    `json:"token,omitempty"`
	Topic        string     `json:"topic"`
}
