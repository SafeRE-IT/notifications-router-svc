package resources

type DeliveryStatus string

const (
	DeliveryStatusNotSent  DeliveryStatus = "not_sent"
	DeliveryStatusFailed   DeliveryStatus = "failed"
	DeliveryStatusSent     DeliveryStatus = "sent"
	DeliveryStatusCanceled DeliveryStatus = "canceled"
	DeliveryStatusSkipped  DeliveryStatus = "skipped"
)

type NotificationPriority int32

const (
	NotificationsPriorityLowest NotificationPriority = iota + 1
	NotificationsPriorityLow
	NotificationsPriorityMedium
	NotificationsPriorityHigh
	NotificationsPriorityHighest
)
