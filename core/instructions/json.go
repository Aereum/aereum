package instructions

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type jsonBuilder struct {
	encode strings.Builder
}

func (j *jsonBuilder) putGeneral(fieldName, value string) {
	if j.encode.Len() > 0 {
		fmt.Fprintf(&j.encode, `,"%v":%v`, fieldName, value)
	} else {
		fmt.Fprintf(&j.encode, `"%v":%v`, fieldName, value)
	}
}

func (j *jsonBuilder) PutTime(fieldName string, t time.Time) {
	j.putGeneral(fieldName, t.Format(time.RFC3339))
}

func (j *jsonBuilder) PutUint64(fieldName string, value uint64) {
	j.putGeneral(fieldName, fmt.Sprintf("%v", value))
}

func (j *jsonBuilder) PutHex(fieldName string, value []byte) {
	if len(value) == 0 {
		return
	}
	j.putGeneral(fieldName, fmt.Sprintf(`"0x%v"`, hex.EncodeToString(value)))
}

func (j *jsonBuilder) PutBase64(fieldName string, value []byte) {
	if len(value) == 0 {
		return
	}
	j.putGeneral(fieldName, fmt.Sprintf(`"%v"`, base64.StdEncoding.EncodeToString(value)))
}

func (j *jsonBuilder) PutString(fieldName, value string) {
	j.putGeneral(fieldName, fmt.Sprintf(`"%v"`, value))
}

func (j *jsonBuilder) PutJSON(fieldName, value string) {
	j.putGeneral(fieldName, value)
}

func (j *jsonBuilder) ToString() string {
	return fmt.Sprintf(`{%v}`, j.encode.String())
}

func (j *jsonBuilder) PutReciepientArray(fieldName string, receipients []Recipient) {
	if len(receipients) == 0 {
		return
	}
	array := &jsonBuilder{}
	array.encode.WriteRune('[')
	first := true
	for _, r := range receipients {
		if !first {
			array.encode.WriteRune(',')
		}
		fmt.Fprintf(&array.encode, `{"token":"%v","value":%v`, base64.StdEncoding.EncodeToString(r.Token), r.Value)
	}
	array.encode.WriteRune(']')
	j.PutString(fieldName, array.encode.String())
}

func (j *jsonBuilder) PutTokenCiphers(fieldName string, tc TokenCiphers) {
	if len(tc) == 0 {
		return
	}
	array := &jsonBuilder{}
	array.encode.WriteRune('[')
	first := true
	for _, r := range tc {
		if !first {
			array.encode.WriteRune(',')
		}
		fmt.Fprintf(&array.encode, `{"token":"%v","cipher":%v`, base64.StdEncoding.EncodeToString(r.token), base64.StdEncoding.EncodeToString(r.cipher))
	}
	array.encode.WriteRune(']')
	j.PutString(fieldName, array.encode.String())
}

func (a *authoredInstruction) JSON(kind byte, bulk *jsonBuilder) string {
	b := &jsonBuilder{}
	b.PutUint64("version", 0)
	b.PutUint64("instructionType", uint64(kind))
	b.PutUint64("epoch", a.epoch)
	b.PutHex("author", a.author)
	fmt.Fprintf(&b.encode, ",%v", bulk.encode.String())
	b.PutHex("wallet", a.wallet)
	b.PutUint64("fee", a.fee)
	b.PutHex("attorney", a.attorney)
	b.PutBase64("signature", a.signature)
	b.PutBase64("walletSignature", a.walletSignature)
	return b.ToString()
}

func (j *JoinNetwork) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutString("caption", j.caption)
	bulk.PutJSON("details", j.details)
	return j.authored.JSON(iJoinNetwork, bulk)
}

func (j *UpdateInfo) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutJSON("details", j.details)
	return j.authored.JSON(iUpdateInfo, bulk)
}

func (j *GrantPowerOfAttorney) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("details", j.attorney)
	return j.authored.JSON(iGrantPowerOfAttorney, bulk)
}

func (j *RevokePowerOfAttorney) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("details", j.attorney)
	return j.authored.JSON(iRevokePowerOfAttorney, bulk)
}
func (j *CreateEphemeral) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("ephemeralToken", j.ephemeralToken)
	bulk.PutUint64("expiry", j.expiry)
	return j.authored.JSON(iCreateEphemeral, bulk)
}

func (j *SecureChannel) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("tokenRange", j.tokenRange)
	bulk.PutUint64("nonce", j.nonce)
	bulk.PutHex("encryptedNonce", j.encryptedNonce)
	bulk.PutHex("encryptedNonce", j.content)
	return j.authored.JSON(iSecureChannel, bulk)
}

func (j *CreateAudience) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("audience", j.audience)
	bulk.PutHex("submission", j.submission)
	bulk.PutHex("moderation", j.moderation)
	bulk.PutHex("audienceKey", j.audienceKey)
	bulk.PutHex("submissionKey", j.submissionKey)
	bulk.PutHex("moderationKey", j.moderationKey)
	bulk.PutUint64("flag", uint64(j.flag))
	bulk.PutString("description", j.description)
	return j.authored.JSON(iCreateAudience, bulk)
}

func (j *JoinAudience) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("audience", j.audience)
	bulk.PutString("presentation", j.presentation)
	return j.authored.JSON(iJoinAudience, bulk)
}

func (j *AcceptJoinAudience) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("audience", j.audience)
	bulk.PutHex("member", j.member)
	bulk.PutHex("read", j.read)
	bulk.PutHex("submit", j.submit)
	bulk.PutHex("moderate", j.moderate)
	bulk.PutBase64("modSignature", j.modSignature)
	return j.authored.JSON(iAcceptJoinRequest, bulk)
}

func (j *UpdateAudience) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("audience", j.audience)
	bulk.PutHex("submission", j.submission)
	bulk.PutHex("moderation", j.moderation)
	bulk.PutHex("submissionKey", j.submissionKey)
	bulk.PutHex("moderationKey", j.moderationKey)
	bulk.PutUint64("flag", uint64(j.flag))
	bulk.PutString("description", j.description)
	bulk.PutTokenCiphers("readMembers", j.readMembers)
	bulk.PutTokenCiphers("subMembers", j.subMembers)
	bulk.PutTokenCiphers("modMembers", j.modMembers)
	bulk.PutBase64("audSignature", j.audSignature)
	return j.authored.JSON(iUpdateAudience, bulk)
}

func (j *SponsorshipOffer) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("audience", j.audience)
	bulk.PutString("contentType", j.contentType)
	bulk.PutBase64("content", j.content)
	bulk.PutUint64("expiry", j.expiry)
	bulk.PutUint64("revenue", j.revenue)
	return j.authored.JSON(iSponsorshipOffer, bulk)
}

func (j *SponsorshipAcceptance) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("audience", j.audience)
	bulk.PutJSON("offer", j.offer.JSON())
	bulk.PutBase64("modSignature", j.modSignature)
	return j.authored.JSON(iSponsorshipAcceptance, bulk)
}

func (j *React) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutHex("hash", j.hash)
	bulk.PutUint64("reaction", uint64(j.reaction))
	return j.authored.JSON(iReact, bulk)
}

func (j *Content) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutUint64("version", 0)
	bulk.PutUint64("instructionType", uint64(iContent))
	bulk.PutUint64("epoch", j.epoch)
	bulk.PutUint64("published", j.published)
	bulk.PutHex("author", j.author)
	bulk.PutHex("audience", j.audience)
	bulk.PutString("contentType", j.contentType)
	bulk.PutBase64("content", j.content)
	bulk.PutHex("hash", j.hash)
	bulk.PutHex("wallet", j.wallet)
	bulk.PutUint64("fee", j.fee)
	bulk.PutHex("attorney", j.attorney)
	bulk.PutBase64("signature", j.signature)
	bulk.PutBase64("walletSignature", j.walletSignature)
	return bulk.ToString()
}

func (j *Transfer) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutUint64("version", 0)
	bulk.PutUint64("instructionType", uint64(iTransfer))
	bulk.PutUint64("epoch", j.epoch)
	bulk.PutHex("from", j.From)
	bulk.PutReciepientArray("to", j.To)
	bulk.PutString("reason", j.Reason)
	bulk.PutUint64("fee", j.Fee)
	bulk.PutBase64("signature", j.Signature)
	return bulk.ToString()
}

func (j *Deposit) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutUint64("version", 0)
	bulk.PutUint64("instructionType", uint64(iDeposit))
	bulk.PutUint64("epoch", j.epoch)
	bulk.PutHex("token", j.Token)
	bulk.PutUint64("value", j.Value)
	bulk.PutUint64("fee", j.Fee)
	bulk.PutBase64("signature", j.Signature)
	return bulk.ToString()
}

func (j *Withdraw) JSON() string {
	bulk := &jsonBuilder{}
	bulk.PutUint64("version", 0)
	bulk.PutUint64("instructionType", uint64(iDeposit))
	bulk.PutUint64("epoch", j.epoch)
	bulk.PutHex("token", j.Token)
	bulk.PutUint64("value", j.Value)
	bulk.PutUint64("fee", j.Fee)
	bulk.PutBase64("signature", j.Signature)
	return bulk.ToString()
}

func (b *Block) JSONSimple() string {
	bulk := &jsonBuilder{}
	bulk.PutUint64("epoch", b.Epoch)
	bulk.PutHex("parent", b.Parent[:])
	bulk.PutHex("publisher", b.Publisher)
	bulk.PutTime("publishedAt", b.PublishedAt)
	bulk.PutUint64("instructionsCount", uint64(len(b.Instructions)))
	bulk.PutHex("hash", b.Parent[:])
	bulk.PutUint64("feesCollectes", b.FeesCollected)
	bulk.PutBase64("signature", b.Signature)
	return bulk.ToString()
}
