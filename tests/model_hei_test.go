package tests

import (
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/utils"
	"testing"
)

func TestFilterHEI(t *testing.T) {
	heis, err := services.GetHEIService().FilterHEI("北京", models.VocationalHEI, models.PublicHEI, "双高计划", 0)
	fmt.Println(heis, err)
	fmt.Println(utils.ToIndentedJSON(heis))
}
