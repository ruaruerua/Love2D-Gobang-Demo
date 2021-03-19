package Game

import (
	"Server/Const"
)

type Board struct {
	Size     int
	FirstPlayer *Player
	LastPlayer   *Player
	End      bool
	Winner   string
	Snapshot []int
	StepNum int
	historyList  [][]int
}

type Player struct {
	ID string
	Score int
	State int
}

type History struct {
	AllSteps []Chess
	Result string
	FirstID string
	LastID string
}

type Chess struct {
	X int
	Y int
}

func (b *Board)StartGame(){
	b.Snapshot = make([]int,b.Size*b.Size,b.Size*b.Size)
	b.StepNum = 1
	b.End = false
}

func NewBoard()*Board{
	return &Board{
		Size:        Const.BoardSizeDefault,
		FirstPlayer:nil,
		LastPlayer:  nil,
		End:         true,
		Winner:      "",
		Snapshot:    nil,
		StepNum:     1,
		historyList: nil,
	}
}

func (b *Board)ExchangeSite()bool{
	if b.End{
		b.FirstPlayer,b.LastPlayer = b.LastPlayer,b.FirstPlayer
		return true
	}
	return false
}

//func (b *Board)

func (b *Board)Play(player string,chess Chess)int{
	curTurn := b.StepNum % 2 == 1
	if curTurn && b.FirstPlayer.ID != player{
		return Const.PlayTurnErr
	}
	if !curTurn && b.LastPlayer.ID != player{
		return Const.PlayTurnErr
	}
	if b.checkRange(chess){
		return Const.PlayOutRange
	}
	if b.checkHasChess(chess){
		return Const.PlayHasChess
	}

	if curTurn{
		b.Snapshot[chess.Y * b.Size + chess.X] = b.StepNum
	} else {
		b.Snapshot[chess.Y * b.Size + chess.X] = b.StepNum
	}
	b.StepNum++
	if b.checkWin(chess){
		b.EndGame()
		return Const.PlayEnd
	}

	return Const.PlayOK
}

func (b *Board)checkWin(chess Chess)bool{
	countA,conutB := 0,0
	countA = b.countChess(chess,1,0)
	conutB = b.countChess(chess,-1,0)
	if conutB + countA + 1 >= 5{
		return true
	}

	countA,conutB = 0,0
	countA = b.countChess(chess,0,1)
	conutB = b.countChess(chess,0,-1)
	if conutB + countA + 1 >= 5{
		return true
	}

	countA,conutB = 0,0
	countA = b.countChess(chess,1,1)
	conutB = b.countChess(chess,-1,-1)
	if conutB + countA + 1 >= 5{
		return true
	}

	countA,conutB = 0,0
	countA = b.countChess(chess,1,-1)
	conutB = b.countChess(chess,-1,1)
	if conutB + countA + 1 >= 5{
		return true
	}

	return false
}

func (b *Board)countChess(chess Chess,addX int,addY int)int{
	ret := 0
	srcChess := chess
	for{
		chess.X += addX
		chess.Y += addY
		if b.checkRange(chess) || !b.isSameChess(chess,srcChess){
			return ret
		}
		ret++
	}
}

func (b *Board)isSameChess(chessA Chess,chessB Chess)bool {
	return b.Snapshot[chessA.Y * b.Size + chessA.X] > 0 &&
		b.Snapshot[chessB.Y * b.Size + chessB.X]>0 &&
		b.Snapshot[chessA.Y * b.Size + chessA.X] % 2 == b.Snapshot[chessB.Y * b.Size + chessB.X] % 2
}

func (b *Board)checkHasChess(chess Chess)bool{
	return b.Snapshot[chess.Y * b.Size + chess.X] != 0
}

func (b *Board)checkRange(chess Chess)bool{
	return chess.X >= b.Size || chess.X<0 ||chess.Y >= b.Size || chess.Y<0
}

func (b *Board)Regret(){
}

func (b *Board)EndGame(){
	b.End = true
	if b.historyList == nil{
		b.historyList = make([][]int,0)
	}
	b.historyList = append(b.historyList, b.Snapshot)
	var win,lose *Player
	if b.StepNum % 2 == 1{
		win = b.FirstPlayer
		lose = b.LastPlayer
	} else {
		win = b.LastPlayer
		lose = b.FirstPlayer
	}
	b.Winner = win.ID
	win.Score++
	lose.Score--
}

func (b *Board)Draw(){
	b.End = true
	b.Winner = "draw"
}


func (b *Board) Defect(playerID string)bool{
	//if b.End == true || b.StepNum == 1{
	//	return false
	//}

	if b.FirstPlayer != nil && playerID != b.FirstPlayer.ID &&
		b.LastPlayer != nil &&playerID != b.LastPlayer.ID{
		return false
	}

	b.End = true
	if playerID == b.FirstPlayer.ID{
		b.LastPlayer.Score++
		b.FirstPlayer.Score--
	} else {
		b.FirstPlayer.Score++
		b.LastPlayer.Score--
	}
	b.Winner = playerID
	return true
}



