package response

import (
	"time"

	"rederinghub.io/utils/constants/dao_artist"
	"rederinghub.io/utils/constants/dao_artist_voted"
)

type DaoArtist struct {
	BaseEntity     `json:",inline"`
	SeqId          uint              `json:"seq_id"`
	CreatedBy      string            `json:"created_by"`
	User           *UserForDao       `json:"user"`
	ExpiredAt      time.Time         `json:"expired_at"`
	Status         dao_artist.Status `json:"status"`
	DaoArtistVoted []*DaoArtistVoted `json:"dao_artist_voted"`
	Action         *ActionDaoArtist  `json:"action"`
	TotalReport    int64             `json:"total_report"`
	TotalVerify    int64             `json:"total_verify"`
}

func (s DaoArtist) Expired() bool {
	return s.ExpiredAt.UTC().Unix() < time.Now().UTC().Unix()
}

func (s *DaoArtist) SetFields(fns ...func(*DaoArtist)) {
	for _, fn := range fns {
		fn(s)
	}
}
func (s DaoArtist) WithAction(action *ActionDaoArtist) func(*DaoArtist) {
	return func(dp *DaoArtist) {
		dp.Action = action
	}
}

type ActionDaoArtist struct {
	CanVote bool `json:"can_vote"`
}

type DaoArtistVoted struct {
	CreatedBy   string                  `json:"created_by"`
	DisplayName string                  `json:"display_name"`
	DaoArtistId string                  `json:"dao_project_id"`
	Status      dao_artist_voted.Status `json:"status"`
}

func (s *DaoArtistVoted) SetFields(fns ...func(*DaoArtistVoted)) {
	for _, fn := range fns {
		fn(s)
	}
}
func (s DaoArtistVoted) WithDisplayName(displayName string) func(*DaoArtistVoted) {
	return func(dp *DaoArtistVoted) {
		dp.DisplayName = displayName
	}
}
