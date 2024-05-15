package api

import (
	"context"
	"fmt"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"net/http"
	"strings"
	"time"
)

func GenerateSmsCode(width int) string {
	var sb strings.Builder
	for i := 0; i < width; i++ {
		sb.WriteRune(rune('0' + rand.Intn(10)))
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {

	form := forms.SendSmsForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	mobile := form.Mobile
	smsCode := GenerateSmsCode(6)

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	// 创建API请求并设置参数
	request := dysmsapi.CreateSendSmsRequest()

	// 该参数值为假设值，请您根据实际情况进行填写
	request.PhoneNumbers = mobile

	// 该参数值为假设值，请您根据实际情况进行填写
	request.SignName = "酷酷租车"
	request.TemplateCode = "SMS_465981476"
	request.TemplateParam = "{\"code\":" + smsCode + "}"

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 打印您需要的返回值，此处打印的是此次请求的 RequestId
	//fmt.Println(response)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//fmt.Println(global.ServerConfig.RedisInfo.Expire)
	rdb.Set(context.Background(), mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
	//val, err := rdb.Get(context.Background(), mobile).Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(val)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": response.Message,
	})
}
