import os, sys
import json

import pika
import requests


USER = 'rabbituser'
PASSWORD = 'password'
BROKER_HOSTNAME = 'rabbitmq'
QUEUE_NAME = 'predictions'

WEB_HOSTNAME = 'backend'
WEB_PORT = '80'


CONNECTION_URL = f'amqp://{USER}:{PASSWORD}@{BROKER_HOSTNAME}:5672/%2f'
print(CONNECTION_URL)


def handle_message(ch, method, properties, body):
    body_str = body.decode('utf-8')

    payload = {}
    payload_json = json.dumps(payload)

    requests.post(f'http://{WEB_HOSTNAME}:{WEB_PORT}/predictions/', data=payload_json)


def main():
    connection = pika.BlockingConnection(pika.URLParameters(CONNECTION_URL))
    channel = connection.channel()

    channel.queue_declare(queue=QUEUE_NAME)
    channel.basic_consume(queue=QUEUE_NAME, auto_ack=True, on_message_callback=handle_message)
    channel.start_consuming()


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print('Interrupted')
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)