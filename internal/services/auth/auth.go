package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/libs/jwt"
	"sso/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// для этого Сервисного слоя нужен Storage
type Auth struct {
	logger       *slog.Logger
	userSaver    IUserSaver
	userProvider IUserProvider
	appProvider  IAppProvider
	tokenTTL     time.Duration
}

// 1 конкретный интерфейс для данного сервиса
type IUserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

// Получает/достает/извлекает user
type IUserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type IAppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(logger *slog.Logger, userSaver IUserSaver, userProvider IUserProvider, appProvider IAppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		logger:       logger,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Register(ctx context.Context, email string, password string) (userId int64, err error) {
	const op = "auth.Register"

	logger := a.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	logger.Info("user is registering")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error("failed to generate password hash", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userId, err = a.userSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		logger.Error("failed to generate password hash", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("user has been registered")

	return userId, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string, appId int) (token string, err error) {
	const op = "auth.Login"

	logger := a.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	logger.Info("user is logging")

	// достаем пользователя по email
	user, err := a.userProvider.User(ctx, email)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.logger.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		logger.Error("failed to get user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	// идем проверять и сравнивать пароли
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		logger.Error("failed to decode password", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// достаем id приложения, в котором пользователь хочет залогинится
	app, err := a.appProvider.App(ctx, appId)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("user successfully logged in")

	// идем создавать токен в зависимости от ключа конкретного приложения
	// приватный ключ будет храниться как на стороне КЛИЕНТСКОГО приложения (UrlShortener) так и на стороне СЕРВЕРНОГО приложения (SSO)
	// секретный ключ (secret) будет храниться в Базе Данных
	jwtToken, err := jwt.GenerateNewToken(user, app, a.tokenTTL)

	if err != nil {
		logger.Error("failed to generate token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return jwtToken, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "auth.Login"

	logger := a.logger.With(
		slog.String("op", op),
		slog.Int64("user_id", userId),
	)

	logger.Info("cheching whether user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("checked whether user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
