package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"GoHeadless/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Claims represents the JWT payload
type Claims struct {
	UserID      string `json:"user_id"`
	RoleID      string `json:"role_id"`
	IsSuperuser bool   `json:"is_superuser"`
	jwt.RegisteredClaims
}

// LoginRequest is the payload for a login action
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse is what is returned after a successful login
type LoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

// RegisterRequest is the payload for creating a new user
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleID   string `json:"role_id" validate:"required"`
}

// Service is the interface for the auth service
type Service interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*domain.User, error)
	ValidateToken(tokenStr string) (*Claims, error)
	
	// Admin Operations
	GetAllRoles(ctx context.Context) ([]domain.Role, error)
	CreateRole(ctx context.Context, role *domain.Role) error
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id string) error
	
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	DeleteUser(ctx context.Context, id string) error

	GetStats(ctx context.Context) (map[string]int64, error)
}

type service struct {
	repo        Repository
	collRepo    domain.CollectionRepository
	recordRepo  domain.RecordRepository
	secret      string
}

func NewService(repo Repository, collRepo domain.CollectionRepository, recordRepo domain.RecordRepository) Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme-set-jwt-secret-in-env"
	}
	return &service{
		repo:       repo,
		collRepo:   collRepo,
		recordRepo: recordRepo,
		secret:     secret,
	}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.ActiveStatus {
		return nil, errors.New("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{Token: token, User: *user}, nil
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	roleID, err := primitive.ObjectIDFromHex(req.RoleID)
	if err != nil {
		return nil, errors.New("invalid role_id")
	}

	user := &domain.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		RoleID:       roleID.Hex(),
		ActiveStatus: true,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	return claims, nil
}

func (s *service) generateToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID:      user.ID.Hex(),
		RoleID:      user.RoleID,
		IsSuperuser: user.IsInitialAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *service) GetAllRoles(ctx context.Context) ([]domain.Role, error) {
	return s.repo.GetAllRoles(ctx)
}

func (s *service) CreateRole(ctx context.Context, role *domain.Role) error {
	if role.Name == "" {
		return errors.New("role name is required")
	}
	return s.repo.CreateRole(ctx, role)
}

func (s *service) UpdateRole(ctx context.Context, role *domain.Role) error {
	if role.ID.IsZero() {
		return errors.New("role id is required for update")
	}
	return s.repo.UpdateRole(ctx, role)
}

func (s *service) DeleteRole(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid role id")
	}
	return s.repo.DeleteRole(ctx, oid)
}

func (s *service) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].PasswordHash = "" 
	}
	return users, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user id")
	}
	return s.repo.DeleteUser(ctx, oid)
}

func (s *service) GetStats(ctx context.Context) (map[string]int64, error) {
	userCount, _ := s.repo.CountUsers(ctx)
	roleCount, _ := s.repo.CountRoles(ctx)
	collCount, _ := s.collRepo.CountCollections(ctx)
	
	colls, _ := s.collRepo.FindAll(ctx)
	var recordCount int64
	for _, c := range colls {
		count, _ := s.recordRepo.CountRecords(ctx, c.Slug)
		recordCount += count
	}

	return map[string]int64{
		"users":       userCount,
		"roles":       roleCount,
		"collections": collCount,
		"records":     recordCount,
	}, nil
}
