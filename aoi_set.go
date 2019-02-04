// /////////////////////////////////////////////////////////////////////////////
// aoi 包模型

package aoi

// /////////////////////////////////////////////////////////////////////////////
// AoiSet 对象

// Aoi 保存对象
type AoiSet map[*Aoi]struct{}

// 添加1个 aoi 对象
func (this AoiSet) Add(aoi *Aoi) {
	this[aoi] = struct{}{}
}

// 移除1个 aoi 对象
func (this AoiSet) Remove(aoi *Aoi) {
	delete(this, aoi)
}

// 是否包含某个 aoi 对象
func (this AoiSet) Contains(aoi *Aoi) (ok bool) {
	_, ok = this[aoi]

	return
}
