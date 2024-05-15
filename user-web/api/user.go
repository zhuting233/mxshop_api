package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErr2Http(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误:" + e.Message(),
				})
			}
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(200, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": global.RemoveTopStruct(errs.Translate(global.Trans))})
}

func GetUserList(ctx *gin.Context) {

	userClient := global.UserSrvClient
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 [用户列表] 失败")
		HandleGrpcErr2Http(err, ctx)
		return
	}

	//claims, _ := ctx.Get("claims")
	//current := claims.(*models.CustomClaims)
	//zap.S().Infof("访问用户: %d", current.ID)
	result := make([]any, 0)
	for _, value := range rsp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			//Birthday: time.Unix(int64(value.BirthDay), 0),
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func PassWordLogin(c *gin.Context) {
	form := forms.PassWordLoginForm{}
	if err := c.ShouldBindJSON(&form); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if !store.Verify(form.CaptchaId, form.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	userClient := global.UserSrvClient

	rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: form.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
				return
			}
		}
	}

	if passRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Password:          form.PassWord,
		EncrpytedPassword: rsp.Password,
	}); err != nil || !passRsp.Success {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"password": "登录失败",
		})
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(rsp.Id),
		NickName:    rsp.NickName,
		AuthorityId: uint(rsp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "cao",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"nick_name":  rsp.NickName,
		"token":      token,
		"expired_at": claims.ExpiresAt * 1000,
	})
}

func Register(c *gin.Context) {
	form := forms.RegisterForm{}
	if err := c.ShouldBindJSON(&form); err != nil {
		HandleValidatorError(c, err)
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := rdb.Get(context.Background(), form.Mobile).Result()
	if err != nil || form.Code != val {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	}

	userClient := global.UserSrvClient
	user, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: form.Mobile,
		PassWord: form.PassWord,
		Mobile:   form.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register] 查询 [新建用户] 失败: %s", err.Error())
		HandleGrpcErr2Http(err, c)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "cao",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": claims.ExpiresAt * 1000,
	})
}
