package utils

import "github.com/gofrs/uuid"

func UUID() string {
	v4, _ := uuid.NewV4()
	return v4.String()
}
