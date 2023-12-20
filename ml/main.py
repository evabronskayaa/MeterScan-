import os
import sys
import json

import requests
import pika
from PIL import Image

from text_detection import detect_objects_on_image, process_ocr_result, concat_all_results

USER = 'user'
PASSWORD = 'password'
BROKER_HOSTNAME = 'rabbitmq'
PORT = 5672
QUEUE_NAME = 'predictions'

WEB_HOSTNAME = 'web'
WEB_PORT = '8080'

CONNECTION_URL = f'amqp://{USER}:{PASSWORD}@{BROKER_HOSTNAME}:{PORT}/'


def handle_message(ch, method, properties, body):
    body_str = body.decode('utf-8')
    request = json.loads(body_str)
    
    index = request.get('index')
    image_url = request.get('image')

    image_data = Image.open(requests.get(image_url, stream = True).raw)

    contours = detect_objects_on_image(image_data)
    ocr_result = process_ocr_result()
    concat_results = concat_all_results(contours, ocr_result)
    
    payload = {
        'index': index,
        'results': concat_results
    }

    payload_json = json.dumps(payload)

    requests.post(f'http://{WEB_HOSTNAME}:{WEB_PORT}/predictions/', data=payload_json)


def main():
    connection = pika.BlockingConnection(pika.URLParameters(CONNECTION_URL))
    channel = connection.channel()

    channel.queue_declare(queue=QUEUE_NAME, durable=True)
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