package jwt

import (
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/walnuts1018/mpeg-dash-encoder/config"
)

func TestJWT(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	RegisterFailHandler(Fail)
	RunSpecs(t, "JWT Suite")
}

var _ = Describe("JWT", func() {
	JwtSigningKey := config.JWTSigningKey("signingKey")
	manager := NewManager(JwtSigningKey)
	fakeManager := NewManager(config.JWTSigningKey("fakeSigningKey"))

	entityIDsA := []string{
		"1",
		"2",
		"3",
	}

	entityIDsEmpty := []string{}

	It("Normal", func() {
		By("Create User Token")
		token, err := manager.CreateUserToken(entityIDsA)
		Expect(err).NotTo(HaveOccurred())

		By("Get Media IDs From Token")
		mediaIDs, err := manager.GetMediaIDsFromToken(token)
		Expect(err).NotTo(HaveOccurred())
		Expect(mediaIDs).To(Equal(entityIDsA))
	})

	It("Empty Media IDs", func() {
		By("Create User Token")
		token, err := manager.CreateUserToken(entityIDsEmpty)
		Expect(err).NotTo(HaveOccurred())

		By("Get Media IDs From Token")
		mediaIDs, err := manager.GetMediaIDsFromToken(token)
		Expect(err).NotTo(HaveOccurred())
		Expect(mediaIDs).To(Equal([]string{}))
	})

	It("Invalid Token", func() {
		By("Create User Token with Fake Manager")
		token, err := fakeManager.CreateUserToken(entityIDsA)
		Expect(err).NotTo(HaveOccurred())

		By("Check Fake Token can be parsed by Fake Manager")
		mediaIDs, err := fakeManager.GetMediaIDsFromToken(token)
		Expect(err).NotTo(HaveOccurred())
		Expect(mediaIDs).To(Equal(entityIDsA))

		By("Check Fake Token will be rejected by Real Manager")
		_, err = manager.GetMediaIDsFromToken(token)
		Expect(err).To(HaveOccurred())
		slog.Debug("error", slog.Any("error", err))
	})
})
