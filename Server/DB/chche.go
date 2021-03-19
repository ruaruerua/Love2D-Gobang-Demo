package DB

import (
	"log"
	"strconv"
	"time"
)

func Set(k string,v interface{},timeout time.Duration) error{
	state := redisdb.Set(k,v,timeout)
	return state.Err()
}

//func SetWithTTL(k,v interface{},ttl time.Duration){
//	redisdb.Do("set",k,v,"ex",ttl.Seconds())
//	Set()
//}

//func Get(k string)interface{}{
//	ret,_ := redisdb.Do("get",k)
//	return ret
//}

func GetState(k string)int{
	ret,_ := redisdb.Get(k).Result()
	r,err := strconv.Atoi(ret)
	if err != nil{
		log.Println(err)
	}
	return r
}

func GetAllPlayerState()[]string{
	cmd := redisdb.Do("keys", "*")
	var ret []string
	if cmd != nil{
		val := cmd.Val().([]interface{})
		for _,v := range val {
			ret = append(ret,v.(string) + "|" + strconv.Itoa(GetState(v.(string))))
		}
	}
	return ret
}

func GetPlayerState(id string)string{
	return  id + "|" + strconv.Itoa(GetState(id))
}

func Del(k string){
	redisdb.Del(k)
}

func RankList(){

}
