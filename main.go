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
	//"strings"
	//"bufio"
	"io/ioutil"
	"encoding/json"
	"unicode"
	"strconv"
)

type config struct {
	breakThresh int32 `json:"breakThresh"`
	defaultVsync bool `json:"defaultVsync"`
}

type ItemStack struct {
	amnt int8
	itype int32
}

type playerPos struct {
	x, y int16
}
type Block struct {
	btype int32
	breakStage int32
}

func placeBlock(x int16, y int16, tilePos *[21][21]Block, inv *[10]ItemStack, selSlot int8) {
	ok := true
	if x>19 || x<0 || y>19 || y<0 {
		ok=false
	}
	if tilePos[x][y].btype>0 {
		ok=false
	}
	if inv[selSlot].itype==0&&inv[selSlot].amnt==0 {
		ok=false
	}
	if ok {
		tilePos[x][y].btype = inv[selSlot].itype
		inv[selSlot].amnt--
	}
}

func moveCheck(x int16, y int16, tilePos [21][21]Block) (ok bool) {
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

func GiveItem(inv *[10]ItemStack, ite ItemStack) {
	iltd:=ite.amnt
	for i := 0; i < 10; i++ { 
		if inv[i].itype==ite.itype {
			inv[i].itype=ite.itype
			for iltd!=0&&inv[i].amnt<85 {
				inv[i].amnt+=1
				iltd-=1
			}
		}
	}
	for i := 0; i < 10; i++ {
		if inv[i].itype==ite.itype||inv[i].itype==0 {
			if inv[i].amnt<85&&iltd>0 {
				inv[i].itype=ite.itype
				for iltd!=0&&inv[i].amnt<85 {
					inv[i].amnt+=1
					iltd-=1
				}
			}
		}
	}
}

func deleteBrokenBlocks(tilePos *[21][21]Block, inv [10]ItemStack, breakThresh int32) ([10]ItemStack) {
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			if tilePos[x][y].breakStage >= breakThresh {	
				GiveItem(&inv, ItemStack{1, tilePos[x][y].btype})
				tilePos[x][y].btype = 0
				tilePos[x][y].breakStage = 0
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

func loadConfig(filename string) (config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return config{}, err
	}

	var c config
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return config{}, err
	}

	return c, nil
}

func run() {

	/*loadedConfig, err := loadConfig(GetExecPath()+"/config.json")
	if err != nil {
		panic(err)
	}

	/*reader := bufio.NewReader(os.Stdin)
	fmt.Print("VSync on?(y/n)")
	ans, _ := reader.ReadString('\n')
	ans = strings.Replace(ans, "\n", "", -1)
	sync := loadedConfig.defaultVsync
	if ans=="n" {
		sync = false
	}
	if ans=="y" {
		sync = true
	}*/

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
	num0 := pixel.NewSprite(asset, pixel.R(0, 16, 16, 32))
	num1 := pixel.NewSprite(asset, pixel.R(16, 16, 32, 32))
	num2 := pixel.NewSprite(asset, pixel.R(32, 16, 48, 32))
	num3 := pixel.NewSprite(asset, pixel.R(48, 16, 64, 32))
	num4 := pixel.NewSprite(asset, pixel.R(64, 16, 80, 32))
	num5 := pixel.NewSprite(asset, pixel.R(80, 16, 96, 32))
	num6 := pixel.NewSprite(asset, pixel.R(96, 16, 112, 32))
	num7 := pixel.NewSprite(asset, pixel.R(112, 16, 128, 32))
	num8 := pixel.NewSprite(asset, pixel.R(0, 32, 16, 48))
	num9 := pixel.NewSprite(asset, pixel.R(16, 32, 32, 48))

	woodSpr := pixel.NewSprite(asset, pixel.R(0, 0, 16, 16))
	youSpr := pixel.NewSprite(asset, pixel.R(16, 0, 32, 16))
	invSpr := pixel.NewSprite(asset, pixel.R(32, 0, 48, 16))
	grassSpr := pixel.NewSprite(asset, pixel.R(48, 0, 64, 16))
	draw := pixel.NewBatch(&pixel.TrianglesData{}, asset)
	var (
		player = playerPos{10, 10}
		frames int16 = 0
		second = time.Tick(time.Second)
		selSlot int8
		tilePos [21][21]Block
		randAngles = [4]float64{0, 1.5708, 3.14159, 4.71239}
		inv [10]ItemStack
		holdWASD [4]int16
		grassRotSeed = time.Now().UnixNano()
	)
	inv[0] = ItemStack{83, 1}
	inv[1] = ItemStack{3, 1}
	inv[4] = ItemStack{85, 1}
	tilePos[5][5] = Block{1, 99}
	rand.Seed(time.Now().UnixNano())

	draw.Clear()
	for !win.Closed() {
		//BEGIN CONTROLS
		//placing/mining bloks
		if win.Pressed(pixelgl.KeyUp) {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y+1].btype==0 {
				placeBlock(player.x, player.y+1, &tilePos, &inv, selSlot)
			} else {
				if !win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y+1].btype!=0 {
					tilePos[player.x][player.y+1].breakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyDown)&&player.y!=0 {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y-1].btype==0 {
				placeBlock(player.x, player.y-1, &tilePos, &inv, selSlot)
			} else {
				if !win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x][player.y-1].btype!=0 {
					tilePos[player.x][player.y-1].breakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x+1][player.y].btype==0 {
				placeBlock(player.x+1, player.y, &tilePos, &inv, selSlot)
			} else {
				if !win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x+1][player.y].btype!=0 {
					tilePos[player.x+1][player.y].breakStage++
				}
			}
		}
		if win.Pressed(pixelgl.KeyLeft)&&player.x!=0 {
			if win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x-1][player.y].btype==0 {
				placeBlock(player.x-1, player.y, &tilePos, &inv, selSlot)
			} else {
				if !win.Pressed(pixelgl.KeyLeftShift)&&tilePos[player.x-1][player.y].btype!=0 {
					tilePos[player.x-1][player.y].breakStage++
				}
			}
		}

		inv = deleteBrokenBlocks(&tilePos, inv, 1)
		//selecting slots
		x := win.Typed()
		if len(x)!=0 {
			x = x[len(x)-1:]
			if unicode.IsDigit([]rune(x)[0]) {
				i , _ := strconv.Atoi(x) 
				selSlot = int8(i-1)
			}
		}
		//moving
		//up
		if win.JustPressed(pixelgl.KeyW) {
			if moveCheck(player.x, player.y+1, tilePos) {
				player.y++
			}
		}
		if win.Pressed(pixelgl.KeyW) {
			if holdWASD[0]!=30 {
				holdWASD[0]++
			} else {
				holdWASD[0]-=5
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
			if holdWASD[1]!=20 {
				holdWASD[1]++
			} else {
				holdWASD[1]-=5
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
			if holdWASD[2]!=20 {
				holdWASD[2]++
			} else {
				holdWASD[2]-=5
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
			if holdWASD[3]!=20 {
				holdWASD[3]++
			} else {
				holdWASD[3]-=5
				if moveCheck(player.x+1, player.y, tilePos) {
					player.x++
				}
			}
		} else {
			holdWASD[3] = 0
		}
		//END CONTROLS
		//BEGIN RENDERING
		win.Clear(colornames.Forestgreen)
		draw.Clear()

		rand.Seed(grassRotSeed)
		for x := 16; x != 656; x += 32 {
			for y := 80; y != 720; y += 32 {
				grassSpr.Draw(draw, pixel.IM.Moved(pixel.V(float64(x), float64(y))).ScaledXY(pixel.V(float64(x), float64(y)), pixel.V(2, 2)).Rotated(pixel.V(float64(x), float64(y)), randAngles[rand.Intn(3)]))
			}
		}

		for x := 0; x != 20; x++ {
			for y := 0; y != 20; y++ {
				if tilePos[x][y].btype>0 {
					woodSpr.Draw(draw, pixel.IM.Moved(pixel.V(float64(x*32)+8, float64(y*32)+40)).ScaledXY(pixel.V(float64(x*32), float64(y*32)), pixel.V(2, 2)))
				}
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
			switch inv[i].itype {
			case 0:
			case 1:
				woodSpr.Draw(draw, pixel.IM.Moved(pixel.V(float64((i*64)+32), 32)).ScaledXY(pixel.V(float64((i*64)+32), 32), pixel.V(2.25, 2.25)))
				if inv[i].amnt<10 {
					switch inv[i].amnt {
					case 0:
						inv[i].itype = 0
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
					switch inv[i].amnt / 10 { //most signigf=iashudant digits
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

					switch inv[i].amnt % 10  { // least signifiicant digits
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
			default:
			}
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