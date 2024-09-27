package googleapi

import (
	meet "google.golang.org/api/meet/v2"
)

func GetGoogleMeetRoomInfo(spaceName string, meetService *meet.Service) (*meet.Space, error) {
	return meetService.Spaces.Get("spaces/" + spaceName).Do()
}
