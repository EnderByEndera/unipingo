package main

import (
	"context"
	"fmt"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"melodie-site/server/utils"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ImportMajors() {
	majors := make([]map[string]interface{}, 0)
	err := utils.UnmarshallJSONFromFile("majors.json", &majors)
	if err != nil {
		panic(err)
	}
	for i, majorProps := range majors {
		// majorProps

		major := models.Major{
			Code:       fmt.Sprintf("%+v", majorProps["专业代码"]),
			Name:       fmt.Sprintf("%+v", majorProps["专业名称"]),
			FirstLevel: fmt.Sprintf("%+v", majorProps["门类"]),
			Category:   fmt.Sprintf("%+v", majorProps["大类"]),
		}
		fmt.Printf("%d %+v\n", i, major)

		filter := bson.M{"name": major.Name}
		res := db.GetCollection("majors").FindOne(context.TODO(), filter)
		var insertedDocID primitive.ObjectID
		if res.Err() != nil {
			fmt.Println("err", res.Err())
			res, err := db.GetCollection("majors").InsertOne(context.TODO(), major)
			if err != nil {
				panic(err)
			}
			insertedDocID = res.InsertedID.(primitive.ObjectID)
		} else {
			_major := models.Major{}
			err = res.Decode(&_major)
			insertedDocID = _major.ID
			if err != nil {
				panic(err)
			}
		}
		fmt.Println(insertedDocID)

	}
}

func GetHEILevel(level string) models.HEILevel {
	if level == "本科" {
		return models.UniversityHEI
	} else {
		return models.VocationalHEI
	}
}

func GetHEIMode(mode string) models.HEIMode {
	if mode == "民办" {
		return models.PrivateHEI
	} else if mode == "内地与港澳台地区合作办学" {
		return models.Mainland_HK_MO_TW_Cooperative
	} else if mode == "中外合作办学" {
		return models.SinoForeignCooperative
	} else {
		return models.PublicHEI
	}
}

type ShuangGaoItem struct {
	X学校名称  string `json:"学校名称"`
	Y专业群名称 string `json:"专业群名称"`
}

// 只有首字母大写才是公有字段
type HEITagsRaw struct {
	A强基计划     []string                   `json:"强基计划"`
	B招收保送生    []string                   `json:"招收保送生"`
	C招收高水平运动队 []string                   `json:"招收高水平运动队"`
	D招收高水平艺术团 []string                   `json:"招收高水平艺术团"`
	E部委直属     map[string][]string        `json:"部委直属"`
	F双高计划     map[string][]ShuangGaoItem `json:"双高计划"`
}

type ShuangGaoHEI struct {
	X类别    string
	Y专业群名称 []string
}

type HEITagsLoaded struct {
	A强基计划     map[string]bool
	B招收保送生    map[string]bool
	C招收高水平运动队 map[string]bool
	D招收高水平艺术团 map[string]bool
	E所属部委     map[string]string // 学校名称->部委名称
	F双高计划     map[string]ShuangGaoHEI
}

func arrToBoolMap(arr []string) (retM map[string]bool) {
	retM = make(map[string]bool)
	for _, s := range arr {
		retM[s] = true
	}
	return
}

func GetTagsLoaded(heiTags *HEITagsRaw) *HEITagsLoaded {
	htl := HEITagsLoaded{}
	htl.A强基计划 = arrToBoolMap(heiTags.A强基计划)
	htl.B招收保送生 = arrToBoolMap(heiTags.B招收保送生)
	htl.C招收高水平运动队 = arrToBoolMap(heiTags.C招收高水平运动队)
	htl.D招收高水平艺术团 = arrToBoolMap(heiTags.D招收高水平艺术团)
	htl.E所属部委 = make(map[string]string)
	for 部委名称, 学校列表 := range heiTags.E部委直属 {
		for _, 学校名称 := range 学校列表 {
			htl.E所属部委[学校名称] = 部委名称
		}
	}
	htl.F双高计划 = make(map[string]ShuangGaoHEI)
	for 档次, 学校列表 := range heiTags.F双高计划 {
		for _, 学校信息 := range 学校列表 {
			htl.F双高计划[学校信息.X学校名称] = ShuangGaoHEI{X类别: 档次, Y专业群名称: strings.Split(学校信息.Y专业群名称, "、")}
		}
	}
	return &htl

}

func (htl *HEITagsLoaded) addHEIPolicies(hei *models.HEI) {
	if _, ok := htl.A强基计划[hei.Name]; ok {
		hei.Tags = append(hei.Tags, "强基计划")
	}
	if _, ok := htl.B招收保送生[hei.Name]; ok {
		hei.Tags = append(hei.Tags, "保送生")
	}
	if _, ok := htl.C招收高水平运动队[hei.Name]; ok {
		hei.Tags = append(hei.Tags, "高水平运动队")
	}
	if _, ok := htl.D招收高水平艺术团[hei.Name]; ok {
		hei.Tags = append(hei.Tags, "高水平艺术团")
	}
	if 部委名称, ok := htl.E所属部委[hei.Name]; ok {
		hei.Tags = append(hei.Tags, "部委直属")
		hei.Tags = append(hei.Tags, 部委名称+"直属")
	}
	if _, ok := htl.F双高计划[hei.Name]; ok {
		hei.Tags = append(hei.Tags, "双高计划")
	}
}

func ImportHEIs() {
	heiTags := HEITagsRaw{}
	err := utils.UnmarshallJSONFromFile("hei_tags.json", &heiTags)
	if err != nil {
		panic(err)
	}
	// fmt.Println(heiTags)
	htl := GetTagsLoaded(&heiTags)
	fmt.Println(htl)

	heis := make([]map[string]interface{}, 0)
	err = utils.UnmarshallJSONFromFile("heis.json", &heis)
	if err != nil {
		panic(err)
	}

	province := ""
	for i, majorProps := range heis {
		// majorProps
		if majorProps["code"] == nil {
			province = strings.Split(majorProps["序号"].(string), "(")[0]
			continue
		}
		hei := models.HEI{
			Code: fmt.Sprintf("%d", int(majorProps["code"].(float64))),
			Name: fmt.Sprintf("%+v", majorProps["name"]),
			Location: models.Location{
				Provincial: province,
				Municipal:  fmt.Sprintf("%+v", majorProps["region"]),
			},
			Level:    GetHEILevel(fmt.Sprintf("%+v", majorProps["level"])),
			Mode:     GetHEIMode(fmt.Sprintf("%+v", majorProps["备注"])),
			Policies: make([]string, 0),
			Tags:     make([]string, 0),
		}
		htl.addHEIPolicies(&hei)
		// fmt.Printf("%d %+v\n", i, hei)

		filter := bson.M{"name": hei.Name}
		res := db.GetCollection("heis").FindOne(context.TODO(), filter)
		var insertedDocID primitive.ObjectID
		if res.Err() != nil {
			fmt.Println("err", i, res.Err())
			res, err := db.GetCollection("heis").InsertOne(context.TODO(), hei)
			if err != nil {
				panic(err)
			}
			insertedDocID = res.InsertedID.(primitive.ObjectID)
		} else {
			panic("already existed this record!")
			// _major := models.Major{}
			// err = res.Decode(&_major)
			// insertedDocID = _major.ID
			// if err != nil {
			// 	panic(err)
			// }
		}
		fmt.Println(insertedDocID)

	}
}

func main() {
	ImportMajors()
	// ImportHEIs()
	// fmt.Println("数据库初始化成功！")
}
