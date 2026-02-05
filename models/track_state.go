package models

type TrackInteractionState string

const (
	TrackStateLiked    TrackInteractionState = "liked"
	TrackStateDisliked TrackInteractionState = "disliked"
	TrackStateNeutral  TrackInteractionState = "neutral"
)
