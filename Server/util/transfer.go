package util

import (
	"Server/Game"
	"Server/Net"
)

func Game2NetRoom(room *Game.Room)*Net.Message{
	var audience []string
	if room.Audience != nil{
		for k,_ := range room.Audience{
			audience = append(audience,k)
		}
	}

	return Net.PackMsg(&Net.RoomDataMsg{
		RoomID:    room.Id,
		FirstPlayer:   Game2NetPlayer(room.FirstPlayer),
		LastPlayer:   Game2NetPlayer(room.LastPlayer),
		Audience:  audience,
		Size:      room.Size,
	//	Snapshot:  room.Snapshot,
		End:       room.End,
		Winner:    room.Winner,
		RoomState: room.RoomState,
	})
}

func Game2NetPlayer(player *Game.Player)*Net.Player{
	if player==nil{
		return nil
	}
	return &Net.Player{
		ID:    player.ID,
		Score: player.Score,
	}
}

func Net2GamePlayer(player *Net.Player)*Game.Player{
	if player==nil{
		return nil
	}
	return &Game.Player{
		ID:    player.ID,
		Score: player.Score,
	}
}