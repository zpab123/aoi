// /////////////////////////////////////////////////////////////////////////////
// 灯塔

package aoi

import "log"

// /////////////////////////////////////////////////////////////////////////////
// TowerAoiManager

// 灯塔管理
type TowerAoiManager struct {
	minX       Coord     // x 轴起点
	maxX       Coord     // x 轴终点
	minY       Coord     // y 轴起点
	maxY       Coord     // y 轴终点
	towerRange Coord     // 灯塔范围
	towers     [][]tower // 灯塔数组
	xTowerNum  int       // x 轴方向灯塔数量
	yTowerNum  int       // y 轴方向灯塔数量
}

// 新建1个 TowerAoiManager 对象
func NewTowerAoiManager(minX, maxX, minY, maxY Coord, towerRange Coord) AOIManager {
	aoiman := &TowerAoiManager{minX: minX, maxX: maxX, minY: minY, maxY: maxY, towerRange: towerRange}
	aoiman.init()

	return aoiman
}

// 某个 aoi 对象进入
func (this *TowerAoiManager) Enter(aoi *Aoi, x, y Coord) {
	aoi.x, aoi.y = x, y
	obj := &aoiobj{aoi: aoi}
	aoi.implData = obj

	// 添加观察者
	this.visitWatchedTowers(x, y, aoi.dist, func(tower *tower) {
		tower.addWatcher(obj)
	})

	// 添加对象
	t := this.getTowerXY(x, y)
	t.addObj(obj, nil)
}

// 某个 aoi 对象离开
func (this *TowerAoiManager) Leave(aoi *Aoi) {
	obj := aoi.implData.(*aoiobj)
	obj.tower.removeObj(obj, true)

	this.visitWatchedTowers(aoi.x, aoi.y, aoi.dist, func(tower *tower) {
		tower.removeWatcher(obj)
	})
}

// 某个 aoi 对象移动
func (this *TowerAoiManager) Moved(aoi *Aoi, x, y Coord) {
	oldx, oldy := aoi.x, aoi.y
	aoi.x, aoi.y = x, y
	obj := aoi.implData.(*aoiobj)
	t0 := obj.tower
	t1 := this.getTowerXY(x, y)

	if t0 != t1 {
		t0.removeObj(obj, false)
		t1.addObj(obj, t0)
	}

	oximin, oximax, oyimin, oyimax := this.getWatchedTowers(oldx, oldy, aoi.dist)
	ximin, ximax, yimin, yimax := this.getWatchedTowers(x, y, aoi.dist)

	// 通知旧观察者
	for xi := oximin; xi <= oximax; xi++ {
		for yi := oyimin; yi <= oyimax; yi++ {
			if xi >= ximin && xi <= ximax && yi >= yimin && yi <= yimax {
				continue
			}

			tower := &this.towers[xi][yi]
			tower.removeWatcher(obj)
		}
	}

	// 添加新观察者
	for xi := ximin; xi <= ximax; xi++ {
		for yi := yimin; yi <= yimax; yi++ {
			if xi >= oximin && xi <= oximax && yi >= oyimin && yi <= oyimax {
				continue
			}

			tower := &this.towers[xi][yi]
			tower.addWatcher(obj)
		}
	}
}

// 计算 x,y 所处灯塔坐标
func (this *TowerAoiManager) transXY(x, y Coord) (int, int) {
	xi := int((x - this.minX) / this.towerRange)
	yi := int((y - this.minY) / this.towerRange)

	return this.normalizeXi(xi), this.normalizeYi(yi)
}

// 修正 xi
func (this *TowerAoiManager) normalizeXi(xi int) int {
	if xi < 0 {
		xi = 0
	} else if xi >= this.xTowerNum {
		xi = this.xTowerNum - 1
	}

	return xi
}

// 修正 yi
func (this *TowerAoiManager) normalizeYi(yi int) int {
	if yi < 0 {
		yi = 0
	} else if yi >= this.yTowerNum {
		yi = this.yTowerNum - 1
	}

	return yi
}

// 获取 x,y 坐标对应的灯塔指针
func (this *TowerAoiManager) getTowerXY(x, y Coord) *tower {
	xi, yi := this.transXY(x, y)

	return &this.towers[xi][yi]
}

// 根据 x,y aoiDistance 计算 灯塔坐标范围
func (this *TowerAoiManager) getWatchedTowers(x, y Coord, aoiDistance Coord) (int, int, int, int) {
	ximin, yimin := this.transXY(x-aoiDistance, y-aoiDistance)
	ximax, yimax := this.transXY(x+aoiDistance, y+aoiDistance)
	//aoiTowerNum := int(aoiDistance/this.towerRange) + 1
	//ximid, yimid := this.transXY(x, y)
	//ximin, ximax := this.normalizeXi(ximid-aoiTowerNum), this.normalizeXi(ximid+aoiTowerNum)
	//yimin, yimax := this.normalizeYi(yimid-aoiTowerNum), this.normalizeYi(yimid+aoiTowerNum)
	return ximin, ximax, yimin, yimax
}

// 遍历 x,y aoiDistance 范围 内的所有灯塔
func (this *TowerAoiManager) visitWatchedTowers(x, y Coord, aoiDistance Coord, f func(*tower)) {
	ximin, ximax, yimin, yimax := this.getWatchedTowers(x, y, aoiDistance)
	for xi := ximin; xi <= ximax; xi++ {
		for yi := yimin; yi <= yimax; yi++ {
			tower := &this.towers[xi][yi]
			f(tower)
		}
	}
}

// 初始化灯塔数据
func (this *TowerAoiManager) init() {
	numXSlots := int((this.maxX-this.minX)/this.towerRange) + 1
	this.xTowerNum = numXSlots
	numYSlots := int((this.maxY-this.minY)/this.towerRange) + 1
	this.yTowerNum = numYSlots

	this.towers = make([][]tower, numXSlots)
	for i := 0; i < numXSlots; i++ {
		this.towers[i] = make([]tower, numYSlots)
		for j := 0; j < numYSlots; j++ {
			this.towers[i][j].init()
		}
	}
}

// /////////////////////////////////////////////////////////////////////////////
// tower

// 灯塔对象
type tower struct {
	objs     map[*aoiobj]struct{} // 对象
	watchers map[*aoiobj]struct{} // 观察者
}

// 数据初始化
func (this *tower) init() {
	this.objs = map[*aoiobj]struct{}{}
	this.watchers = map[*aoiobj]struct{}{}
}

// 添加1个对象
func (this *tower) addObj(obj *aoiobj, fromOtherTower *tower) {
	// 加入对象列表
	obj.tower = this
	this.objs[obj] = struct{}{}

	// 通知观察者
	if fromOtherTower == nil {
		for watcher := range this.watchers {
			if watcher == obj {
				continue
			}
			watcher.aoi.callback.OnEnterAoi(obj.aoi)
		}
	} else {
		// 从其他 tower 移动到本 tower
		for watcher := range fromOtherTower.watchers {
			if watcher == obj {
				continue // 对象自己
			}

			if _, ok := this.watchers[watcher]; ok {
				continue // 观察者重合
			}

			watcher.aoi.callback.OnLeaveAoi(obj.aoi) // 通知离开
		}

		// 通知本 tower 进入消息
		for watcher := range this.watchers {
			if watcher == obj {
				continue // 对象自己
			}

			if _, ok := fromOtherTower.watchers[watcher]; ok {
				continue // 观察者重合
			}

			watcher.aoi.callback.OnEnterAoi(obj.aoi) // 通知进入
		}
	}
}

// 移除1个 对象
func (this *tower) removeObj(obj *aoiobj, notifyWatchers bool) {
	obj.tower = nil
	delete(this.objs, obj)

	if notifyWatchers {
		for watcher := range this.watchers {
			if watcher == obj {
				continue
			}

			watcher.aoi.callback.OnLeaveAoi(obj.aoi)
		}
	}
}

// 添加1个观察者
func (this *tower) addWatcher(obj *aoiobj) {
	if _, ok := this.watchers[obj]; ok {
		log.Panicf("添加观察者异常：重复添加")
	}

	this.watchers[obj] = struct{}{}

	// 通知其他对象
	for neighbor := range this.objs {
		if neighbor == obj {
			continue
		}

		obj.aoi.callback.OnEnterAoi(neighbor.aoi)
	}
}

// 移除1个观察者
func (this *tower) removeWatcher(obj *aoiobj) {
	if _, ok := this.watchers[obj]; !ok {
		log.Panicf("添加观察者异常：观察者不存在")
	}

	delete(this.watchers, obj)

	// 通知其他对象
	for neighbor := range this.objs {
		if neighbor == obj {
			continue
		}

		obj.aoi.callback.OnLeaveAoi(neighbor.aoi)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// aoiobj

// 灯塔监控对象
type aoiobj struct {
	aoi   *Aoi   // aoi 对象
	tower *tower // 灯塔对象
}
