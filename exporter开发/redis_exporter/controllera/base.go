package controllera

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type baseSelect struct {
	redisconn *redis.Client
}

func newbaseSelect(conredis *redis.Client) *baseSelect {
	return &baseSelect{
		redisconn: conredis,
	}
}

func (b *baseSelect) redisele(grepvalue string) (float64, error) {
	var str string

	name := b.redisconn.Info()

	err := name.Scan(&str)
	if err != nil {
		log.Print("[error]read redis info is error :", err)
	}

	regex, err := regexp.Compile("\r\n")
	if err != nil {
		log.Print("[error]read redis info regexp the is :", err)
	}

	str1 := regex.Split(str, -1)

	for _, v := range str1 {
		if strings.Contains(v, grepvalue) {
			value := strings.Split(v, ":")
			if num, err := strconv.Atoi(value[1]); err == nil {
				return float64(num), nil
			} else {
				log.Print("[error]string strconv is error : ", err)
			}
		}
	}
	return 0, fmt.Errorf("grepvalue is not found")
}
