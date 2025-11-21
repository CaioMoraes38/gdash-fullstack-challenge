import requests
import pika
import json
import time
import schedule
from datetime import datetime

RABBITMQ_HOST = 'localhost' 
QUEUE_NAME = 'weather_data'
CITY_NAME = "Birigui"
LAT = "-21.2886"
LON = "-50.3411"

def get_weather_data():
    print(f"[{datetime.now()}] Buscando dados para {CITY_NAME}...")
    
    try:
        url = f"https://api.open-meteo.com/v1/forecast?latitude={LAT}&longitude={LON}&current=temperature_2m,relative_humidity_2m,is_day,precipitation,wind_speed_10m&timezone=America%2FSao_Paulo"
        
        response = requests.get(url)
        data = response.json()
        
        current = data['current']
        
        payload = {
            "city": CITY_NAME,
            "temperature": current['temperature_2m'],
            "humidity": current['relative_humidity_2m'],
            "wind_speed": current['wind_speed_10m'],
            "precipitation": current['precipitation'],
            "is_day": current['is_day'] == 1,
            "timestamp": datetime.now().isoformat()
        }
        
        send_to_rabbitmq(payload)
        
    except Exception as e:
        print(f"❌ Erro ao buscar dados: {e}")

def send_to_rabbitmq(payload):
    try:
        # Conecta ao RabbitMQ
        connection = pika.BlockingConnection(
            pika.ConnectionParameters(
                host=RABBITMQ_HOST,
                port=5672,
                credentials=pika.PlainCredentials('CaioMoraes', 'caio1234') 
            )
        )
        channel = connection.channel()
        
        # Garante que a fila existe
        channel.queue_declare(queue=QUEUE_NAME, durable=True)
        
        # Publica a mensagem
        channel.basic_publish(
            exchange='',
            routing_key=QUEUE_NAME,
            body=json.dumps(payload),
            properties=pika.BasicProperties(
                delivery_mode=2,  
            )
        )
        
        print(f"✅ Dados enviados para fila: {payload['temperature']}°C")
        connection.close()
    except Exception as e:
        print(f"❌ Erro no RabbitMQ (Verifique se o Docker está rodando): {e}")

if __name__ == "__main__":
    print("🚀 Coletor de Clima Iniciado (Ctrl+C para parar)")
    
    get_weather_data()
    
    schedule.every(1).minutes.do(get_weather_data)
    
    while True:
        schedule.run_pending()
        time.sleep(1)