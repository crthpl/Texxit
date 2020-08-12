package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"time"
	"math/rand"
	"path/filepath"
	"os"
	"image"
	_ "image/png"
	"golang.org/x/image/colornames"
)

type ItemStack struct {
	amnt int8
	itype int32
}

type playerPos struct {
	x, y int16
}
type Block struct {
	btype int32
	breakStage int16
}

func placeBlock(x int16, y int16, tilePos *[20][20]Block, blockType int32) {
	ok := true
	if x>19 {
		ok=false
	}
	if x<0 {
		ok=false
	}
	if y>19 {
		ok=false
	}
	if y<0 {
		ok=false
	}
	if ok {
		tilePos[x][y].btype = blockType
	}
}

func moveCheck(x int16, y int16, tilePos [20][20]Block) (ok bool) {
	if x>19 {
		return false
	}
	if x<0 {
		return false
	}
	if y>19 {
		return false
	}
	if y<0 {
		return false
	}
	if tilePos[x][y].btype>0 {
		return false
	}
	return true
}

func deleteBrokenBlocks(tilePos *[20][20]Block) {
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			if tilePos[x][y].breakStage == 100 {
				tilePos[x][y].btype = 0
			}
		}
	}
}

func LoadPicture(path string) (pixel.Picture, error) {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
	curPath := filepath.Dir(ex)
	path = curPath+"/assets/"+path
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{ //the settings for the window
		Title:  "Online Game",
		Bounds: pixel.R(0, 0, 640, 704),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	//load the images from the disc
	grass, err := LoadPicture("grass.png")	//loading the grass tile
	tile, err := LoadPicture("wood.png")	//loading the tile tile
	you, err := LoadPicture("you.png")		//loading you
	invp, err := LoadPicture("inv.png")	//loading the inv hotbar slors
	if err != nil {
		panic(err)
	}
	//make the graphics batches and sprites
	grasses := pixel.NewBatch(&pixel.TrianglesData{}, grass)
	grassSpr := pixel.NewSprite(grass, grass.Bounds())
	tileSpr := pixel.NewSprite(tile, tile.Bounds())
	youSpr := pixel.NewSprite(you, you.Bounds())
	invs := pixel.NewBatch(&pixel.TrianglesData{}, invp)
	invSpr := pixel.NewSprite(invp, invp.Bounds())

	var (
		player = playerPos{10, 10}
		frames int = 0
		second = time.Tick(time.Second)
		selSlot int8
		tilePos [20][20]Block
		randAngles = [4]float64{0, 1.5708, 3.14159, 4.71239}
		//inv [52]ItemStack
	)
	rand.Seed(time.Now().UnixNano())

	grasses.Clear()

	for x := 16; x != 656; x += 32 {
		for y := 80; y != 720; y += 32 {
			grassSpr.Draw(grasses, pixel.IM.Moved(pixel.V(float64(x), float64(y))).ScaledXY(pixel.V(float64(x), float64(y)), pixel.V(2, 2)).Rotated(pixel.V(float64(x), float64(y)), randAngles[rand.Intn(3)]))
		}
	}

	for !win.Closed() {
		//BEGIN CONTROLS
		//placing/mining bloks
		if win.Pressed(pixelgl.KeyUp) {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y+1].btype==0 {
				placeBlock(player.x, player.y+1, &tilePos, 1)
			} else {
				if !(tilePos[player.x][player.y+1].btype==0) {
					tilePos[player.x][player.y+1].breakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyDown) {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y-1].btype==0 {
				placeBlock(player.x, player.y-1, &tilePos, 1)
			} else {
				if !(tilePos[player.x][player.y-1].btype==0) {
					tilePos[player.x][player.y-1].breakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x+1][player.y].btype==0 {
				placeBlock(player.x+1, player.y, &tilePos, 1)
			} else {
				if !(tilePos[player.x+1][player.y].btype==0) {
					tilePos[player.x+1][player.y].breakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyLeft) {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x-1][player.y].btype==0 {
				placeBlock(player.x-1, player.y, &tilePos, 1)
			} else {
				if !(tilePos[player.x-1][player.y].btype==0) {
					tilePos[player.x-1][player.y].breakStage++
				}
			}
		}

		deleteBrokenBlocks(&tilePos)
		//selecting slots
		if win.JustPressed(pixelgl.Key1) {
			selSlot = 0
		}
		if win.JustPressed(pixelgl.Key2) {
			selSlot = 1
		}
		if win.JustPressed(pixelgl.Key3) {
			selSlot = 2
		}
		if win.JustPressed(pixelgl.Key4) {
			selSlot = 3
		}
		if win.JustPressed(pixelgl.Key5) {
			selSlot = 4
		}
		if win.JustPressed(pixelgl.Key6) {
			selSlot = 5
		}
		if win.JustPressed(pixelgl.Key7) {
			selSlot = 6
		}
		if win.JustPressed(pixelgl.Key8) {
			selSlot = 7
		}
		if win.JustPressed(pixelgl.Key9) {
			selSlot = 8
		}
		if win.JustPressed(pixelgl.Key0) {
			selSlot = 9
		}
		//moving
		if win.JustPressed(pixelgl.KeyW) {
			if moveCheck(player.x, player.y+1, tilePos) {
				player.y++
			}
		}
		if win.JustPressed(pixelgl.KeyA) {
			if moveCheck(player.x-1, player.y, tilePos) {
				player.x--
			}
		}
		if win.JustPressed(pixelgl.KeyS) {
			if moveCheck(player.x, player.y-1, tilePos) {
				player.y--
			}
		}
		if win.JustPressed(pixelgl.KeyD) {
			if moveCheck(player.x+1, player.y, tilePos) {
				player.x++
			}
		}
		//END CONTROLS
		//BEGIN RENDERING
		win.Clear(colornames.Forestgreen)
		grasses.Draw(win)
		for x := 0; x != 20; x++ {
			for y := 0; y != 20; y++ {
				if tilePos[x][y].btype>0 {
					tileSpr.Draw(win, pixel.IM.Moved(pixel.V(float64(x*32)+8, float64(y*32)+40)).ScaledXY(pixel.V(float64(x*32), float64(y*32)), pixel.V(2, 2)))
				}
			}
		}
		youSpr.Draw(win, pixel.IM.Moved(pixel.V(float64(player.x*16)+8, float64(player.y*16)+40)).ScaledXY(pixel.V(1, 1), pixel.V(2, 2)))
		invs.Clear()
		for x := 32; x != 672; x += 64 {
			curSlot := (x+32)/64-1
			if int8(curSlot)==selSlot {
				invSpr.Draw(invs, pixel.IM.Moved(pixel.V(float64(x), float64(32))).ScaledXY(pixel.V(float64(x), float64(32)), pixel.V(4.6, 4.6)))
			} else {
				invSpr.Draw(invs, pixel.IM.Moved(pixel.V(float64(x), float64(32))).ScaledXY(pixel.V(float64(x), float64(32)), pixel.V(4, 4)))
			}
			/*switch inv[curSlot].itype {
			case 0:
			case 1:
				tileSpr.Draw(win, pixel.IM.Moved(pixel.V(float64(x), 32)).ScaledXY(pixel.V(float64(x), 32), pixel.V(2.3, 2.3)))
				//if inv[curSlot].amnt<1 {
				//	fmt.Fprintln(basicTxt, inv[curSlot].amnt)
				//	basicTxt.WriteRune(rune(inv[curSlot].amnt))
				//}
			default:
			}*/
		}
		invs.Draw(win)
		win.Update()
		//END RENDERING
		frames++ //Fps displaying stuff
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0 //End Fps displaying stuff
		default:
		}
	}
}

func main() {
	pixelgl.Run(run) // all the graphics stuff (end everything else...)
}