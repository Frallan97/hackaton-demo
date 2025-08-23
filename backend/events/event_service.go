package events

// EventService provides a high-level interface for event operations
type EventService struct {
	eventBus EventBus
}

// NewEventService creates a new event service
func NewEventService(eventBus EventBus) *EventService {
	return &EventService{
		eventBus: eventBus,
	}
}

// PublishUserEvent publishes a user-related event
func (es *EventService) PublishUserEvent(eventType string, userID int, email, name string, additionalData map[string]interface{}) error {
	data := BuildUserEventData(userID, email, name)

	// Merge additional data
	for k, v := range additionalData {
		data[k] = v
	}

	return es.eventBus.Publish(TopicUsers, eventType, data, &userID)
}

// PublishAuthEvent publishes an authentication event
func (es *EventService) PublishAuthEvent(eventType string, userID int, email, action string, success bool, additionalData map[string]interface{}) error {
	data := BuildAuthEventData(userID, email, action, success)

	// Merge additional data
	for k, v := range additionalData {
		data[k] = v
	}

	return es.eventBus.Publish(TopicAuth, eventType, data, &userID)
}

// PublishRoleEvent publishes a role-related event
func (es *EventService) PublishRoleEvent(eventType string, userID int, roleID int, roleName string, additionalData map[string]interface{}) error {
	data := BuildRoleEventData(roleID, roleName)
	data[DataKeyUserID] = userID

	// Merge additional data
	for k, v := range additionalData {
		data[k] = v
	}

	return es.eventBus.Publish(TopicRoles, eventType, data, &userID)
}

// PublishOrgEvent publishes an organization-related event
func (es *EventService) PublishOrgEvent(eventType string, userID int, orgID int, orgName string, additionalData map[string]interface{}) error {
	data := BuildOrgEventData(orgID, orgName)
	data[DataKeyUserID] = userID

	// Merge additional data
	for k, v := range additionalData {
		data[k] = v
	}

	return es.eventBus.Publish(TopicOrganizations, eventType, data, &userID)
}

// PublishAdminEvent publishes an admin action event
func (es *EventService) PublishAdminEvent(userID int, action, details string, additionalData map[string]interface{}) error {
	data := BuildAdminEventData(userID, action, details)

	// Merge additional data
	for k, v := range additionalData {
		data[k] = v
	}

	return es.eventBus.Publish(TopicAdmin, EventTypeAdminAction, data, &userID)
}

// PublishSystemEvent publishes a system event
func (es *EventService) PublishSystemEvent(eventType string, data map[string]interface{}) error {
	return es.eventBus.Publish(TopicSystem, eventType, data, nil)
}

// SubscribeToTopic subscribes to a specific topic
func (es *EventService) SubscribeToTopic(topic string) (<-chan Event, error) {
	return es.eventBus.Subscribe(topic)
}

// SubscribeToUserEvents subscribes to all user events for a specific user
func (es *EventService) SubscribeToUserEvents(userID int) (<-chan Event, error) {
	// Subscribe to user-specific topic
	topic := TopicUsers
	return es.eventBus.Subscribe(topic)
}

// SubscribeToAuthEvents subscribes to authentication events
func (es *EventService) SubscribeToAuthEvents() (<-chan Event, error) {
	return es.eventBus.Subscribe(TopicAuth)
}

// SubscribeToAdminEvents subscribes to admin events
func (es *EventService) SubscribeToAdminEvents() (<-chan Event, error) {
	return es.eventBus.Subscribe(TopicAdmin)
}

// GetEventStats returns event bus statistics
func (es *EventService) GetEventStats() map[string]interface{} {
	return es.eventBus.GetEventStats()
}

// Shutdown gracefully shuts down the event service
func (es *EventService) Shutdown() {
	es.eventBus.Shutdown()
}

// PublishUserCreated publishes a user created event
func (es *EventService) PublishUserCreated(userID int, email, name string) error {
	return es.PublishUserEvent(EventTypeUserCreated, userID, email, name, nil)
}

// PublishUserLogin publishes a user login event
func (es *EventService) PublishUserLogin(userID int, email, name string) error {
	return es.PublishUserEvent(EventTypeUserLogin, userID, email, name, nil)
}

// PublishUserLogout publishes a user logout event
func (es *EventService) PublishUserLogout(userID int, email, name string) error {
	return es.PublishUserEvent(EventTypeUserLogout, userID, email, name, nil)
}

// PublishAuthSuccess publishes an authentication success event
func (es *EventService) PublishAuthSuccess(userID int, email, action string) error {
	return es.PublishAuthEvent(EventTypeAuthSuccess, userID, email, action, true, nil)
}

// PublishAuthFailure publishes an authentication failure event
func (es *EventService) PublishAuthFailure(userID int, email, action string, errorDetails string) error {
	data := map[string]interface{}{
		DataKeyError: errorDetails,
	}
	return es.PublishAuthEvent(EventTypeAuthFailure, userID, email, action, false, data)
}

// PublishRoleAssigned publishes a role assigned event
func (es *EventService) PublishRoleAssigned(userID int, roleID int, roleName string) error {
	return es.PublishRoleEvent(EventTypeRoleAssigned, userID, roleID, roleName, nil)
}

// PublishRoleRemoved publishes a role removed event
func (es *EventService) PublishRoleRemoved(userID int, roleID int, roleName string) error {
	return es.PublishRoleEvent(EventTypeRoleRemoved, userID, roleID, roleName, nil)
}

// PublishUserAddedToOrg publishes a user added to organization event
func (es *EventService) PublishUserAddedToOrg(userID int, orgID int, orgName string) error {
	return es.PublishOrgEvent(EventTypeUserAddedToOrg, userID, orgID, orgName, nil)
}

// PublishUserRemovedFromOrg publishes a user removed from organization event
func (es *EventService) PublishUserRemovedFromOrg(userID int, orgID int, orgName string) error {
	return es.PublishOrgEvent(EventTypeUserRemovedFromOrg, userID, orgID, orgName, nil)
}

// PublishSystemStartup publishes a system startup event
func (es *EventService) PublishSystemStartup() error {
	data := map[string]interface{}{
		DataKeyTimestamp: "system_started",
	}
	return es.PublishSystemEvent(EventTypeSystemStartup, data)
}

// PublishSystemError publishes a system error event
func (es *EventService) PublishSystemError(errorDetails string) error {
	data := map[string]interface{}{
		DataKeyError: errorDetails,
	}
	return es.PublishSystemEvent(EventTypeSystemError, data)
}
