package contracts

import (
	contextcontracts "project-survey-proceeder/internal/context/contracts"
)

type IServiceProvider interface {
	GetContextFiller() contextcontracts.IRequestFiller
}
