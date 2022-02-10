package tribunal

import "geoip/pkg/model"

const (
	// ReasonEqualHash reason type
	ReasonEqualHash = "equal_hash"

	// ReasonNoPreviousLogin no previous logins found in db
	ReasonNoPreviousLogin = "no_previous_login"

	// ReasonMatchDeviceIDCountry current deviceID and Country matches previous login
	ReasonMatchDeviceIDCountry = "match_deviceid_country"

	// ReasonMatchDeviceIDNotCountry current deviceID match but not country
	ReasonMatchDeviceIDNotCountry = "match_deviceid_not_country"

	// SentenceNone rendering no sentence
	SentenceNone = "none"

	// SentenceNotice just notice the user of the activity
	SentenceNotice = "notice"

	// SentenceBlock block the user from futher access until challenge
	SentenceBlock = "block"

	// SentenceChallenge challenge the user to prove it's legitimate
	SentenceChallenge = "challenge"
)

// Verdict holds the verdict
type Verdict struct {
	Reason       string
	Sentence     string
	LoginEventID string
	FishScore    int
}

type Client struct {
	Previous []*model.LoginEvent
	Current  *model.LoginEvent
}

// Decision conclude a verdict
func (c *Client) Decision() *Verdict {
	return nil
}

func (c *Client) EstimatePhishingScore() *Verdict {
	// No previous logins found = 0
	if c.Previous == nil {
		return &Verdict{
			Reason:       ReasonNoPreviousLogin,
			Sentence:     SentenceNone,
			LoginEventID: "",
			FishScore:    0,
		}
	}

	// hash is equal = 0
	for _, previous := range c.Previous {
		if previous.Hash == c.Current.Hash {
			return &Verdict{
				Reason:       ReasonEqualHash,
				Sentence:     SentenceNone,
				LoginEventID: previous.ID,
				FishScore:    0,
			}
		}
	}

	// Check deviceId in different combination
	for _, previous := range c.Previous {
		if previous.DeviceID == c.Current.DeviceID {
			if previous.IP.Country == c.Current.IP.Country {
				return &Verdict{
					Reason:       ReasonMatchDeviceIDCountry,
					Sentence:     SentenceNone,
					LoginEventID: previous.ID,
					FishScore:    0,
				}
			}
			// equal deviceId, but unequal country
			return &Verdict{
				Reason:       ReasonMatchDeviceIDNotCountry,
				Sentence:     SentenceNone,
				LoginEventID: previous.ID,
				FishScore:    50,
			}
		} else {

		}
	}

	return &Verdict{
		Reason: ReasonEqualHash,
	}
}
