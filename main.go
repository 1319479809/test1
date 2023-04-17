package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义接收数据的结构体 9.
type SendPost struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Type     int    `form:"type" json:"type"  ` //1：刷脸 2：扫码 3：刷卡
	Data     string `form:"data" json:"data"  ` //扫描二维码获取到的文本信息、刷卡卡号和刷脸人员 ID
	Time     int64  `form:"time" json:"time"`
	DeviceSn string `form:"deviceSn" json:"deviceSn"`
}

type RevPost struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Code    int    `form:"code" json:"code" uri:"code" xml:"code" `             //code < 0 ：失败 code >= 0：成功
	Cmd     int    `form:"cmd" json:"cmd" uri:"cmd" xml:"cmd""`                 //1：输出继电器
	Message string `form:"message" json:"message" uri:"message" xml:"message" ` //响应描述
}

type ReturnResult struct { //返回
	Msg     string `json:"msg"`
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Result  int    `json:"result"`
}

// 远程控制设备
type SendControl struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Pass string          `json:"Pass"` //设备密码
	Data SendControlData `json:"data"` //详见<Data 数据>表
}
type SendControlData struct {
	Commend   int    `json:"command"`   //1：远程开门，2：播放音频， 100：重启设备
	VoiceData string `json:"voiceData"` //仅在 command=2 生效

}

// 人员注册
type SendPersonCreat struct {
	Pass                     string `form:"pass" json:"pass"`                                         //设备密码
	EmployeeNumber           string `form:"employeeNumber" json:"employeeNumber"`                     //人员 Id
	Name                     string `form:"name" json:"name"`                                         //人员名称
	Gender                   string `form:"gender" json:"gender"`                                     //性别
	Nationa                  string `form:"nationa" json:"nationa"`                                   //国家
	DepartmentName           string `form:"departmentName" json:"departmentName"`                     //部门名称
	IdCardNumber             string `form:"idCardNumber" json:"idCardNumber"`                         //人员卡号
	Mobile                   string `form:"mobile" json:"mobile"`                                     //手机号
	AccessRight              int    `form:"accessRight" json:"accessRight"`                           // 权限模式
	TemporaryAccessStartTime int    `form:"temporaryAccessStartTime" json:"temporaryAccessStartTime"` //开始时间戳
	TemporaryAccessEndTime   int    `form:"temporaryAccessEndTime" json:"temporaryAccessEndTime"`     //结束时间戳
	TemporaryAccessTimes     int    `form:"temporaryAccessTimes" json:"temporaryAccessTimes"`         //通行次数
	AccessCardNumber         string `form:"accessCardNumber" json:"accessCardNumber"`                 //门禁卡号
	Remarks                  string `form:"remarks" json:"remarks"`                                   //备注
	PhotoFromCapture         int    `form:"photoFromCapture" json:"photoFromCapture"`                 //拍照注册
	Replace                  int    `form:"replace" json:"replace"`                                   //替换标志
	RegisterBase64           string `form:"registerBase64" json:"registerBase64"`                     //照片 base64 数据
}

// 人员删除
type SendPersonDelete struct {
	Pass           string `form:"pass" json:"pass"  `                    //设备密码
	EmployeeNumber string `form:"employeeNumber" json:"employeeNumber" ` //人员 Id
}

// 人员查询
type SendPersonFind struct {
	Pass           string `json:"pass"`            //设备密码
	PicLarge       int    `json:"picLarge"`        //是否返回注册时照片，0 或者不传不返回， 1：返回
	EmployeeNumber string `json:"employeeNumber" ` //人员 Id
	Name           string `json:"name"`            //人员姓名
	Length         int    `json:"length"`          //每页最大数量
	Index          int    `json:"index"`           //页码
}

// 同步白名单
type SendPersonWhiteListSync struct {
	Pass       string          `form:"pass" json:"pass"  ` //设备密码
	SyncFlag   int             `form:"syncFlag" json:"syncFlag"`
	WhiteLists []WhiteListSync `form:"whiteLists" json:"whiteLists" ` //人员 Id
}
type WhiteListSync struct {
	UserType        int    `form:"userType" json:"userType" `              //101-二维码，202-卡，303-人脸
	UserId          string `form:"userId" json:"userId"`                   //页码
	BeginTime       int    `form:"beginTime" json:"beginTime"`             //时间戳
	EndTime         int    `form:"endTime" json:"endTime"`                 //时间戳
	RepeatType      int    `form:"repeatType" json:"repeatType"`           //小于等于 0-不重复，1-每天重复，2-每周重 复，3-每月重复
	RepeatBeginTime int    `form:"repeatBeginTime" json:"repeatBeginTime"` //1：开始时间为距 0 点的时间 2：开始时间为周几 3：开始时间为某日
	RepeatEndTime   int    `form:"repeatEndTime" json:"repeatEndTime"`     //1：结束时间为距 0 点的时间 2：结束时间为周几 3：结束时间为某日
	SyncType        int    `form:"syncType" json:"syncType"`               //当 syncFlag=2 时有效，1-增加，2-修改，3- 删 除
}

// 查询白名单
type SendPersonWhiteListFind struct {
	Pass      string `form:"pass" json:"pass" `          //设备密码
	UserType  int    `form:"userType" json:"userType"  ` //101-二维码，202-卡，303-人脸
	UserId    string `form:"userId" json:"userId"`       //页码
	BeginTime int    `form:"beginTime" json:"beginTime"` //时间戳
	EndTime   int    `form:"endTime" json:"endTime"`     //时间戳
	Length    int    `form:"length" json:"length"`       //每页最大数量
	Index     int    `form:"index" json:"index"`         //页码
}

// 人员注册（feature）
type SendPersonRegisterFeats struct {
	Pass  string               `form:"pass" json:"pass" `   //设备密码
	Users []PersonRegisterData `form:"users" json:"users" ` //设备密码
}
type PersonRegisterData struct {
	EmployeeNumber           []string `form:"employeeNumber" json:"employeeNumber"  `                   //人员 Id
	Name                     string   `form:"name" json:"name"`                                         //人员名称
	Feature                  string   `form:"feature" json:"feature" "`                                 //人员特征值 16 进制字符串长度 2048 字节
	Gender                   string   `form:"gender" json:"gender"`                                     //性别
	Nationa                  string   `form:"nationa" json:"nationa"`                                   //国家
	DepartmentName           string   `form:"departmentName" json:"departmentName"`                     //部门名称
	IdCardNumber             string   `form:"idCardNumber" json:"idCardNumber"`                         //人员卡号
	Mobile                   string   `form:"mobile" json:"mobile"`                                     //手机号
	AccessRight              int      `form:"accessRight" json:"accessRight"`                           // 权限模式
	TemporaryAccessStartTime int      `form:"temporaryAccessStartTime" json:"temporaryAccessStartTime"` //开始时间戳
	TemporaryAccessEndTime   int      `form:"temporaryAccessEndTime" json:"temporaryAccessEndTime"`     //结束时间戳
	TemporaryAccessTimes     int      `form:"temporaryAccessTimes" json:"temporaryAccessTimes"`         //通行次数
	AccessCardNumber         string   `form:"accessCardNumber" json:"accessCardNumber"`                 //门禁卡号
	Remarks                  string   `form:"remarks" json:"remarks"`                                   //备注
	Replace                  int      `form:"replace" json:"replace"`                                   //替换标志
}

var DeviceUrl = "http://192.168.31.171:8060"

// http://192.168.2.19:80/sendTest
func main() {
	r := gin.Default()
	v1 := r.Group("/device")
	{
		v1.POST("/control", deviceControl) //远程控制设备
	}
	v2 := r.Group("/person")
	{
		v2.POST("/create", personCreate)               //人员注册.
		v2.POST("/delete", personDelete)               //人员删除
		v2.POST("/find", personFind)                   //人员查询
		v2.POST("/whiteListSync", personWhiteListSync) //同步白名单
		v2.POST("/whiteListFind", personWhiteListFind) //查询白名单
		v2.POST("/registerFeats", personRegisterFeats) //人员注册（feature）

	}
	v3 := r.Group("/v1")
	{
		v3.POST("/post", sendPost)
	}
	//定义默认路由
	r.NoRoute(func(c *gin.Context) {
		fmt.Println("test4")
		c.JSON(http.StatusNotFound, gin.H{
			"status": 0,
			"error":  "success",
		})
	})
	r.Run(":8080")
}

// http://127.0.0.1:8080/v1/post
func sendPost(c *gin.Context) {
	//type := c.DefaultQuery("type", "type")
	var SJson SendPost
	var rJson RevPost
	fmt.Println("sendPost ", c.Request)
	fmt.Println("SJson=", SJson)
	if err := c.ShouldBindJSON(&SJson); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//c.String(-1, "false")
		rJson.Code = -1
		rJson.Message = err.Error()
		fmt.Println("rJson=", rJson)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rJson.Code = 1
	rJson.Cmd = 1
	rJson.Message = "success"

	fmt.Println("rJson=", rJson)
	c.JSON(http.StatusOK, rJson)
	//c.String(0, "success")

}

// 设备控制
func deviceControl(c *gin.Context) {
	//urlP := "http://192.168.31.171:8060/device/control"
	urlP := DeviceUrl + "/device/control"
	fmt.Println("urlP=", urlP)
	// fmt.Println("c=", c.Request.Body)
	// 参数
	var SControl SendControl
	if err := c.ShouldBindJSON(&SControl); err != nil {

		fmt.Println("SControl=", SControl)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("SControl 1=", SControl)

	bytesData, _ := json.Marshal(SControl)
	fmt.Println("SControl 2=", SControl)
	// 使用http.Post请求
	resp, err := http.Post(urlP, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println("err2=", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err3=", err.Error())
		return
	}

	// 打印返回结果
	fmt.Println(string(body), resp.StatusCode)

	c.String(200, string(body)) //直接返回body

}

// 人员注册
func personCreate(c *gin.Context) {
	urlP := DeviceUrl + "/person/create"
	fmt.Println("create urlP=", urlP)
	// 参数
	var SPersonCreat SendPersonCreat
	if err := c.ShouldBindJSON(&SPersonCreat); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("SPersonCreat  1=", SPersonCreat)

	bytesData, _ := json.Marshal(SPersonCreat)
	fmt.Println("SPersonCreat  2=", SPersonCreat)
	// 使用http.Post请求
	resp, err := http.Post(urlP, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println("SPersonCreat err2=", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("SPersonCreat err3=", err.Error())
		return
	}

	// 打印返回结果
	fmt.Println(string(body), resp.StatusCode)

	c.String(200, string(body)) //直接返回body
}

// 人员删除
func personDelete(c *gin.Context) {
	urlP := DeviceUrl + "/person/delete"
	fmt.Println("urlP=", urlP)
	// 参数
	var SPersonDelete SendPersonDelete
	if err := c.ShouldBindJSON(&SPersonDelete); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("SPersonDelete 1=", SPersonDelete)

	bytesData, _ := json.Marshal(SPersonDelete)
	fmt.Println("SPersonDelete 2=", SPersonDelete)
	// 使用http.Post请求
	resp, err := http.Post(urlP, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println("SPersonDelete err2=", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("SPersonDelete err3=", err.Error())
		return
	}

	// 打印返回结果
	fmt.Println(string(body), resp.StatusCode)

	c.String(200, string(body)) //直接返回body
}

// 人员查询
func personFind(c *gin.Context) {
	urlP := DeviceUrl + "/person/find"
	fmt.Println("urlP=", urlP)
	// fmt.Println("personFind ", c.Request)
	// 参数
	var SPersonFind SendPersonFind
	if err := c.ShouldBindJSON(&SPersonFind); err != nil {

		fmt.Println("urlP=", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("SPersonFind 1=", SPersonFind)

	bytesData, _ := json.Marshal(SPersonFind)
	fmt.Println("SPersonFind 2=", SPersonFind)
	// 使用http.Post请求
	resp, err := http.Post(urlP, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println("SPersonFind err2=", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("SPersonFind err3=", err.Error())
		return
	}

	// 打印返回结果
	fmt.Println(string(body), resp.StatusCode)

	c.String(200, string(body)) //直接返回body
}

// 同步白名单
func personWhiteListSync(c *gin.Context) {
	urlP := DeviceUrl + "/person/whiteListSync"
	fmt.Println("urlP=", urlP)
	// 参数
	var SPersonWhiteListSync SendPersonWhiteListSync
	if err := c.ShouldBindJSON(&SPersonWhiteListSync); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("SPersonWhiteListSync 1=", SPersonWhiteListSync)

	bytesData, _ := json.Marshal(SPersonWhiteListSync)
	fmt.Println("SPersonWhiteListSync 2=", SPersonWhiteListSync)
	// 使用http.Post请求
	resp, err := http.Post(urlP, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println("SPersonWhiteListSync err2=", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("SPersonWhiteListSync err3=", err.Error())
		return
	}

	// 打印返回结果
	fmt.Println(string(body), resp.StatusCode)

	c.String(200, string(body)) //直接返回body
}

// 查询白名单
func personWhiteListFind(c *gin.Context) {
	urlP := DeviceUrl + "/person/whiteListFind"
	fmt.Println("urlP=", urlP)
	// 参数
	var SPersonWhiteListFind SendPersonWhiteListFind
	if err := c.ShouldBindJSON(&SPersonWhiteListFind); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("SPersonWhiteListFind 1=", SPersonWhiteListFind)

	bytesData, _ := json.Marshal(SPersonWhiteListFind)
	fmt.Println("SPersonWhiteListSync 2=", SPersonWhiteListFind)
	// 使用http.Post请求
	resp, err := http.Post(urlP, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println("SPersonWhiteListFind err2=", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("SPersonWhiteListFind err3=", err.Error())
		return
	}

	// 打印返回结果
	fmt.Println(string(body), resp.StatusCode)

	c.String(200, string(body)) //直接返回body
}

// 人员注册（feature）
func personRegisterFeats(c *gin.Context) {
	urlP := DeviceUrl + "/person/registerFeats"
	fmt.Println("urlP=", urlP)
	// 参数
	var SPersonRegisterFeats SendPersonRegisterFeats
	if err := c.ShouldBindJSON(&SPersonRegisterFeats); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("SPersonRegisterFeats 1=", SPersonRegisterFeats)

	bytesData, _ := json.Marshal(SPersonRegisterFeats)
	fmt.Println("SPersonWhiteListSync 2=", SPersonRegisterFeats)
	// 使用http.Post请求
	resp, err := http.Post(urlP, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println("SPersonRegisterFeats err2=", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("SPersonRegisterFeats err3=", err.Error())
		return
	}

	// 打印返回结果
	fmt.Println(string(body), resp.StatusCode)

	c.String(200, string(body)) //直接返回body
}
