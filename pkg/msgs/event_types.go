package msgs

type EventType string

const (
	EventVersion EventType = "1.0"

	EventStatus EventType = "statusbar:status"
	EventError  EventType = "statusbar:error"

	EventManager         EventType = "manager:change"
	EventProjectsUpdated EventType = "projects:updated"

	EventAppInit    EventType = "app:initialized"
	EventAppReady   EventType = "app:ready"
	EventViewChange EventType = "app:view-changed"
	EventRefresh    EventType = "app:refresh"

	EventDataLoaded EventType = "data:loaded"

	EventTabCycle EventType = "hotkey:tab-cycle"

	EventImagesChanged EventType = "images:changed"
)

var AllMessages = []struct {
	Value  EventType `json:"value"`
	TSName string    `json:"tsname"`
}{
	{EventStatus, "STATUS"},
	{EventError, "ERROR"},
	{EventManager, "MANAGER"},
	{EventProjectsUpdated, "PROJECTS_UPDATED"},
	{EventAppInit, "APP_INIT"},
	{EventAppReady, "APP_READY"},
	{EventViewChange, "VIEW_CHANGE"},
	{EventRefresh, "REFRESH"},
	{EventVersion, "VERSION"},
	{EventDataLoaded, "DATA_LOADED"},
	{EventTabCycle, "TAB_CYCLE"},
	{EventImagesChanged, "IMAGES_CHANGED"},
}
