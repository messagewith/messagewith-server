package auth

import "errors"

var (
	ErrUserAlreadyLoggedIn         = errors.New("user is already logged in")
	ErrUserNotLoggedIn             = errors.New("user is not logged in")
	ErrUserBadPassword             = errors.New("bad password")
	ErrUserNicknameAlreadyUsed     = errors.New("user with specified nickname already exists")
	ErrUserEmailAlreadyUsed        = errors.New("user with specified e-mail already exists")
	ErrNoUserWithSpecifiedEmail    = errors.New("there is no user with the specified e-mail")
	ErrNoUserWithSpecifiedId       = errors.New("there is no user with specified id")
	ErrNoUserWithSpecifiedNickname = errors.New("there is no user with specified nickname")
	ErrInvalidID                   = errors.New("invalid id")
)
