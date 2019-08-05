package wxpay

/**
- XMLName字段，如上所述，会省略
- 具有标签"-"的字段会省略
- 具有标签"name,attr"的字段会成为该XML元素的名为name的属性
- 具有标签",attr"的字段会成为该XML元素的名为字段名的属性
- 具有标签",chardata"的字段会作为字符数据写入，而非XML元素
- 具有标签",innerxml"的字段会原样写入，而不会经过正常的序列化过程
- 具有标签",comment"的字段作为XML注释写入，而不经过正常的序列化过程，该字段内不能有"--"字符串
- 标签中包含"omitempty"选项的字段如果为空值会省略
  空值为false、0、nil指针、nil接口、长度为0的数组、切片、映射
- 匿名字段（其标签无效）会被处理为其字段是外层结构体的字段
- 如果一个字段的标签为"a>b>c"，则元素c将会嵌套进其上层元素a和b中。如果该字段相邻的字段标签指定了同样的上层元素，则会放在同一个XML元素里。
*/

type Request struct {
	AppId    string `xml:"app_id"`
	MchId    string `xml:"mac_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
	SignType string `xml:"sign_type,omitempty"`
}

type Response struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"app_id"`
	MchId      string `xml:"mac_id"`
	DeviceInfo string `xml:"device_info,omitempty"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
}
