package keeper_test

import (
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("dymension msg_server.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress
	var nonOwner i.TestValidatorAddress
	_ = nonOwner

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		nonOwner = i.GenerateTestValidatorAddress("NonOwner")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("AnnounceValidator - DYMENSION KAS - (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x9695e09597f3111b183700e06d6f1a7d50ea1aee"                       // looks ok
		validatorPrivKey := "c18908a1bbe0ec588cd6522d2b02af3076a2f2c562a09bb8bf5a40f6e9a0ef1b" // looks busted

		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId,
			Creator:         creator.Address,
		})

		// Assert
		Expect(err).To(BeNil())
		validateAnnouncement(s, mailboxId.String(), validatorAddress, []string{
			storageLocation,
		})
	})
})
