package events

// Event types for authentication
const (
	// User events
	EventTypeUserCreated         = "user.created"
	EventTypeUserUpdated         = "user.updated"
	EventTypeUserDeleted         = "user.deleted"
	EventTypeUserLogin           = "user.login"
	EventTypeUserLogout          = "user.logout"
	EventTypeUserPasswordChanged = "user.password_changed"

	// Authentication events
	EventTypeAuthSuccess      = "auth.success"
	EventTypeAuthFailure      = "auth.failure"
	EventTypeAuthTokenRefresh = "auth.token_refresh"
	EventTypeAuthTokenExpired = "auth.token_expired"

	// Role events
	EventTypeRoleAssigned = "role.assigned"
	EventTypeRoleRemoved  = "role.removed"
	EventTypeRoleCreated  = "role.created"
	EventTypeRoleUpdated  = "role.updated"
	EventTypeRoleDeleted  = "role.deleted"

	// Organization events
	EventTypeOrgCreated         = "organization.created"
	EventTypeOrgUpdated         = "organization.updated"
	EventTypeOrgDeleted         = "organization.deleted"
	EventTypeUserAddedToOrg     = "organization.user_added"
	EventTypeUserRemovedFromOrg = "organization.user_removed"

	// Admin events
	EventTypeAdminAction = "admin.action"
	EventTypeAdminLogin  = "admin.login"
	EventTypeAdminLogout = "admin.logout"

	// System events
	EventTypeSystemStartup  = "system.startup"
	EventTypeSystemShutdown = "system.shutdown"
	EventTypeSystemError    = "system.error"
	EventTypeSystemWarning  = "system.warning"
)

// Event topics for publishing
const (
	TopicAuth          = "auth"
	TopicUsers         = "users"
	TopicRoles         = "roles"
	TopicOrganizations = "organizations"
	TopicAdmin         = "admin"
	TopicSystem        = "system"
	TopicAll           = "all" // Broadcast to all topics
)

// Event data keys
const (
	DataKeyUserID    = "user_id"
	DataKeyEmail     = "email"
	DataKeyName      = "name"
	DataKeyRoleID    = "role_id"
	DataKeyRoleName  = "role_name"
	DataKeyOrgID     = "organization_id"
	DataKeyOrgName   = "organization_name"
	DataKeyAction    = "action"
	DataKeyDetails   = "details"
	DataKeyIPAddress = "ip_address"
	DataKeyUserAgent = "user_agent"
	DataKeyTimestamp = "timestamp"
	DataKeyError     = "error"
	DataKeySuccess   = "success"
)

// Common event data builders
func BuildUserEventData(userID int, email, name string) map[string]interface{} {
	return map[string]interface{}{
		DataKeyUserID: userID,
		DataKeyEmail:  email,
		DataKeyName:   name,
	}
}

func BuildRoleEventData(roleID int, roleName string) map[string]interface{} {
	return map[string]interface{}{
		DataKeyRoleID:   roleID,
		DataKeyRoleName: roleName,
	}
}

func BuildOrgEventData(orgID int, orgName string) map[string]interface{} {
	return map[string]interface{}{
		DataKeyOrgID:   orgID,
		DataKeyOrgName: orgName,
	}
}

func BuildAuthEventData(userID int, email, action string, success bool) map[string]interface{} {
	return map[string]interface{}{
		DataKeyUserID:  userID,
		DataKeyEmail:   email,
		DataKeyAction:  action,
		DataKeySuccess: success,
	}
}

func BuildAdminEventData(userID int, action, details string) map[string]interface{} {
	return map[string]interface{}{
		DataKeyUserID:  userID,
		DataKeyAction:  action,
		DataKeyDetails: details,
	}
}
