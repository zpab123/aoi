// /////////////////////////////////////////////////////////////////////////////
// Aoi 包模型

package aoi

// /////////////////////////////////////////////////////////////////////////////
// public api

// 对1个 *Aoi 对象初始化数据
func InitAOI(a *Aoi, dist Coord, data interface{}, callback IAoiCallback) {
	a.distance = dist
	a.Data = data
	a.callback = callback
}

// /////////////////////////////////////////////////////////////////////////////
// aoi 对象

// 兴趣区域对象
type Aoi struct {
	x         Coord       // x 轴坐标点
	y         Coord       // y 轴坐标点
	distance  Coord       // 感应距离
	entity    IAoiEntity  // 具备 aoi 特性的实体
	aoiObject *aoiObject  // aoiObject 对象
	Data      interface{} // 数据
}
