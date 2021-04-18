package srvcont

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFilePath(t *testing.T) {
	assert.Equal(t,
		"acme",
		GetBucketDomainNameFromHost("acme.sepet.devingen.io"),
		"incorrect domain",
	)
	assert.Equal(t,
		"localhost",
		GetBucketDomainNameFromHost("localhost"),
		"incorrect domain",
	)
}
