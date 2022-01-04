package instructions

import (
	"encoding/base64"
	"fmt"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

func putReciepientArray(j *util.JSONBuilder, fieldName string, receipients []Recipient) {
	if len(receipients) == 0 {
		return
	}
	array := &util.JSONBuilder{}
	array.Encode.WriteRune('[')
	first := true
	for _, r := range receipients {
		if !first {
			array.Encode.WriteRune(',')
		}
		fmt.Fprintf(&array.Encode, `{"token":"%v","value":%v`, base64.StdEncoding.EncodeToString(r.Token[:]), r.Value)
	}
	array.Encode.WriteRune(']')
	j.PutString(fieldName, array.Encode.String())
}

func putTokenCiphersJSON(j *util.JSONBuilder, fieldName string, tc TokenCiphers) {
	if len(tc) == 0 {
		return
	}
	array := &util.JSONBuilder{}
	array.Encode.WriteRune('[')
	first := true
	for _, r := range tc {
		if !first {
			array.Encode.WriteRune(',')
		}
		fmt.Fprintf(&array.Encode, `{"token":"%v","cipher":%v`, base64.StdEncoding.EncodeToString(r.Token[:]), base64.StdEncoding.EncodeToString(r.Cipher))
	}
	array.Encode.WriteRune(']')
	j.PutString(fieldName, array.Encode.String())
}

func (a *AuthoredInstruction) JSON(kind byte, bulk *util.JSONBuilder) string {
	b := &util.JSONBuilder{}
	b.PutUint64("version", 0)
	b.PutUint64("instructionType", uint64(kind))
	b.PutUint64("epoch", a.epoch)
	b.PutHex("author", a.Author[:])
	fmt.Fprintf(&b.Encode, ",%v", bulk.Encode.String())
	if a.Wallet != crypto.ZeroToken {
		b.PutHex("wallet", a.Wallet[:])
	}
	b.PutUint64("fee", a.Fee)
	if a.Attorney != crypto.ZeroToken {
		b.PutHex("attorney", a.Attorney[:])
	}
	b.PutBase64("signature", a.signature[:])
	b.PutBase64("walletSignature", a.walletSignature[:])
	return b.ToString()
}

func (j *JoinNetwork) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutString("caption", j.Caption)
	bulk.PutJSON("details", j.Details)
	return j.Authored.JSON(IJoinNetwork, bulk)
}

func (j *UpdateInfo) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutJSON("details", j.Details)
	return j.Authored.JSON(IUpdateInfo, bulk)
}

func (j *GrantPowerOfAttorney) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("details", j.Attorney[:])
	return j.Authored.JSON(IGrantPowerOfAttorney, bulk)
}

func (j *RevokePowerOfAttorney) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("details", j.Attorney[:])
	return j.Authored.JSON(IRevokePowerOfAttorney, bulk)
}
func (j *CreateEphemeral) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("ephemeralToken", j.EphemeralToken[:])
	bulk.PutUint64("expiry", j.Expiry)
	return j.Authored.JSON(ICreateEphemeral, bulk)
}

func (j *SecureChannel) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("tokenRange", j.TokenRange)
	bulk.PutUint64("nonce", j.Nonce)
	bulk.PutHex("encryptedNonce", j.EncryptedNonce)
	bulk.PutHex("encryptedNonce", j.Content)
	return j.Authored.JSON(ISecureChannel, bulk)
}

func (j *CreateStage) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("audience", j.Audience[:])
	bulk.PutHex("submission", j.Submission[:])
	bulk.PutHex("moderation", j.Moderation[:])
	bulk.PutUint64("flag", uint64(j.Flag))
	bulk.PutString("description", j.Description)
	return j.Authored.JSON(ICreateAudience, bulk)
}

func (j *JoinStage) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("audience", j.Audience[:])
	bulk.PutString("presentation", j.Presentation)
	return j.Authored.JSON(IJoinAudience, bulk)
}

func (j *AcceptJoinStage) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("audience", j.Stage[:])
	bulk.PutHex("member", j.Member[:])
	bulk.PutHex("read", j.Read)
	bulk.PutHex("submit", j.Submit)
	bulk.PutHex("moderate", j.Moderate)
	bulk.PutBase64("modSignature", j.modSignature[:])
	return j.Authored.JSON(IAcceptJoinRequest, bulk)
}

func (j *UpdateStage) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("audience", j.Stage[:])
	bulk.PutHex("submission", j.Submission[:])
	bulk.PutHex("moderation", j.Moderation[:])
	bulk.PutUint64("flag", uint64(j.Flag))
	bulk.PutString("description", j.Description)
	putTokenCiphersJSON(bulk, "readMembers", j.ReadMembers)
	putTokenCiphersJSON(bulk, "subMembers", j.SubMembers)
	putTokenCiphersJSON(bulk, "modMembers", j.ModMembers)
	bulk.PutBase64("audSignature", j.audSignature[:])
	return j.Authored.JSON(IUpdateAudience, bulk)
}

func (j *SponsorshipOffer) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("audience", j.Stage[:])
	bulk.PutString("contentType", j.ContentType)
	bulk.PutBase64("content", j.Content)
	bulk.PutUint64("expiry", j.Expiry)
	bulk.PutUint64("revenue", j.Revenue)
	return j.Authored.JSON(ISponsorshipOffer, bulk)
}

func (j *SponsorshipAcceptance) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("audience", j.Stage[:])
	bulk.PutJSON("offer", j.Offer.JSON())
	bulk.PutBase64("modSignature", j.modSignature[:])
	return j.Authored.JSON(ISponsorshipAcceptance, bulk)
}

func (j *React) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutHex("hash", j.Hash)
	bulk.PutUint64("reaction", uint64(j.Reaction))
	return j.Authored.JSON(IReact, bulk)
}

func (j *Content) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutUint64("version", 0)
	bulk.PutUint64("instructionType", uint64(IContent))
	bulk.PutUint64("epoch", j.epoch)
	bulk.PutUint64("published", j.Published)
	bulk.PutHex("author", j.Author[:])
	bulk.PutHex("audience", j.Audience[:])
	bulk.PutString("contentType", j.ContentType)
	bulk.PutBase64("content", j.Content)
	bulk.PutHex("hash", j.Hash)
	if j.Wallet != crypto.ZeroToken {
		bulk.PutHex("wallet", j.Wallet[:])
	}
	bulk.PutUint64("fee", j.Fee)
	if j.Attorney != crypto.ZeroToken {
		bulk.PutHex("attorney", j.Attorney[:])
	}
	bulk.PutBase64("signature", j.Signature[:])
	bulk.PutBase64("walletSignature", j.WalletSignature[:])
	return bulk.ToString()
}

func (j *Transfer) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutUint64("version", 0)
	bulk.PutUint64("instructionType", uint64(ITransfer))
	bulk.PutUint64("epoch", j.epoch)
	bulk.PutHex("from", j.From[:])
	putReciepientArray(bulk, "to", j.To)
	bulk.PutString("reason", j.Reason)
	bulk.PutUint64("fee", j.Fee)
	bulk.PutBase64("signature", j.Signature[:])
	return bulk.ToString()
}

func (j *Deposit) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutUint64("version", 0)
	bulk.PutUint64("instructionType", uint64(IDeposit))
	bulk.PutUint64("epoch", j.epoch)
	bulk.PutHex("token", j.Token[:])
	bulk.PutUint64("value", j.Value)
	bulk.PutUint64("fee", j.Fee)
	bulk.PutBase64("signature", j.Signature[:])
	return bulk.ToString()
}

func (j *Withdraw) JSON() string {
	bulk := &util.JSONBuilder{}
	bulk.PutUint64("version", 0)
	bulk.PutUint64("instructionType", uint64(IDeposit))
	bulk.PutUint64("epoch", j.epoch)
	bulk.PutHex("token", j.Token[:])
	bulk.PutUint64("value", j.Value)
	bulk.PutUint64("fee", j.Fee)
	bulk.PutBase64("signature", j.Signature[:])
	return bulk.ToString()
}
