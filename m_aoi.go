// /////////////////////////////////////////////////////////////////////////////
// aoi 包模型

package aoi

// /////////////////////////////////////////////////////////////////////////////
// 接口

//
type IAoiCallback interface {
	OnEnterAoi(other *Aoi)
	OnLeaveAoi(other *Aoi)
}

// aoi 管理接口
type IAoiManager interface {
	Enter(aoi *Aoi, x, y Coord) // 进入
	Leave(aoi *Aoi)             // 离开
	Moved(aoi *Aoi, x, y Coord) // 移动
}
