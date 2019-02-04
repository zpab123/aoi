// /////////////////////////////////////////////////////////////////////////////
// Aoi 包模型

package aoi

// 坐标值数值类型
type Coord float32

// /////////////////////////////////////////////////////////////////////////////
// public api

// 对1个 *Aoi 对象初始化数据
func InitAOI(a *Aoi, dist Coord, data interface{}, callback IAoiCallback) {
	a.dist = dist
	a.Data = data
	a.callback = callback
}

// /////////////////////////////////////////////////////////////////////////////
// Aoi 包模型

// aoi 对象
type Aoi struct {
	x        Coord        // x 轴坐标点
	y        Coord        // y 轴坐标点
	dist     Coord        // 距离
	callback IAoiCallback // 未知
	implData interface{}  // 未知
	Data     interface{}  // 数据
}
