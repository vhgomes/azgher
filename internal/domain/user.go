package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuthProvider string

const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
)

type User struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Name  string    `json:"name" db:"name"`
	Email string    `json:"email" db:"email"`

	// Campo de senha caso o usuário seja cadastrado por email, ficará vazio caso seja por Google OAuth
	PasswordHash string `json:"-" db:"password_hash"`

	// Google OAuth
	GoogleID  string `json:"google_id,omitempty" db:"google_id"`
	AvatarURL string `json:"avatar_url,omitempty" db:"avatar_url"`

	// Provider(s) usados — útil se quiser permitir múltiplos vínculos
	// Não persistido diretamente; calculado ou vindo de outra tabela/query
	Providers []AuthProvider `json:"providers" db:"-"`

	EmailVerified bool `json:"email_verified" db:"email_verified"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

func NewUser(
	id uuid.UUID,
	name, email string,
	passwordHash, googleID, avatarURL *string,
	emailVerified bool,
	createdAt, updatedAt time.Time,
) *User {
	user := &User{
		ID:            id,
		Name:          name,
		Email:         email,
		EmailVerified: emailVerified,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	if passwordHash != nil {
		user.PasswordHash = *passwordHash
		user.Providers = append(user.Providers, AuthProviderEmail)
	}
	if googleID != nil {
		user.GoogleID = *googleID
		user.Providers = append(user.Providers, AuthProviderGoogle)
	}
	if avatarURL != nil {
		user.AvatarURL = *avatarURL
	}

	return user
}

func (u *User) IsGoogleUser() bool {
	return u.GoogleID != ""
}

func (u *User) IsEmailUser() bool {
	return u.PasswordHash != ""
}
