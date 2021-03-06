
package main
import (
  "github.com/Syfaro/telegram-bot-api"
  "log"
  "net/http"
  "io/ioutil"
  "github.com/PuerkitoBio/goquery"
  "strings"
  "encoding/json"
  //"reflect"
  //"fmt"
  "os"
  "io"
  "bufio"
  "time"
  "sync"
  "sphere_telegram_bot/Include"
  "encoding/hex"
)

type SeenPerson struct{
	Id int64
	Name string
	Seen bool
}

type RetUrl struct{
	Url string
	Id string
}

type ConnConf struct{
	Token string `json:"Token"`
	PSWD string  `json:"PSWD"`
}


const (

	helpMsg = 	"Choose one of the following commands:\n"+
				"/all - to display all news\n" +
				"/vk - to display news from vk.com\n" +
				"/insta - to display news from instagram\n" + 
				"/fb - to display news from facebook\n"
	helloMsg = "Hi, i'm a SphereDemonis bot. I can view SphereDemonis news from all sources.\n" +
				helpMsg

	notAccessFunc = "Sorry, I can only view news from vk.com/spheredemonis so far :-)"	

	PersonPath = "person.json"	

	NewsPath = "newsid.txt"

	ConfPath = "conf.json"

	key = "4w8TLujJ72ssItNIH92kKeR3pPGVzIuu"


)

var urls = map[string]string{
	"/vk" : "https://vk.com/wall-61308004?own=1",
}

var linkClass = map[string]string{
	"/vk" : "a.wi_date",
}

var retLinks = map[string]string{
	"/vk" : "https://vk.com/spheredemonis/",
}

func SavePerson(data SeenPerson, FilePath string) {
    
	jsonData, _ := json.MarshalIndent(data, "", " ")
 
	//_ = ioutil.WriteFile("test.json", file, 0644)

	f, err := os.OpenFile(FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
	    log.Fatal(err)
	}
	
	defer f.Close()
	
	if _, err = f.Write(jsonData); err != nil {
	    log.Fatal(err)
	}
}

func SaveConf(data ConnConf, FilePath string) {
    
	jsonData, _ := json.MarshalIndent(data, "", " ")
 
	//_ = ioutil.WriteFile("test.json", file, 0644)

	f, err := os.OpenFile(FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
	    log.Fatal(err)
	}
	
	defer f.Close()
	
	if _, err = f.Write(jsonData); err != nil {
	    log.Fatal(err)
	}
}

func SaveNews(data string, FilePath string) {
    

	f, err := os.OpenFile(FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
	    log.Fatal(err)
	}
	
	defer f.Close()
	data +="\n"
	if _, err = f.WriteString(data); err != nil {
	    log.Fatal(err)
	}
}

func ReadPerson(FilePath string)(map[int64]bool, []SeenPerson) {
	//f, err := os.Open(FilePath)

	f, err := os.OpenFile(FilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	line := make([]byte, 32)

	retPerson := []SeenPerson{}

	seenId := make(map[int64]bool,1024)

	reader:=bufio.NewReader(f)

	var user SeenPerson

	for {	

		line, err = reader.ReadBytes('}')

		if err == io.EOF {
			break
		}

		err = json.Unmarshal(line, &user)

		
		if err != nil {
			log.Fatal(err)
		}

		retPerson = append(retPerson, user)
		seenId[user.Id] = user.Seen
	}

	return seenId, retPerson

}


func decodeSecret(data string) string {
	src := []byte(data)

    dst := make([]byte, hex.DecodedLen(len(src)))
    n, err := hex.Decode(dst, src)
    if err != nil {
        log.Fatal(err)
    }

	decode, err := Include.Decrypt([]byte(key), dst[:n])
	if err != nil {
		log.Printf("err dec tok")
		log.Fatal(err)
	}

	return string(decode)
}

func GetConf(FilePath string) ConnConf{

	plan, _ := ioutil.ReadFile(FilePath)
	var retConf ConnConf
	//var jsondata JsonConf
	err := json.Unmarshal(plan, &retConf)

	if err != nil {
		log.Fatal(err)
	}

	retConf.PSWD = decodeSecret(retConf.PSWD)
	retConf.Token = decodeSecret(retConf.Token)

	return retConf

}

func ReadNewsFile(FilePath string) map[string]bool {
	//f, err := os.Open(FilePath)

	f, err := os.OpenFile(FilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	line := make([]byte, 32)


	seenId := make(map[string]bool,1024)

	reader:=bufio.NewReader(f)

	//var user SeenPerson

	for {	

		line, err = reader.ReadBytes('\n')

		if err == io.EOF {
			break
		}
		

		seenId[string(line)] = true
	}

	return seenId

}

func getNews(t string) ([]RetUrl, error) {
	var retUrl []RetUrl

	var err error
	if url, ok := urls[t]; ok {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Err %v", err)
			return retUrl, err
		}
		defer resp.Body.Close()
		//log.Printf("URL %s", url)if err == nil {
		//body, _ := ioutil.ReadAll(resp.Body)
		//log.Printf("http.Get body %#v\n\n\n", string(body))
		doc, err := goquery.NewDocumentFromReader(resp.Body)

		if err != nil {
			
			log.Fatal(err)
			return retUrl, err
		}

		// Find the review items
		doc.Find(linkClass[t]).Each(func(i int, link *goquery.Selection) {
		
    		band, ok := link.Attr("href")
    	
    		if ok {
				band = strings.Replace(band, "/", "", -1)
				retUrl = append(retUrl, RetUrl{retLinks[t] + band + "?w=" + band, band})

        		//log.Printf("ahre: %v", retUrl)				
				
			}

        	//return retUrl, band
 
		})
	}

	return retUrl, err
}


//type seenID map[int64]*SeenPerson

func main() {
  	flog, err := os.OpenFile("info.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
    	log.Fatal(err)
	}
	defer flog.Close()

	infoLog := log.New(flog, "INFO\t", log.Ldate|log.Ltime)

	infoLog.Printf("start app")


	conf := GetConf(ConfPath)


  	// ???????????????????????? ?? ???????? ?? ?????????????? ????????????
  	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		infoLog.Panic(err)
		log.Panic(err)
	}

	

	bot.Debug = true
	infoLog.Printf("Authorized on account %s", bot.Self.UserName)

	seenID, person := ReadPerson(PersonPath)

	seenNews := ReadNewsFile(NewsPath)

	rss, err := getNews("/vk")
	
	if err == nil {
		for _, item := range rss {
			if seenNews[item.Id] != true {
				seenNews[item.Id] = true
				SaveNews(item.Id, NewsPath)
			}
		}
	}

	//rss, err := getNews("/vk")

	//log.Printf("rd: %v",person)
	//log.Printf("seen: %v",seenID)

	// ???????????????????????????? ??????????, ???????? ?????????? ?????????????????? ???????????????????? ???? API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	//err = bot.UpdatesChan(ucfg)
	updates, err := bot.GetUpdatesChan(ucfg)

	mu := &sync.Mutex{}
	
	go func(mu *sync.Mutex){

		c := time.Tick(5*time.Second)

		for {
			select{
			case <-c:
				mu.Lock()
				rss, err := getNews("/vk")
				mu.Unlock()
				//log.Println("_timer ")
				if err == nil {
					for _, item := range rss {
						if seenNews[item.Id] != true {
							mu.Lock()
							seenNews[item.Id] = true
							SaveNews(item.Id, NewsPath)
							mu.Unlock()
							for _, pers := range person {
 		    	  					mu.Lock()
 		    	  					msg := tgbotapi.NewMessage(pers.Id, item.Url)
 		    	      				bot.Send(msg)
 		    	      				mu.Unlock()
 		    	  				}		
						}
					}
				}
			}
		}
	}(mu)


	// ???????????? ???????????????????? ???? ????????????
	for update := range updates {

        
        if update.Message == nil {

            continue
        }

        UserTxt := update.Message.Text

        if UserTxt != "" {
        	if seenID[update.Message.Chat.ID] != true {
        		seenID[update.Message.Chat.ID] = true
        		p := SeenPerson{update.Message.Chat.ID, update.Message.Chat.FirstName, true}
        		person = append(person, p)
  				SavePerson(p, PersonPath)
    		}
        	lines := strings.Split(string(UserTxt), " ")

        	//log.Printf("lines %v", lines)

        	switch lines[0] {
        	case "/start":

                msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMsg)
                bot.Send(msg)
            case "/vk":
            	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "News from vk.com:")
                bot.Send(msg)
                rss, err := getNews(lines[0])
                //log.Printf("rss: %v, err: %v", rss, err)
                if err == nil {
                	for _, item := range rss {
                		//log.Printf("items: %s", item)
                		msg := tgbotapi.NewMessage(update.Message.Chat.ID, item.Url)
                		bot.Send(msg)
                	}
                }
            case "/fb":
            	msg := tgbotapi.NewMessage(update.Message.Chat.ID, notAccessFunc)
                bot.Send(msg)        
            case "/insta":
            	msg := tgbotapi.NewMessage(update.Message.Chat.ID, notAccessFunc)
                bot.Send(msg)
            case "/all":
            	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "All news:")
                bot.Send(msg) 
                rss, err := getNews("/vk")
                //log.Printf("rss: %v, err: %v", rss, err)
                if err == nil {
                	for _, item := range rss {
                		//log.Printf("items: %s", item)
                		msg := tgbotapi.NewMessage(update.Message.Chat.ID, item.Url)
                		bot.Send(msg)
                	}
                } 
            case "/upd":
            	if len(lines) > 2 {
            		pss := lines[1]
            		if pss == conf.PSWD {
            			upd := strings.Join(lines[2:], " ")

            			for _, item := range person {
            				msg := tgbotapi.NewMessage(item.Id, upd)
                			bot.Send(msg)
            			}

            		}
            	}    
                         
            default:
            	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMsg)
                bot.Send(msg)    
        	}
        }


        //msg := tgbotapi.NewMessage(update.Message.Chat.ID, UserTxt)
        //bot.Send(msg)
    }
}