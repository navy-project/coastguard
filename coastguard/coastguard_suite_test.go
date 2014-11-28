package coastguard_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "testing"
    )

func TestCoastguard(t *testing.T) {
  RegisterFailHandler(Fail)
    RunSpecs(t, "Coastguard Suite")
}

