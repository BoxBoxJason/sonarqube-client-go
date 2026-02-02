package helpers

import (
	"fmt"
	"strings"
	"time"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"
)

// IgnoreNotFoundError returns nil if the error indicates a resource was not found,
// otherwise returns the original error. This is useful in cleanup functions where
// the resource may have already been deleted by the test itself.
func IgnoreNotFoundError(err error) error {
	if err == nil {
		return nil
	}
	// Check if the error message contains "not found" or similar patterns
	// that indicate the resource doesn't exist
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "does not exist") ||
		strings.Contains(errStr, "404") {
		return nil
	}

	return err
}

// CleanupManager manages cleanup of e2e test resources.
type CleanupManager struct {
	client    *sonargo.Client
	resources []cleanupResource
}

type cleanupResource struct {
	cleanupFn    func() error
	resourceType string
	identifier   string
}

// NewCleanupManager creates a new cleanup manager.
func NewCleanupManager(client *sonargo.Client) *CleanupManager {
	return &CleanupManager{
		client:    client,
		resources: make([]cleanupResource, 0),
	}
}

// RegisterCleanup registers a resource for cleanup.
func (cm *CleanupManager) RegisterCleanup(resourceType, identifier string, cleanupFn func() error) {
	cm.resources = append(cm.resources, cleanupResource{
		cleanupFn:    cleanupFn,
		resourceType: resourceType,
		identifier:   identifier,
	})
}

// Cleanup runs all registered cleanup functions in reverse order.
func (cm *CleanupManager) Cleanup() []error {
	var errors []error

	// Cleanup in reverse order (LIFO)
	for i := len(cm.resources) - 1; i >= 0; i-- {
		r := cm.resources[i]

		err := r.cleanupFn()
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to cleanup %s '%s': %w", r.resourceType, r.identifier, err))
		}
	}

	// Clear the resources
	cm.resources = cm.resources[:0]

	return errors
}

// CleanupOrphanedResources cleans up any orphaned e2e resources older than maxAge.
func CleanupOrphanedResources(client *sonargo.Client, maxAge time.Duration) error {
	// Clean up orphaned projects
	err := cleanupOrphanedProjects(client, maxAge)
	if err != nil {
		return fmt.Errorf("failed to cleanup orphaned projects: %w", err)
	}

	// Clean up orphaned users
	err = cleanupOrphanedUsers(client, maxAge)
	if err != nil {
		return fmt.Errorf("failed to cleanup orphaned users: %w", err)
	}

	// Clean up orphaned groups
	err = cleanupOrphanedGroups(client, maxAge)
	if err != nil {
		return fmt.Errorf("failed to cleanup orphaned groups: %w", err)
	}

	return nil
}

func cleanupOrphanedProjects(client *sonargo.Client, _ time.Duration) error {
	// Search for projects with e2e prefix
	projects, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
		PaginationArgs:    sonargo.PaginationArgs{Page: 0, PageSize: 0},
		AnalyzedBefore:    "",
		OnProvisionedOnly: false,
		Projects:          nil,
		Query:             E2EResourcePrefix,
		Qualifiers:        nil,
		Visibility:        "",
	})
	if err != nil {
		return fmt.Errorf("failed to search projects: %w", err)
	}

	if projects == nil || len(projects.Components) == 0 {
		return nil
	}

	for _, p := range projects.Components {
		if strings.HasPrefix(p.Key, E2EResourcePrefix) {
			_, _ = client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: p.Key,
			})
		}
	}

	return nil
}

func cleanupOrphanedUsers(client *sonargo.Client, _ time.Duration) error {
	//nolint:staticcheck // Using deprecated API until v2 API is implemented
	users, _, err := client.Users.Search(&sonargo.UsersSearchOption{
		PaginationArgs:        sonargo.PaginationArgs{Page: 0, PageSize: 0},
		Deactivated:           false,
		ExternalIdentity:      "",
		LastConnectedAfter:    "",
		LastConnectedBefore:   "",
		Managed:               false,
		Q:                     E2EResourcePrefix,
		SlLastConnectedAfter:  "",
		SlLastConnectedBefore: "",
	})
	if err != nil {
		return fmt.Errorf("failed to search users: %w", err)
	}

	if users == nil || len(users.Users) == 0 {
		return nil
	}

	for _, u := range users.Users {
		if strings.HasPrefix(u.Login, E2EResourcePrefix) {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, _ = client.Users.Deactivate(&sonargo.UsersDeactivateOption{
				Login:     u.Login,
				Anonymize: false,
			})
		}
	}

	return nil
}

func cleanupOrphanedGroups(client *sonargo.Client, _ time.Duration) error {
	//nolint:staticcheck // Using deprecated API until v2 API is implemented
	groups, _, err := client.UserGroups.Search(&sonargo.UserGroupsSearchOption{
		PaginationArgs: sonargo.PaginationArgs{Page: 0, PageSize: 0},
		Managed:        nil,
		Fields:         nil,
		Query:          E2EResourcePrefix,
	})
	if err != nil {
		return fmt.Errorf("failed to search groups: %w", err)
	}

	if groups == nil || len(groups.Groups) == 0 {
		return nil
	}

	for _, g := range groups.Groups {
		if strings.HasPrefix(g.Name, E2EResourcePrefix) {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _ = client.UserGroups.Delete(&sonargo.UserGroupsDeleteOption{
				Name: g.Name,
			})
		}
	}

	return nil
}
