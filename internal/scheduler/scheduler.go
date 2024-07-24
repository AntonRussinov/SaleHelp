package scheduler

import (
	"log"
	"salehelp/internal/client"
	"salehelp/internal/db"
	"time"
)

func StartScheduler() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		processClients()
	}
}

func processClients() {
	clients := []string{"client1", "client2"} // Получите список клиентов из вашей БД или другого источника
	for _, clientID := range clients {
		go func(id string) {
			resp, err := client.GetClientData(id)
			if err != nil {
				log.Printf("Error fetching data for client %s: %v", id, err)
				return
			}
			// Преобразование данных и сохранение в БД
			err = db.SaveClientData(resp.Body())
			if err != nil {
				log.Printf("Error saving data for client %s: %v", id, err)
				return
			}
			// Дополнительные манипуляции с данными
		}(clientID)
	}
}
