package instructions

import (
	"reflect"
	"testing"
)

func TestSponsorshipOffer(t *testing.T) {
	// mapTest := map
	sponsorshipOffer := author.NewSponsorshipOffer(audienceTest, "test", mapTest, 12, 100, 10, 2000)
	sponsorshipOffer2 := ParseSponsorshipOffer(sponsorshipOffer.Serialize())
	if sponsorshipOffer2 == nil {
		t.Error("could not parse SponsorshipOffer")
		return
	}
	if !reflect.DeepEqual(sponsorshipOffer, sponsorshipOffer2) {
		t.Error("Parse and Serialize not working for SponsorshipOffer")
	}
}

func TestSponsorshipAcceptance(t *testing.T) {
	sponsorshipOffer := author.NewSponsorshipOffer(audienceTest, "test", mapTest, 12, 100, 10, 2000)
	sponsorshipAcceptance := author.NewSponsorshipAcceptance(audienceTest, sponsorshipOffer, 10, 2000)
	sponsorshipAcceptance2 := ParseSponsorshipAcceptance(sponsorshipAcceptance.Serialize())
	if sponsorshipAcceptance2 == nil {
		t.Error("could not parse SponsorshipAcceptance")
		return
	}
	if !reflect.DeepEqual(sponsorshipAcceptance, sponsorshipAcceptance2) {
		t.Error("Parse and Serialize not working for SponsorshipAcceptance")
	}

}
