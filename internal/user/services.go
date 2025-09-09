package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
    db     *gorm.DB
    logger common.Logger
}

func NewService(db *gorm.DB, logger common.Logger) *Service {
    return &Service{
        db:     db,
        logger: logger,
    }
}

func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    if req.Password != req.ConfirmPassword {
        return nil, errors.New("passwords do not match")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        Email:     req.Email,
        Password:  string(hashedPassword),
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Phone:     req.Phone,
        IsActive:  true,
        Role:      RoleMember,
    }
    
    if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
        return nil, err
    }
    
    s.logger.Info("User created successfully", "user_id", user.ID)
    return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
    var user User
    if err := s.db.WithContext(ctx).Preload("Family").First(&user, "id = ?", userID).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, req *UpdateUserRequest) (*User, error) {
    var user User
    if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
        return nil, err
    }
    
    // Update fields
    if req.FirstName != "" {
        user.FirstName = req.FirstName
    }
    if req.LastName != "" {
        user.LastName = req.LastName
    }
    if req.Phone != "" {
        user.Phone = req.Phone
    }
    if req.DateOfBirth != "" {
        dob, err := time.Parse("2006-01-02", req.DateOfBirth)
        if err != nil {
            return nil, errors.New("invalid date format")
        }
        user.DateOfBirth = &dob
    }
    
    if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (s *Service) GetFamilyMembers(ctx context.Context, userID uuid.UUID) ([]User, error) {
    var user User
    if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
        return nil, err
    }
    
    if user.FamilyID == nil {
        return []User{}, nil
    }
    
    var members []User
    if err := s.db.WithContext(ctx).Where("family_id = ?", user.FamilyID).Find(&members).Error; err != nil {
        return nil, err
    }
    
    return members, nil
}

func (s *Service) CreateAddress(ctx context.Context, userID uuid.UUID, req *CreateAddressRequest) (*User, error) {
    var user User
    if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
        return nil, err
    }

    address := Address{
        Country:     req.Country,
        City:        req.City,
        PostalCode:  req.PostalCode,
        Street:      req.Street,
        HouseNumber: req.HouseNumber,
    }

    user.Address = append(user.Address, address)

    if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
        return nil, err
    }

    return &user, nil
}

func (s *Service) UpdateAddress(ctx context.Context, userID uuid.UUID, req *UpdateAddressRequest) (*User, error) {
    var user User
    if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
        return nil, err
    }

    if len(user.Address) == 0 {
        return nil, errors.New("no address to update")
    }
    // Find the address to update via the address id
    var address *Address
    for i := range user.Address {
        if user.Address[i].ID == req.ID {
            address = &user.Address[i]
            break
        }
    }
    if address == nil {
        return nil, errors.New("address not found")
    }

    if req.Country != "" {
        address.Country = req.Country
    }
    if req.City != "" {
        address.City = req.City
    }
    if req.PostalCode != "" {
        address.PostalCode = req.PostalCode
    }
    if req.Street != "" {
        address.Street = req.Street
    }
    if req.HouseNumber != 0 {
        address.HouseNumber = req.HouseNumber
    }

    if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
        return nil, err
    }

    return &user, nil
}

func (s *Service) ListAllUsers(ctx context.Context) ([]User, error) {
    var users []User
    if err := s.db.WithContext(ctx).Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (s *Service) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
    return s.db.WithContext(ctx).
        Model(&User{}).
        Where("id = ?", userID).
        Update("is_active", false).
        Error
}

func (s *Service) ActivateUser(ctx context.Context, userID uuid.UUID) error {
    return s.db.WithContext(ctx).
        Model(&User{}).
        Where("id = ?", userID).
        Update("is_active", true).
        Error
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
    // Implementation for forgot password logic
    return nil
}

func (s *Service) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
    // Implementation for reset password logic
    return nil
}

func (s *Service) GetAllActiveUsers(ctx context.Context) ([]User, error) {
    var users []User
    if err := s.db.WithContext(ctx).Where("is_active = ?", true).Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (s *Service) GetAllInactiveUsers(ctx context.Context) ([]User, error) {
    var users []User
    if err := s.db.WithContext(ctx).Where("is_active = ?", false).Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
