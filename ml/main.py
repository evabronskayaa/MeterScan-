import os
import sys
import json
from io import BytesIO
import base64

import requests
import pika

from text_detection import detect_objects_on_image, process_ocr_result, concat_all_results

USER = 'rabbituser'
PASSWORD = 'password'
BROKER_HOSTNAME = 'rabbitmq'
QUEUE_NAME = 'predictions'

WEB_HOSTNAME = 'backend'
WEB_PORT = '80'

CONNECTION_URL = f'amqp://{USER}:{PASSWORD}@{BROKER_HOSTNAME}:5672/%2f'


def handle_message(ch, method, properties, body):
    body_str = body.decode('utf-8')
    request = json.loads(body_str)
    
    index = request.get('index')
    image_data = request.get('image')
    image = BytesIO(base64.b64decode(image_data))

    contours = detect_objects_on_image(image)
    ocr_result = process_ocr_result()

    concat_results = concat_all_results(contours, ocr_result)
    
    payload = {
        'index': index,
        'results': concat_results
    }
    
    payload_json = json.dumps(payload)

    response = requests.post(f'http://{WEB_HOSTNAME}:{WEB_PORT}/predictions/', data=payload_json)


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