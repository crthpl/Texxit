package lib

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"fmt"
	"golang.org/x/image/colornames"
)

type selCraft struct {
	x, y int
}

func CraftGUI(gui *int16, inv *[10]ItemStack, asset pixel.Picture, pixs [8][8]*pixel.Sprite, opics [100]*pixel.Sprite) {
	craftCfg := pixelgl.WindowConfig{ //the settings for the window
		Title:  "Crafting",
		Bounds: pixel.R(0, 0, 768, 320),
		VSync:  true,
		AlwaysOnTop: true,
	}

	Cwin, err := pixelgl.NewWindow(craftCfg)
	if err != nil {
		panic(err)
	}

	invSpr := pixel.NewSprite(asset, pixel.R(32, 0, 48, 16))
	tiles := pixel.NewBatch(&pixel.TrianglesData{}, asset)
	crafts := ReadCrafts("crafts.json")
	for !Cwin.Closed() {
		Cwin.Clear(colornames.Forestgreen)
		selCarft:=selCraft{0, 0}
		if Cwin.JustPressed(pixelgl.KeyC) {
			break
		}
		var ic int
		var mov pixel.Matrix
		var savedic int
		for x := 0; x < 12; x++ {
			for y := 0; y < 5; y++ {
				ic++
				if len(crafts.Crafts)>=ic {
					if !(x==selCarft.x&&y==selCarft.y) {
						invSpr.Draw(tiles, pixel.IM.Moved(pixel.V(float64(x*64)+32, float64(y*64)+32)).Scaled(pixel.V(float64(x*64)+32, float64(y*64)+32), 4))
						opics[crafts.Crafts[ic-1].Result.Itype].Draw(tiles, pixel.IM.Moved(pixel.V(float64(x*64)+32, float64(y*64)+32)).Scaled(pixel.V(float64(x*64)+32, float64(y*64)+32), 2))
					} else {
						savedic = ic
						mov = pixel.IM.Moved(pixel.V(float64(x*64)+32, float64(y*64)+32)).Scaled(pixel.V(float64(x*64)+32, float64(y*64)+32), 2.5)
					}
					
				}
			}
		}

		invSpr.Draw(tiles, pixel.IM.Moved(pixel.V(float64(selCarft.x*64)+32, float64(selCarft.y*64)+32)).Scaled(pixel.V(float64(selCarft.x*64)+32, float64(selCarft.y*64)+32), 5))
		opics[crafts.Crafts[savedic-1].Result.Itype].Draw(tiles, mov)
		adi := AddItems(*inv)
		craft1 := crafts.Crafts[0]
		for _, v := range craft1.Reqs {
			if int8(adi[v.Itype])>=v.Amnt {
				
			} else {
				fmt.Println("nop")
			}
		}
		tiles.Draw(Cwin)
		Cwin.Update()
	}
	Cwin.Destroy()
	*gui = 0
}
