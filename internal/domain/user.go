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

func (u *User) IsGoogleUser() bool {
	return u.GoogleID != ""
}

func (u *User) IsEmailUser() bool {
	return u.PasswordHash != ""
}

/*
 * CREATE TABLE users (
 *     id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 *     name           TEXT NOT NULL,
 *     email          TEXT NOT NULL UNIQUE,
 *     password_hash  TEXT,
 *     google_id      TEXT UNIQUE,
 *     avatar_url     TEXT,
 *     email_verified BOOLEAN NOT NULL DEFAULT FALSE,
 *     created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
 *     updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
 *     deleted_at     TIMESTAMPTZ
 * );
 *
 * CREATE INDEX idx_users_deleted_at ON users (deleted_at);
 */
