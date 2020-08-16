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
	"unicode"
	//"strconv"
	"./lib"
	"flag"
)

type playerPos struct {
	x, y int16
}

func placeBlock(x int16, y int16, tilePos *[21][21]lib.Block, inv *[10]lib.ItemStack, selSlot int8) {
	ok := true
	if x>19 || x<0 || y>19 || y<0 {
		ok=false
	}
	if tilePos[x][y].Btype>0 {
		ok=false
	}
	if inv[selSlot].Itype==0&&inv[selSlot].Amnt==0 {
		ok=false
	}
	if ok {
		tilePos[x][y].Btype = inv[selSlot].Itype
		inv[selSlot].Amnt--
	}
}

func moveCheck(x int16, y int16, tilePos [21][21]lib.Block) (ok bool) {
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
	if tilePos[x][y].Btype>0 {
		return false
	}
	return true
}



func deleteBrokenBlocks(tilePos *[21][21]lib.Block, inv [10]lib.ItemStack, breakThresh int32) ([10]lib.ItemStack) {
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			if tilePos[x][y].BreakStage >= breakThresh {	
				lib.GiveItem(&inv, lib.IS(1, tilePos[x][y].Btype))
				tilePos[x][y].Btype = 0
				tilePos[x][y].BreakStage = 0
			}
		}
	}
	return inv
}

func GetExecPath() string {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    return filepath.Dir(ex)
}

func LoadPicture(path string) (pixel.Picture, error) {
	curPath := GetExecPath()
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

	config := lib.ReadJson("config.json")
	bt:=int(config.BreakThresh)
	wt:=int(config.WalkThresh)
	ws:=int(config.WalkSpeed)
	flag.IntVar(&bt, "bt", bt, "How many frames it takes to break blocks")
	flag.IntVar(&wt, "wt", wt, "How long you have to hold WASD until you move faster in frames")
	flag.IntVar(&ws, "ws", ws, "How many frames between walking (after WalkThresh been surpased)")
	flag.Parse()
	config.BreakThresh = int32(bt)
	config.WalkThresh = int16(wt)
	config.WalkSpeed = int16(ws)

	cfg := pixelgl.WindowConfig{ //the settings for the window
		Title:  "Texxit",
		Bounds: pixel.R(0, 0, 640, 704),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	//load the images from the disc
	asset, err := LoadPicture("asset.png")	//loading all the assets
	if err != nil {
		panic(err)
	}
	//make the graphics batches and sprites and stuff
	var pixs [8][8]*pixel.Sprite
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			pixs[x][y] = pixel.NewSprite(asset, pixel.R(float64(x*16), float64(y*16), float64(x*16)+16, float64(y*16)+16))
		}
	}

	numSpr[0] := pixs[0][1]
	numSpr[1] := pixs[1][1]
	numSpr[2] := pixs[2][1]
	numSpr[1] := pixs[3][1]
	numSpr[1] := pixs[4][1]
	numSpr[1] := pixs[5][1]
	numSpr[1] := pixs[6][1]
	numSpr[1] := pixs[7][1]
	numSpr[1] := pixs[0][2]
	numSpr[1] := pixs[1][2]
	youSpr := pixs[1][0]
	invSpr := pixs[2][0]
	//guiSpr := pixel.NewSprite(asset, pixel.R(32, 32, 56, 48))
	//arrowSpr := pixel.NewSprite(asset, pixel.R(56, 32, 64, 48))

	var opics [100]*pixel.Sprite
	opics[0] = pixs[3][0]
	opics[1] = pixs[0][0]
	opics[2] = pixs[5][0]
	opics[3] = pixs[4][0]

	draw := pixel.NewBatch(&pixel.TrianglesData{}, asset)
	var (
		player = playerPos{9, 9}
		frames int16 = 0
		second = time.Tick(time.Second)
		selSlot int8
		tilePos [21][21]lib.Block
		//randAngles = [4]float64{0, 1.5708, 3.14159, 4.71239}
		inv [10]lib.ItemStack
		holdWASD [4]int16
		//grassRotSeed = time.Now().UnixNano()
		gui int16 = 0
	)
	lib.Gen(&tilePos)
	if tilePos[player.x][player.y].Btype>0{
		player = playerPos{10, 10}
	}
	//tilePos[5][5] = Block{2, 49}
	rand.Seed(time.Now().UnixNano())

	draw.Clear()
	for !win.Closed() {
		//BEGIN CONTROLS
		//placing/mining bloks
		if win.Pressed(pixelgl.KeyUp) {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y+1].Btype==0 {
				placeBlock(player.x, player.y+1, &tilePos, &inv, selSlot)
			} else {
				if !win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y+1].Btype!=0 {
					tilePos[player.x][player.y+1].BreakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyDown)&&player.y!=0 {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y-1].Btype==0 {
				placeBlock(player.x, player.y-1, &tilePos, &inv, selSlot)
			} else {
				if !win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y-1].Btype!=0 {
					tilePos[player.x][player.y-1].BreakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x+1][player.y].Btype==0 {
				placeBlock(player.x+1, player.y, &tilePos, &inv, selSlot)
			} else {
				if !win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x+1][player.y].Btype!=0 {
					tilePos[player.x+1][player.y].BreakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyLeft)&&player.x!=0 {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x-1][player.y].Btype==0 {
				placeBlock(player.x-1, player.y, &tilePos, &inv, selSlot)
			} else {
				if !win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x-1][player.y].Btype!=0 {
					tilePos[player.x-1][player.y].BreakStage++
				}
			}
		}

		inv = deleteBrokenBlocks(&tilePos, inv, config.BreakThresh)
		//selecting slots
		x := win.Typed()
		if len(x)!=0 {
			ok:=false
			var xnum int
			if unicode.IsDigit([]rune(x)[0]) {
				ok = true
			}
			for i := 0; i < 5; i++ {
				if len(x)>i+1 {
					if unicode.IsDigit([]rune(x)[i+1]) {
						ok = true
						xnum = i+1
					}
				}
			}
			if ok {
				i := int(x[xnum])-48
				if err!=nil {
					panic(err)
				}
				if i==0 {
					i=10
				}
				i--
				if !win.Pressed(pixelgl.KeySpace) { 
					selSlot = int8(i)
				} else {
					x:=inv[i]
					y:=inv[selSlot]
					inv[selSlot] = x
					inv[i] = y
				}
			}
		}

		if win.Pressed(pixelgl.KeySpace)&&win.Pressed(pixelgl.KeyLeftShift) {
			var x rune
			if len(win.Typed())!=0 {
				x=rune(win.Typed()[0])
			}
			var i int
			switch x {
			case 32:
				i=-1
			case 33:
				i=0
			case 64:
				i=1
			case 35:
				i=2
			case 36:
				i=3
			case 37:
				i=4
			case 94:
				i=5
			case 38:
				i=6
			case 42:
				i=7
			case 40:
				i=8
			case 41:
				i=9
			default:
				i=-1
			}
			if i!=-1 {
				if inv[i].Itype==inv[selSlot].Itype||inv[i].Itype==0 {
					x:=inv[selSlot].Amnt/2
					if inv[i].Itype==0 {
						inv[i].Itype = inv[selSlot].Itype
					}
					for x>0 {
						x--
						inv[selSlot].Amnt--
						inv[i].Amnt++
						if inv[i].Amnt>85 {
							break
						}
					}
					if inv[i].Amnt==0 {
						inv[i].Itype = 0
					}
				}
			}
		}


		if win.JustPressed(pixelgl.KeyK)&&win.Pressed(pixelgl.KeyY)&&win.Pressed(pixelgl.KeyP)&&win.Pressed(pixelgl.KeyG) {
			lib.GiveItem(&inv, lib.IS(83, 1))
		}


		//managine inventory items
		if win.Pressed(pixelgl.KeyQ)&&win.Pressed(pixelgl.KeySpace) {
			x:=inv[selSlot]
			inv[selSlot]=lib.IS(0,0)
			lib.GiveItem(&inv, x)
		}

		//moving
		//up
		if win.JustPressed(pixelgl.KeyW) {
			if moveCheck(player.x, player.y+1, tilePos) {
				player.y++
			}
		}
		if win.Pressed(pixelgl.KeyW) {
			if holdWASD[0]!=config.WalkThresh {
				holdWASD[0]++
			} else {
				holdWASD[0]-=config.WalkSpeed
				if moveCheck(player.x, player.y+1, tilePos) {
					player.y++
				}
			}
		} else {
			holdWASD[0] = 0
		}
		//left
		if win.JustPressed(pixelgl.KeyA) {
			if moveCheck(player.x-1, player.y, tilePos) {
				player.x--
			}
		}
		if win.Pressed(pixelgl.KeyA) {
			if holdWASD[1]!=config.WalkThresh {
				holdWASD[1]++
			} else {
				holdWASD[1]-=config.WalkSpeed
				if moveCheck(player.x-1, player.y, tilePos) {
				player.x--
			}
			}
		} else {
			holdWASD[1] = 0
		}
		//down
		if win.JustPressed(pixelgl.KeyS) {
			if moveCheck(player.x, player.y-1, tilePos) {
				player.y--
			}
		}
		if win.Pressed(pixelgl.KeyS) {
			if holdWASD[2]!=config.WalkThresh {
				holdWASD[2]++
			} else {
				holdWASD[2]-=config.WalkSpeed
				if moveCheck(player.x, player.y-1, tilePos) {
					player.y--
				}
			}
		} else {
			holdWASD[2] = 0
		}
		//right
		if win.JustPressed(pixelgl.KeyD) {
			if moveCheck(player.x+1, player.y, tilePos) {
				player.x++
			}
		}
		if win.Pressed(pixelgl.KeyD) {
			if holdWASD[3]!=config.WalkThresh {
				holdWASD[3]++
			} else {
				holdWASD[3]-=config.WalkSpeed
				if moveCheck(player.x+1, player.y, tilePos) {
					player.x++
				}
			}
		} else {
			holdWASD[3] = 0
		}

		//crafting
		if win.JustPressed(pixelgl.KeyC) {
			if gui==0 {
				gui = 1
			} else {
				gui = 0
			}
		}

		//END CONTROLS
		//BEGIN RENDERING
		win.Clear(colornames.Forestgreen)
		draw.Clear()

		for x := 0; x != 20; x++ {
			for y := 0; y != 20; y++ {
				move:=pixel.IM.Moved(pixel.V(float64(x*32)+8, float64(y*32)+40)).ScaledXY(pixel.V(float64(x*32), float64(y*32)), pixel.V(2, 2))
				opics[tilePos[x][y].Btype].Draw(draw, move)
			}
		}
		youSpr.Draw(draw, pixel.IM.Moved(pixel.V(float64(player.x*16)+8, float64(player.y*16)+40)).ScaledXY(pixel.V(1, 1), pixel.V(2, 2)))
		for x := 32; x != 672; x += 64 {
			curSlot := (x+32)/64-1
			if int8(curSlot)==selSlot {
				invSpr.Draw(draw, pixel.IM.Moved(pixel.V(float64(x), float64(32))).ScaledXY(pixel.V(float64(x), float64(32)), pixel.V(4.6, 4.6)))
			} else {
				invSpr.Draw(draw, pixel.IM.Moved(pixel.V(float64(x), float64(32))).ScaledXY(pixel.V(float64(x), float64(32)), pixel.V(4, 4)))
			}
		}
		

		for i := 0; i < 10; i++ {
			move:=pixel.IM.Moved(pixel.V(float64((i*64)+32), 32)).ScaledXY(pixel.V(float64((i*64)+32), 32), pixel.V(2.25, 2.25))
			if inv[i].Itype!=0 {
				opics[inv[i].Itype].Draw(draw, move)
			}
			if inv[i].Amnt<10 {
				switch inv[i].Amnt {
				case 0:
					inv[i].Itype = 0
				case 1:
				case 2:
					num2.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+50), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 3:
					num3.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+50), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 4:
					num4.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+50), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 5:
					num5.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+50), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 6:
					num6.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+50), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 7:
					num7.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+50), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 8:
					num8.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+50), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 9:
					num9.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+50), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				}
			} else { // number i bigger then 10 
				switch inv[i].Amnt / 10 { //most signigf=iashudant digits
				case 0:
					num0.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 1:
					num1.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 2:
					num2.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 3:
					num3.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 4:
					num4.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 5:
					num5.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 6:
					num6.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 7:
					num7.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 8:
					num8.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 9:
					num9.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+44), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				}

				switch inv[i].Amnt % 10  { // least signifiicant digits
				case 0:
					num0.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 1:
					num1.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 2:
					num2.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 3:
					num3.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 4:
					num4.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 5:
					num5.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 6:
					num6.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 7:
					num7.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 8:
					num8.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				case 9:
					num9.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+56), 20)).ScaledXY(pixel.V(float64((i*64)+50), 20), pixel.V(1.25, 1.25)))
				}
			}
		}
		
		switch gui {
		case 0:
		case 1:
			lib.CraftGUI(&gui, &inv, asset, pixs, opics)
		default:
		}

		draw.Draw(win)
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