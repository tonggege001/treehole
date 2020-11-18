package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
)

var conn redis.Conn

var IPPORT string = "127.0.0.1:6379"

//异常恢复
func RecoverResolve(){
	if err:=recover();err!=nil{
		log.Println("A panic happened!")
		log.Println(err) // 这里的err其实就是panic传入的内容，55
	}
}

// 发送json
func JSON(data *map[string]interface{}) []byte{
	m := make(map[string]interface{})
	m["code"]= 10
	if data == nil{
		ret,_ :=json.Marshal(m)
		return ret
	}

	ret,err := json.Marshal(*data)
	if err!= nil{
		log.Printf("JSON marshall err, data=%v, err=%v",*data, err)
		return []byte{}
	}
	return ret
}

func SendJson(data * map[string]interface{}, w http.ResponseWriter){
	bytes :=JSON(data)
	_, err := w.Write(bytes)
	if err != nil{
		log.Printf("SendJson error, err=%v",err)
	}
	return
}

// Newpost
func NewPost(w http.ResponseWriter, r *http.Request){
	defer RecoverResolve()
	retMap := make(map[string]interface{})

	err := r.ParseForm()
	if err != nil{
		log.Printf("Newpost ParseForm error, err=%v", err)
		retMap["code"] = -1
		SendJson(&retMap, w)
		return
	}

	formData := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil{
		log.Printf("Newpost NewDecoder error, formData=%v\n, err=%v", r.Body, err)
		retMap["code"] = -2
		SendJson(&retMap, w)
		return
	}

	timestamp := formData["timestamp"]
	content := formData["content"]
	nickname := formData["nickname"]
	IP := r.Header.Get("X-FORWARDED-FOR")
	if IP == ""{
		IP = r.RemoteAddr
	}

	if timestamp==nil || content == nil || nickname == nil || content==""{
		log.Printf("Newpost Content error, some variable is nil. IP=%v, timestamp=%v, content=%v, nickname=%v", IP, timestamp, content, nickname)
		retMap["code"] = -3
		SendJson(&retMap, w)
		return
	}

	conn, err = redis.Dial("tcp", IPPORT)
	if err!= nil {
		log.Printf("Newpost redis Dial error, err=%v", err)
		retMap["code"] = -4
		SendJson(&retMap, w)
		return
	}

	res, err := redis.Int(conn.Do("LLEN", "treehole"))
	if err!= nil{
		log.Printf("Newpost redis get treehole LLEN error, err=%v", err)
		retMap["code"] = -5
		SendJson(&retMap, w)
		return
	}

	formData["id"] = res
	_, err = conn.Do("LPUSH", "treehole", JSON(&formData))
	if err != nil{
		log.Printf("Newpost redis treehole LPUSH error, value=%v,err=%v", JSON(&formData),err)
		retMap["code"] = -6
		SendJson(&retMap, w)
		return
	}

	retMap["code"] = 0
	SendJson(&retMap, w)
	return
}

func NewComment(w http.ResponseWriter, r *http.Request){
	defer RecoverResolve()
	retMap := make(map[string]interface{})
	err := r.ParseForm()
	if err!= nil{
		log.Printf("NewComment ParseForm error, err=%v", err)
		retMap["code"] = -1
		SendJson(&retMap, w)
		return
	}

	formData := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil{
		log.Printf("NewComment NewDecoder error, formData=%v\n, err=%v", r.Body, err)
		retMap["code"] = -2
		SendJson(&retMap, w)
		return
	}

	id := formData["id"]
	content := formData["content"]
	nickname := formData["nickname"]
	IP := r.Header.Get("X-FORWARDED-FOR")
	if IP == ""{
		IP = r.RemoteAddr
	}

	if id == nil || content == nil || content==""{
		log.Printf("NewComment id=%v, content=%v error, formData=%v\n, err=%v", id, content, r.Body, err)
		retMap["code"] = -2
		SendJson(&retMap, w)
		return
	}

	if nickname == nil{
		nickname = "匿名"
	}

	conn, err := redis.Dial("tcp", IPPORT)
	if err != nil{
		log.Printf("NewComment Redis Dial error, err=%v", err)
		retMap["code"] = -3
		SendJson(&retMap, w)
		return
	}

	dataMap := make(map[string]interface{})
	dataMap["IP"] = IP
	dataMap["content"] = content
	dataMap["nickname"] = nickname
	_, err = conn.Do("LPUSH", fmt.Sprintf("comment_%v", id),fmt.Sprintf("%v", string(JSON(&dataMap))))
	if err != nil{
		log.Printf("NewComment LPUSH error, err=%v", err)
		retMap["code"] = -4
		SendJson(&retMap, w)
		return
	}

	retMap["code"] = 0
	SendJson(&retMap, w)
	return
}

func UpGood(w http.ResponseWriter, r * http.Request){
	defer RecoverResolve()
	retMap := make(map[string]interface{})
	err := r.ParseForm()
	if err!= nil{
		log.Printf("UpGood ParseForm error, err=%v", err)
		retMap["code"] = -1
		SendJson(&retMap, w)
		return
	}

	formData := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil{
		log.Printf("UpGood NewDecoder error, formData=%v\n, err=%v", r.Body, err)
		retMap["code"] = -2
		SendJson(&retMap, w)
		return
	}

	id := formData["id"]
	IP := r.Header.Get("X-FORWARDED-FOR")
	if IP == ""{
		IP = r.RemoteAddr
	}
	if id == nil{
		log.Printf("UpGood id==nil error, formData=%v\n, err=%v", r.Body, err)
		retMap["code"] = -3
		SendJson(&retMap, w)
		return
	}

	conn, err := redis.Dial("tcp", IPPORT)
	if err != nil{
		log.Printf("UpGood redis Dial error, err=%v", err)
		retMap["code"] = -4
		SendJson(&retMap, w)
		return
	}

	_, err = conn.Do("SADD", fmt.Sprintf("up_%v", id), IP)
	if err != nil{
		log.Printf("UpGood redis SADD error, err=%v", err)
		retMap["code"] = -5
		SendJson(&retMap, w)
		return
	}

	retMap["code"] = 0
	SendJson(&retMap, w)
	return

}

func DownBad(w http.ResponseWriter, r * http.Request){
	defer RecoverResolve()
	retMap := make(map[string]interface{})
	err := r.ParseForm()
	if err!= nil{
		log.Printf("DownBad ParseForm error, err=%v", err)
		retMap["code"] = -1
		SendJson(&retMap, w)
		return
	}

	formData := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil{
		log.Printf("DownBad NewDecoder error, formData=%v\n, err=%v", r.Body, err)
		retMap["code"] = -2
		SendJson(&retMap, w)
		return
	}

	id := formData["id"]
	IP := r.Header.Get("X-FORWARDED-FOR")
	if IP == ""{
		IP = r.RemoteAddr
	}
	if id == nil{
		log.Printf("DownBad id==nil error, formData=%v\n, err=%v", r.Body, err)
		retMap["code"] = -3
		SendJson(&retMap, w)
		return
	}

	conn, err := redis.Dial("tcp", IPPORT)
	if err != nil{
		log.Printf("DownBad redis Dial error, err=%v", err)
		retMap["code"] = -4
		SendJson(&retMap, w)
		return
	}

	_, err = conn.Do("SADD", fmt.Sprintf("down_%v", id), IP)
	if err != nil{
		log.Printf("DownBad redis SADD error, err=%v", err)
		retMap["code"] = -5
		SendJson(&retMap, w)
		return
	}

	retMap["code"] = 0
	SendJson(&retMap, w)
	return

}

// 根据页数获取数据
func GetPost(w http.ResponseWriter, r *http.Request){
	defer RecoverResolve()
	retMap := make(map[string]interface{})
	err := r.ParseForm()
	if err != nil{
		log.Printf("GetPost ParseForm error, err=%v", err)
		retMap["code"] = -1
		SendJson(&retMap, w)
		return
	}

	formData := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil{
		log.Printf("GetPost NewDecoder error, formData=%v\n, err=%v", r.Body, err)
		retMap["code"] = -2
		SendJson(&retMap, w)
		return
	}

	pageCount, _ := formData["page_count"].(int)
	pageNum, _ := formData["page_num"].(int)

	conn, err = redis.Dial("tcp", IPPORT)
	if err!= nil {
		log.Printf("GetPost redis Dial error, err=%v", err)
		retMap["code"] = -4
		SendJson(&retMap, w)
		return
	}

	res, err := redis.Values(conn.Do("LRANGE", "treehole", fmt.Sprintf("%d", pageNum * pageCount),
		fmt.Sprintf("%d", (pageNum+1)* pageCount-1)))

	if err!= nil {
		log.Printf("GetPost redis LRANGE error, err=%v", err)
		retMap["code"] = -5
		SendJson(&retMap, w)
		return
	}

	var postlist = make([]string,0)
	var idlist  = make([]int, 0)
	for _, data := range res{
		postlist = append(postlist, string(data.([]byte)))
		var post  = make(map[string]interface{})
		err = json.Unmarshal(data.([]byte), &post)
		if err != nil{
			log.Printf("GetPost redis Unnarshall error, data=%v, err=%v", data, err)
			continue
		}

		id := int(post["id"].(float64))
		idlist = append(idlist, id)
	}

	// 加载评论、赞、踩
	var comment = make(map[int][]string)
	var upGood = make(map[int]interface{})
	var downBad = make(map[int]interface{})


	for _, id := range idlist{
		commentlist, err := redis.Strings(conn.Do("LRANGE", fmt.Sprintf("comment_%v", id), 0, -1))
		if err != nil{
			log.Printf("GetPost redis get comment by id error, id=%v, err=%v", id, err)
			continue
		}
		comment[id] = commentlist

		res, err := redis.Int(conn.Do("SCARD", fmt.Sprintf("up_%v", id)))
		if err != nil{
			log.Printf("GetPost redis get up by id error, id=%v, err=%v", id, err)
			continue
		}
		upGood[id] = res

		res, err = redis.Int(conn.Do("SCARD", fmt.Sprintf("down_%v", id)))
		if err != nil{
			log.Printf("GetPost redis get down by id error, id=%v, err=%v", id, err)
			continue
		}
		downBad[id] = res
	}

	retMap["postlist"] = postlist
	retMap["commentlist"] = comment
	retMap["up"] = upGood
	retMap["down"] = downBad
	retMap["code"] = 0
	SendJson(&retMap, w)
	return
}
