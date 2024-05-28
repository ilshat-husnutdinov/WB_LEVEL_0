package stansub

import (
	"log"
	"service/service/internal/config"
	libjson "service/service/internal/lib/json"
	"service/service/pkg/cache"
	"service/service/pkg/database"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
)

// printMsg выводит в логи сообщение, полученное из Nats-Streaming
func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}

// ConnectToStan инициирует соединение к Nats-Streaming
func ConnectToStan(clusterID, clientID string) (stan.Conn, error) {
	const op = "stansub.stansub.ConnectToStan"

	sc, err := stan.Connect(clusterID, clientID, stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		log.Fatalf("[%v]: connection lost, reason: [%v]", reason, op)
	}))

	return sc, err
}

// RunSubscriber запускает подписку на канал в Nats-Streaming и обрабатывает полученные данные
func RunSubscriber(sc stan.Conn, conf config.Config, db *sqlx.DB, csh *cache.Cache) {
	const op = "stansub.stansub.RunSubscriber"

	var i int
	mcb := func(msg *stan.Msg) {
		i++
		printMsg(msg, i)
		order, err := libjson.ValidateData(msg.Data)
		if err != nil {
			log.Printf("Message is not valid JSON. Can`t validate:%v", err)

		} else {
			id := order.OrderUID
			err := database.InsertOrder(db, id, order)

			if err != nil {
				log.Printf("Error when inserting data into database: %v", err)
			} else {
				csh.Set(id, order, 0)
			}
		}

	}

	_, err := sc.Subscribe(conf.STAN.Subject, mcb, stan.DurableName(conf.STAN.Durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("[%v]: Nats-Streaming server: listening subject=[%s], clientID=[%s], durable=[%s]\n", op, conf.STAN.Subject, conf.STAN.ClientID, conf.STAN.Durable)
}
