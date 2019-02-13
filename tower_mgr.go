// /////////////////////////////////////////////////////////////////////////////
// 灯塔

package aoi

// /////////////////////////////////////////////////////////////////////////////
// TowerManager

// 灯塔管理
type TowerManager struct {
	startX     Coord     // x 轴起点
	endX       Coord     // x 轴终点
	startY     Coord     // y 轴起点
	endY       Coord     // y 轴终点
	towerRange Coord     // 灯塔范围
	towers     [][]tower // 灯塔数组
	xTowerNum  int       // x 轴方向灯塔数量
	yTowerNum  int       // y 轴方向灯塔数量
}

// 新建1个 TowerManager 对象
func NewTowerManager(startX, endX, startY, endY, towerRange Coord) IAoiManager {
	mgr := &TowerManager{
		startX:     startX,
		endX:       endX,
		startY:     startY,
		endY:       endY,
		towerRange: towerRange,
	}

	mgr.init()

	return mgr
}

// 某个 aoi 对象进入
func (this *TowerManager) Enter(aoi *Aoi, x, y Coord) {
	aoi.x = x
	aoi.y = y

	obj := &aoiObject{
		aoi: aoi,
	}

	aoi.aoiObject = obj

	// 添加观察者
	this.visitTowers(x, y, aoi.distance, func(t *tower) {
		t.addWatcher(obj)
	})

	// 添加对象
	t := this.getTower(x, y)
	t.addAoiObject(obj, nil)
}

// 某个 aoi 对象离开
func (this *TowerManager) Leave(aoi *Aoi) {
	// 移除对象
	obj := aoi.aoiObject
	obj.tower.removeAoiObject(obj, true)

	// 移除观察者
	this.visitTowers(aoi.x, aoi.y, aoi.distance, func(t *tower) {
		t.removeWatcher(obj)
	})
}

// 某个 aoi 对象移动
func (this *TowerManager) Moved(aoi *Aoi, x, y Coord) {
	// 位置更新
	oldx, oldy := aoi.x, aoi.y
	aoi.x, aoi.y = x, y

	// 灯塔变化
	oldTower := aoi.aoiObject.tower
	newTower := this.getTower(x, y)

	if oldTower != newTower {
		oldTower.removeAoiObject(aoi.aoiObject, false)

		newTower.addAoiObject(aoi.aoiObject, oldTower)
	}

	// 新/旧 观察范围
	oldMinXt, oldMaxXt, oldMinYt, oldMaxYt := this.getWatchRange(oldx, oldy, aoi.distance)
	newMinXt, newMaxXt, newMinYt, newMaxYt := this.getWatchRange(x, y, aoi.distance)

	// 观察范围：离开某个旧塔
	for i := oldMinXt; i <= oldMaxXt; i++ {
		for j := oldMinYt; j <= oldMaxYt; j++ {
			if i >= newMinXt && i <= newMaxXt && j >= newMinYt && j <= newMaxYt {
				continue
			}

			t := &this.towers[i][j]

			t.removeWatcher(aoi.aoiObject)
		}
	}

	// 观察范围：进入某个新塔
	for i := newMinXt; i <= newMaxXt; i++ {
		for j := newMinYt; j <= newMaxYt; j++ {
			if i >= oldMinXt && i <= oldMaxXt && j >= oldMinYt && j <= oldMaxYt {
				continue
			}

			t := &this.towers[i][j]

			t.addWatcher(aoi.aoiObject)
		}
	}
}

// 初始化灯塔数据
func (this *TowerManager) init() {
	numXSlots := int((this.endX-this.startX)/this.towerRange) + 1
	this.xTowerNum = numXSlots

	numYSlots := int((this.endY-this.startY)/this.towerRange) + 1
	this.yTowerNum = numYSlots

	this.towers = make([][]tower, numXSlots)

	for i := 0; i < numXSlots; i++ {
		this.towers[i] = make([]tower, numYSlots)
		for j := 0; j < numYSlots; j++ {
			this.towers[i][j].init()
		}
	}
}

// 遍历以 x,y 为中心，aoiDistance 为半径范围内的所有灯塔
func (this *TowerManager) visitTowers(x, y, aoiDistance Coord, f func(t *tower)) {
	minXt, maxXt, minYt, maxYt := this.getWatchRange(x, y, aoiDistance)

	for i := minXt; i < maxXt; i++ {
		for j := minYt; j < maxYt; j++ {
			tower := &this.towers[i][j]

			f(tower)
		}
	}
}

// 根据 x,y aoiDistance 计算观察范围
func (this *TowerManager) getWatchRange(x, y, aoiDistance Coord) (int, int, int, int) {
	minXt, minYt := this.getTowerPos(x-aoiDistance, y-aoiDistance)
	maxXt, maxYt := this.getTowerPos(x+aoiDistance, y+aoiDistance)

	return minXt, maxXt, minYt, maxYt
}

// 计算 x,y 所处灯塔坐标
func (this *TowerManager) getTowerPos(x, y Coord) (int, int) {
	xt := int((x - this.startX) / this.towerRange)
	yt := int((y - this.startY) / this.towerRange)

	return this.normalizeXt(xt), this.normalizeYt(yt)
}

// 修正 xt
func (this *TowerManager) normalizeXt(xt int) int {
	if xt < 0 {
		xt = 0
	} else if xt >= this.xTowerNum {
		xt = this.xTowerNum - 1
	}

	return xt
}

// 修正 yt
func (this *TowerManager) normalizeYt(yt int) int {
	if yt < 0 {
		yt = 0
	} else if yt >= this.yTowerNum {
		yt = this.yTowerNum - 1
	}

	return yt
}

// 获取 x,y 坐标对应的灯塔指针
func (this *TowerManager) getTower(x, y Coord) *tower {
	xt, yt := this.getTowerPos(x, y)

	return &this.towers[xt][yt]
}
