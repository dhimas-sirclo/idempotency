package faker

import (
	"idempotency/model"

	"github.com/go-faker/faker/v4"
)

func Product() (res model.UpsertProductPayload) {
	faker.FakeData(&res)
	return res
}
