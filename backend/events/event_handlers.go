package events

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

// EventHandlerManager manages all event handlers
type EventHandlerManager struct {
	eventBus EventBus
}

// NewEventHandlerManager creates a new event handler manager
func NewEventHandlerManager(eventBus EventBus) *EventHandlerManager {
	em := &EventHandlerManager{
		eventBus: eventBus,
	}

	// Register all handlers
	em.registerHandlers()

	return em
}

// registerHandlers registers all event handlers with the event bus
func (em *EventHandlerManager) registerHandlers() {
	// User events
	em.eventBus.RegisterHandler(EventTypeUserCreated, em.handleUserCreated)
	em.eventBus.RegisterHandler(EventTypeUserLogin, em.handleUserLogin)
	em.eventBus.RegisterHandler(EventTypeUserLogout, em.handleUserLogout)

	// Auth events
	em.eventBus.RegisterHandler(EventTypeAuthSuccess, em.handleAuthSuccess)
	em.eventBus.RegisterHandler(EventTypeAuthFailure, em.handleAuthFailure)

	// Role events
	em.eventBus.RegisterHandler(EventTypeRoleAssigned, em.handleRoleAssigned)
	em.eventBus.RegisterHandler(EventTypeRoleRemoved, em.handleRoleRemoved)

	// Organization events
	em.eventBus.RegisterHandler(EventTypeUserAddedToOrg, em.handleUserAddedToOrg)
	em.eventBus.RegisterHandler(EventTypeUserRemovedFromOrg, em.handleUserRemovedFromOrg)

	// Admin events
	em.eventBus.RegisterHandler(EventTypeAdminAction, em.handleAdminAction)

	// System events
	em.eventBus.RegisterHandler(EventTypeSystemStartup, em.handleSystemStartup)
	em.eventBus.RegisterHandler(EventTypeSystemError, em.handleSystemError)
}

// handleUserCreated handles user creation events
func (em *EventHandlerManager) handleUserCreated(ctx context.Context, event Event) error {
	log.Printf("Handling user created event: %s for user %v", event.ID, event.Data[DataKeyUserID])

	// Here you could:
	// - Send welcome email
	// - Create default profile
	// - Assign default role
	// - Log to audit system
	// - Update analytics

	// Simulate some async work
	time.Sleep(100 * time.Millisecond)

	log.Printf("User created event processed successfully: %s", event.ID)
	return nil
}

// handleUserLogin handles user login events
func (em *EventHandlerManager) handleUserLogin(ctx context.Context, event Event) error {
	log.Printf("Handling user login event: %s for user %v", event.ID, event.Data[DataKeyUserID])

	// Here you could:
	// - Update last login timestamp
	// - Log login attempt
	// - Send security notification if suspicious
	// - Update user session count

	// Simulate some async work
	time.Sleep(50 * time.Millisecond)

	log.Printf("User login event processed successfully: %s", event.ID)
	return nil
}

// handleUserLogout handles user logout events
func (em *EventHandlerManager) handleUserLogout(ctx context.Context, event Event) error {
	log.Printf("Handling user logout event: %s for user %v", event.ID, event.Data[DataKeyUserID])

	// Here you could:
	// - Update session statistics
	// - Log logout time
	// - Clean up temporary data

	// Simulate some async work
	time.Sleep(30 * time.Millisecond)

	log.Printf("User logout event processed successfully: %s", event.ID)
	return nil
}

// handleAuthSuccess handles successful authentication events
func (em *EventHandlerManager) handleAuthSuccess(ctx context.Context, event Event) error {
	log.Printf("Handling auth success event: %s for user %v", event.ID, event.Data[DataKeyUserID])

	// Here you could:
	// - Update authentication metrics
	// - Log successful auth
	// - Update user activity tracking

	// Simulate some async work
	time.Sleep(25 * time.Millisecond)

	log.Printf("Auth success event processed successfully: %s", event.ID)
	return nil
}

// handleAuthFailure handles failed authentication events
func (em *EventHandlerManager) handleAuthFailure(ctx context.Context, event Event) error {
	log.Printf("Handling auth failure event: %s for user %v", event.ID, event.Data[DataKeyUserID])

	// Here you could:
	// - Update failure metrics
	// - Check for brute force attempts
	// - Send security alerts
	// - Log failed attempt details

	// Simulate some async work
	time.Sleep(40 * time.Millisecond)

	log.Printf("Auth failure event processed successfully: %s", event.ID)
	return nil
}

// handleRoleAssigned handles role assignment events
func (em *EventHandlerManager) handleRoleAssigned(ctx context.Context, event Event) error {
	log.Printf("Handling role assigned event: %s for user %v, role %v",
		event.ID, event.Data[DataKeyUserID], event.Data[DataKeyRoleID])

	// Here you could:
	// - Update user permissions cache
	// - Send notification to user
	// - Log role change
	// - Update audit trail

	// Simulate some async work
	time.Sleep(60 * time.Millisecond)

	log.Printf("Role assigned event processed successfully: %s", event.ID)
	return nil
}

// handleRoleRemoved handles role removal events
func (em *EventHandlerManager) handleRoleRemoved(ctx context.Context, event Event) error {
	log.Printf("Handling role removed event: %s for user %v, role %v",
		event.ID, event.Data[DataKeyUserID], event.Data[DataKeyRoleID])

	// Here you could:
	// - Update user permissions cache
	// - Send notification to user
	// - Log role change
	// - Update audit trail

	// Simulate some async work
	time.Sleep(60 * time.Millisecond)

	log.Printf("Role removed event processed successfully: %s", event.ID)
	return nil
}

// handleUserAddedToOrg handles user added to organization events
func (em *EventHandlerManager) handleUserAddedToOrg(ctx context.Context, event Event) error {
	log.Printf("Handling user added to org event: %s for user %v, org %v",
		event.ID, event.Data[DataKeyUserID], event.Data[DataKeyOrgID])

	// Here you could:
	// - Send welcome email to organization
	// - Update organization member count
	// - Assign default organization role
	// - Log membership change

	// Simulate some async work
	time.Sleep(70 * time.Millisecond)

	log.Printf("User added to org event processed successfully: %s", event.ID)
	return nil
}

// handleUserRemovedFromOrg handles user removed from organization events
func (em *EventHandlerManager) handleUserRemovedFromOrg(ctx context.Context, event Event) error {
	log.Printf("Handling user removed from org event: %s for user %v, org %v",
		event.ID, event.Data[DataKeyUserID], event.Data[DataKeyOrgID])

	// Here you could:
	// - Send goodbye email
	// - Update organization member count
	// - Remove organization-specific permissions
	// - Log membership change

	// Simulate some async work
	time.Sleep(70 * time.Millisecond)

	log.Printf("User removed from org event processed successfully: %s", event.ID)
	return nil
}

// handleAdminAction handles admin action events
func (em *EventHandlerManager) handleAdminAction(ctx context.Context, event Event) error {
	log.Printf("Handling admin action event: %s for user %v, action: %s",
		event.ID, event.Data[DataKeyUserID], event.Data[DataKeyAction])

	// Here you could:
	// - Log admin action to audit system
	// - Send notifications to other admins
	// - Update admin activity metrics
	// - Check for suspicious admin actions

	// Simulate some async work
	time.Sleep(80 * time.Millisecond)

	log.Printf("Admin action event processed successfully: %s", event.ID)
	return nil
}

// handleSystemStartup handles system startup events
func (em *EventHandlerManager) handleSystemStartup(ctx context.Context, event Event) error {
	log.Printf("Handling system startup event: %s", event.ID)

	// Here you could:
	// - Initialize system components
	// - Load configuration
	// - Start background services
	// - Send startup notifications

	// Simulate some async work
	time.Sleep(200 * time.Millisecond)

	log.Printf("System startup event processed successfully: %s", event.ID)
	return nil
}

// handleSystemError handles system error events
func (em *EventHandlerManager) handleSystemError(ctx context.Context, event Event) error {
	log.Printf("Handling system error event: %s: %v", event.ID, event.Data[DataKeyError])

	// Here you could:
	// - Send error notifications
	// - Log to error tracking system
	// - Trigger error recovery procedures
	// - Update system health metrics

	// Simulate some async work
	time.Sleep(100 * time.Millisecond)

	log.Printf("System error event processed successfully: %s", event.ID)
	return nil
}

// LogEvent logs an event to the system log
func (em *EventHandlerManager) LogEvent(event Event) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	log.Printf("Event logged: %s", string(eventJSON))
}
