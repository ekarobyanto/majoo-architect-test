package repository_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommentRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Comments Repository Suite")
}
