// /////////////////////////////////////////////////////////////////////////////
// 灯塔

package aoi

import "log"

// /////////////////////////////////////////////////////////////////////////////
// tower 对象

// 灯塔对象
type tower struct {
	aoiObject map[*aoiObject]struct{} // 对象
	watchers  map[*aoiObject]struct{} // 观察者
}

// 数据初始化
func (this *tower) init() {
	this.aoiObject = map[*aoiObject]struct{}{}
	this.watchers = map[*aoiObject]struct{}{}
}

// 添加1个观察者
func (this *tower) addWatcher(watcher *aoiObject) {
	if _, ok := this.watchers[watcher]; ok {
		log.Panic("添加观察者异常：重复添加。")
	}

	this.watchers[watcher] = struct{}{}

	// 其他 aoiObject 进入观察者范围
	for obj := range this.aoiObject {
		if watcher == obj {
			continue
		}

		watcher.aoi.entity.OnAoiEnter(obj.aoi)
	}
}

// 移除1个观察者
func (this *tower) removeWatcher(watcher *aoiObject) {
	if _, ok := this.watchers[watcher]; !ok {
		log.Panic("移除观察者异常：该观察者不存在。")
	}

	delete(this.watchers, watcher)

	// 其他 aoiObject 离开观察者范围
	for obj := range this.aoiObject {
		if watcher == obj {
			continue
		}

		watcher.aoi.entity.OnAoiLeave(obj.aoi)
	}
}

// 添加1个对象
func (this *tower) addAoiObject(obj *aoiObject, otherTower *tower) {
	// 加入对象列表
	obj.tower = this
	this.aoiObject[obj] = struct{}{}

	// 通知观察者
	if nil == otherTower { // 不是从其他 tower 来的
		for watcher := range this.watchers {
			if watcher == obj {
				continue
			}

			watcher.aoi.entity.OnAoiEnter(obj.aoi)
		}
	} else { // 从其他 tower 来的

		// 旧塔： 将对象离开的消息，通知其他观察者
		for watcher := range otherTower.watchers {
			if watcher == obj {
				continue // 对象自己
			}

			if _, ok := this.watchers[watcher]; ok {
				continue // 观察者重合
			}

			watcher.aoi.entity.OnAoiLeave(obj.aoi)
		}

		// 新塔： 将对象进入的消息，通知其他观察者
		for watcher := range this.watchers {
			if watcher == obj {
				continue // 对象自己
			}

			if _, ok := otherTower.watchers[watcher]; ok {
				continue // 观察者重合
			}

			watcher.aoi.entity.OnAoiEnter(obj.aoi)
		}
	}
}

// 添加1个对象
func (this *tower) removeAoiObject(obj *aoiObject, notifyWatchers bool) {
	obj.tower = nil
	delete(this.aoiObject, obj)

	if notifyWatchers {
		for watcher := range this.watchers {
			if watcher == obj {
				continue
			}

			watcher.aoi.entity.OnAoiLeave(obj.aoi)
		}
	}
}
