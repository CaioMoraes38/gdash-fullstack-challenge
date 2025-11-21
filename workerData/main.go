package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type WeatherData struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Timestamp   string  `json:"timestamp"`
}

func main() {
	// 1. Conexão com RabbitMQ
	// Usamos 'guest:guest' ou 'user:password123' dependendo de como seu docker subiu.
	// Se você resetou os volumes como sugeri antes, deve ser user:password123
	connStr := "amqp://CaioMoraes:caio1234@localhost:5672/"
	
	conn, err := amqp.Dial(connStr)
	failOnError(err, "Falha ao conectar no RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Falha ao abrir canal")
	defer ch.Close()

	// Garantimos que a fila existe (igual fizemos no Python)
	q, err := ch.QueueDeclare(
		"weather_data", // nome
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Falha ao declarar fila")

	// 2. Configura o Consumidor
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag
		false,  // auto-ack (IMPORTANTE: false, pois queremos confirmar manualmente só se der tudo certo)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Falha ao registrar consumidor")

	log.Printf(" [*] Worker Go aguardando mensagens. To stop press CTRL+C")

	// 3. Loop infinito de processamento
	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf(" [x] Recebido: %s", d.Body)

		
			var data WeatherData
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("Erro ao ler JSON: %s", err)
				d.Nack(false, false) 
				continue
			}

			// Tenta enviar para a API (que faremos no próximo passo)
			err = sendToAPI(data)
			
			if err == nil {
				// SUCESSO: Confirma pro RabbitMQ que pode apagar a mensagem
				log.Println(" [V] Sucesso! Enviado para API.")
				d.Ack(false)
			} else {
				// ERRO: A API não respondeu. 
				// Por enquanto, vamos dar Ack para não travar seu teste, 
				// mas na vida real faríamos d.Nack (rejeitar) para tentar de novo depois.
				log.Printf(" [!] Erro ao chamar API (Esperado se o NestJS não estiver rodando): %s", err)
				
				// IMPORTANTE: Estou dando Ack aqui só para limpar sua fila no teste.
				// Quando o NestJS estiver pronto, mudaremos lógica de retry.
				d.Ack(false) 
			}
		}
	}()

	<-forever
}

// Função auxiliar para chamar o Backend
func sendToAPI(data WeatherData) error {
	// URL da API que vamos criar
	apiUrl := "http://localhost:3000/api/weather" // Endereço local

	jsonData, _ := json.Marshal(data)
	
	// Simula um timeout rápido de 2 segundos
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return NewError("API retornou erro")
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Gambiarra simples para erro genérico
func NewError(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}