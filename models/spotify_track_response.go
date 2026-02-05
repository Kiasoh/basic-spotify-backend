package models

type SpotifyTrackResponse struct {
	SpotifyTrack
	InteractionState TrackInteractionState `json:"interaction_state,omitempty"`
}
