package usecase

import "github.com/secmon-lab/alertchain/pkg/domain/interfaces"

type UseCase struct {
	db    interfaces.Database
	genAI interfaces.GenAI
}

func New(options ...Option) *UseCase {
	uc := &UseCase{}
	for _, opt := range options {
		opt(uc)
	}
	return uc
}

type Option func(*UseCase)

func WithDatabase(db interfaces.Database) Option {
	return func(uc *UseCase) {
		uc.db = db
	}
}

func WithGenAI(genAI interfaces.GenAI) Option {
	return func(uc *UseCase) {
		uc.genAI = genAI
	}
}
