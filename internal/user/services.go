package user

import (
    "context"
    "errors"
    "time"
    
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type Service struct {
    db     *gorm.DB
    logger Logger
}

type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Debug(msg string, fields ...interface{})
}

func NewService(db *gorm.DB, logger Logger) *Service {
    return &Service{
        db:     db,
        logger: logger,
    }
}

func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Validate password confirmation
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