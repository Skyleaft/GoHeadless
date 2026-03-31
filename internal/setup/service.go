package setup

import (
	"context"
	"log"

	"GoHeadless/internal/auth"
	"GoHeadless/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Service handles the system initialization check
type Service struct {
	repo auth.Repository
}

func NewService(repo auth.Repository) *Service {
	return &Service{repo: repo}
}

// IsSetupRequired returns true if no users exist in the database
func (s *Service) IsSetupRequired(ctx context.Context) (bool, error) {
	count, err := s.repo.CountUsers(ctx)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// Bootstrap seeds the default roles and creates the initial superadmin user
func (s *Service) Bootstrap(ctx context.Context, username, password string) (*domain.User, error) {
	// Create default roles: Superadmin, Admin, Editor, Viewer
	superadminRole, err := s.seedRole(ctx, "Superadmin", "Full wildcard access to everything", []domain.Permission{
		{CollectionSlug: "*", Actions: []string{"create", "read", "update", "delete"}},
	})
	if err != nil {
		return nil, err
	}

	// Seed standard roles (best-effort, ignore duplicates)
	s.seedRole(ctx, "Admin", "Full CRUD on all collections", []domain.Permission{ //nolint:errcheck
		{CollectionSlug: "*", Actions: []string{"create", "read", "update", "delete"}},
	})
	s.seedRole(ctx, "Editor", "Create, read, and update — no delete", []domain.Permission{ //nolint:errcheck
		{CollectionSlug: "*", Actions: []string{"create", "read", "update"}},
	})
	s.seedRole(ctx, "Viewer", "Read-only access", []domain.Permission{ //nolint:errcheck
		{CollectionSlug: "*", Actions: []string{"read"}},
	})

	log.Println("[Setup] Default roles seeded successfully.")

	// Hash the initial superadmin password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:       username,
		PasswordHash:   string(hash),
		RoleID:         superadminRole.ID.Hex(),
		IsInitialAdmin: true,
		ActiveStatus:   true,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	log.Printf("[Setup] Initial superadmin '%s' created successfully.", username)
	return user, nil
}

func (s *Service) seedRole(ctx context.Context, name, description string, perms []domain.Permission) (*domain.Role, error) {
	// Check if role already exists
	existing, err := s.repo.GetRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	role := &domain.Role{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		Permissions: perms,
	}
	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}
