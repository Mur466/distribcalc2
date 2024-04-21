package signin

import (
	"time"

	"github.com/Mur466/distribcalc2/internal/cfg"
	"github.com/Mur466/distribcalc2/internal/entities"
	"github.com/Mur466/distribcalc2/internal/errors"
	l "github.com/Mur466/distribcalc2/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repo
	cfg  *cfg.Config
}

func NewSvc(repo Repo, cfg *cfg.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *Service) Do(request *request) (string, error) {

	if request.Username == "" || request.Password == "" {
		return "", errors.ErrInvalidCreds
	}
	u := s.repo.GetUser(request.Username)
	if u == nil {
		return "", errors.ErrInvalidCreds
	}
	if !s.checkAuth(request, u) {
		return "", errors.ErrInvalidCreds
	}
	tokenString, err := s.GetJWTtoken(u.Username)
	if err != nil {
		l.SLogger.Error("token generation failure: %w", err)
		return "", errors.ErrInvalidCreds
	}

	return tokenString, nil
}

func (s *Service) checkAuth(creds *request, user *entities.User) bool {
	// todo: переделать на хеш
	return creds.Username == user.Username && bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)) == nil
}

func (s *Service) GetJWTtoken(username string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"nbf":      now.Unix(),                                                 // время, с которого токен станет валидным, пример now.Add(1*time.Minute).Unix()
		"exp":      now.Add(time.Duration(s.cfg.AuthTTL) * time.Second).Unix(), // время, с которого токен перестанет быть валидным ("протухнет")
		"iat":      now.Unix(),                                                 // время создания токена
	})

	tokenString, err := token.SignedString([]byte(s.cfg.Secret))
	return tokenString, err

	/*  пример проверки токена
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		l.SLogger.Fatal(err)
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		fmt.Println("user name: ", claims["username"])
	} else {
		panic(err)
	}
	*/
}
