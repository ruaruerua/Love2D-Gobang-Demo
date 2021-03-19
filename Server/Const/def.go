package Const


const (
	Hello = iota
	GAME
	Exit
	Err
	Match
	
	Invite
	Room
	RoomData

	Msg
	Heart
)

const (
	InviteCreate = iota
	InviteBack

	RoomActionIn
	RoomActionExit
	RoomActionStart
	RoomActionChangeSite
	RoomActionAgain
	RoomActionReady
)

const (
	OnlineState = iota
	HallState
	RoomState
	GameState
)

const (
	ErrPlayerNotExist = "player not found in db"
	ErrPlayerOffline  = "player should login"
	ErrPlayIsOnline   = "player has been logged in"
	ErrPlayIsPlaying  = "player is gaming"
	ErrProcess        = "plz transmit the correct message"
	ErrMismatchCMD    = "cmd match err"
	ErrAction         = "body action err"
	ErrRoomNotExist   = "room not exist"
	ErrGame           = "game err"
	ErrHasChess       = "path has chess"
	ErrOutRange       = "chess out range"
	ErrPlayOp = "play result get err result"
	ErrSiteErr = "change site err"
	ErrAgain = "again default"
	ErrStart = "start default plz check the room state"
	ErrReadyGame = "player con't check to ready state"
	ErrNotSupportCmd = "server can not support this cmd"

	MsgMatchBack = "matching"
	MsgExitBack = "exit room"
)

const (
	CmdMsgInfo = "info"
	CmdMsgAllRoomBase = "room"
	CmdMsgRoomData = "roomData"
	CmdMsgPlayerState = "playerState"
)

const (
	PlayerNotEnough = iota
	PlayerNotInSite
	PlayerReady
)

const (
	NoState = iota
	ReadyState
	PlayingState

	NotExist
	AlreadyInRoom
	ExitOldRoom
	OnlyInRoom


)

const (
	PlayChess = iota
	PlayUndo
	PlayUndoBack

	PlayDraw
	PlayDrawBack

	PlayRegret
)

const (
	SiteFirst = iota
	SiteLast
	SiteViewer
	SiteNil
)

const (
	BoardSizeDefault = 15

	PlayOK = iota
	PlayHasChess
	PlayOutRange
	PlayTurnErr
	PlayEnd
)