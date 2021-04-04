package goscrobble

type ProfileResponse struct {
	UUID      string                 `json:"uuid"`
	Username  string                 `json:"username"`
	Scrobbles []ScrobbleResponseItem `json:"scrobbles"`
}

func getProfileForUser(user User) (ProfileResponse, error) {
	resp := ProfileResponse{
		UUID:     user.UUID,
		Username: user.Username,
	}
	scrobbleReq, err := getScrobblesForUser(user.UUID, 10, 1)
	if err != nil {
		return resp, err
	}

	resp.Scrobbles = scrobbleReq.Items
	return resp, nil
}
