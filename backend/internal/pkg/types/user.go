package types

import "github.com/google/uuid"

type UserID uuid.UUID

func (i UserID) UUID() uuid.UUID {
	return uuid.UUID(i)
}

func (i UserID) String() string {
	return uuid.UUID(i).String()
}
