package lib

type ItemStack struct {
	Amnt int8
	Itype uint16
}

func GiveItem(inv *[10]ItemStack, ite ItemStack) {
	iltd:=ite.Amnt
	for i := 0; i < 10; i++ { 
		if inv[i].Itype==ite.Itype {
			inv[i].Itype=ite.Itype
			for iltd!=0&&inv[i].Amnt<85 {
				inv[i].Amnt+=1
				iltd-=1
			}
		}
	}
	for i := 0; i < 10; i++ {
		if inv[i].Itype==ite.Itype||inv[i].Itype==0 {
			if inv[i].Amnt<85&&iltd>0 {
				inv[i].Itype=ite.Itype
				for iltd!=0&&inv[i].Amnt<85 {
					inv[i].Amnt+=1
					iltd-=1
				}
			}
		}
	}
}

func AddItems(inv [10]ItemStack) ([65535]int64) {
	var retList [65535]int64
	for i := 0; i < 10; i++ {
		retList[inv[i].Itype] += int64(inv[i].Amnt)
	}
	return retList
}

func IS(Amnt int8, Itype uint16) (ItemStack) {
	return ItemStack{Amnt, Itype}
}