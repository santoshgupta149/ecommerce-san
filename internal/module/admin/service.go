package admin

import (
	"context"
	"errors"
	"log"
	"time"

	"ecommerce-go/internal/module/admin/dto"

	"github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	repo AdminRepository
}

func NewService(repo AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// CreateAdmin contains business logic — no HTTP, no gin here.
func (s *AdminService) CreateAdmin(ctx context.Context, req dto.CreateAdminRequest) (dto.CreateAdminResponse, error) {
	log.Printf("[admin] service: creating admin first_name=%q last_name=%q email=%q", req.FirstName, req.LastName, req.Email)

	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.CreateAdminResponse{}, err
	}
	if existing != nil {
		return dto.CreateAdminResponse{}, errors.New("email already registered")
	}

	existingMobile, err := s.repo.FindByMobile(ctx, req.Mobile)
	if err != nil {
		return dto.CreateAdminResponse{}, err
	}
	if existingMobile != nil {
		return dto.CreateAdminResponse{}, errors.New("mobile number already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[admin] service: failed to hash password err=%v", err)
		return dto.CreateAdminResponse{}, errors.New("failed to process password")
	}

	now := time.Now().UTC()
	admin := &Admin{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Mobile:    req.Mobile,
		Password:  string(hashedPassword),
		Role:      "admin",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreateAdmin(ctx, admin); err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return dto.CreateAdminResponse{}, errors.New("email or mobile already registered")
		}
		log.Printf("[admin] service: failed to create admin err=%v", err)
		return dto.CreateAdminResponse{}, errors.New("failed to create admin")
	}

	log.Printf("[admin] service: admin created id=%d email=%q", admin.ID, admin.Email)

	return dto.CreateAdminResponse{
		ID:        admin.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}, nil
}
