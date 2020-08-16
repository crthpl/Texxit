package lib

import (
	//"fmt"
	"math/rand"
	"time"
)

type Block struct {
	Btype uint16
	BreakStage int32
}

func Gen(tilePos*[21][21]Block) (){
	rand.Seed(time.Now().UnixNano())
	for x := rand.Intn(3); x < 21; x+=4 {
		for y := rand.Intn(3); y < 21; y+=4 {
			xoff := rand.Intn(3)-1
			yoff := rand.Intn(3)-1
			if !(x+xoff>20)&&!(y+yoff>20)&&!(x+xoff<0)&&!(y+yoff<0) {
				tilePos[x+xoff][y+yoff] = Block{2, 0}
			}
		}
	}
}
