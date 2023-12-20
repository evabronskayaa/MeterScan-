import json
import os
import sys

import grpc
import pika
import requests
from PIL import Image

from proto import database_pb2, database_pb2_grpc
from text_detection import detect_objects_on_image, process_ocr_result, concat_all_results

USER = 'user'
PASSWORD = 'password'
BROKER_HOSTNAME = 'rabbitmq'
PORT = 5672
QUEUE_NAME = 'predictions'

DATABASE = 'database:3333'

CONNECTION_URL = f'amqp://{USER}:{PASSWORD}@{BROKER_HOSTNAME}:{PORT}/'


def handle_message(ch, method, properties, body):
    body_str = body.decode('utf-8')
    request = json.loads(body_str)

    index = request.get('index')
    image_url = request.get('image')

    image_data = Image.open(requests.get(image_url, stream=True).raw)

    contours = detect_objects_on_image(image_data)
    ocr_result = process_ocr_result()

    with grpc.insecure_channel(DATABASE) as channel:
        stub = database_pb2_grpc.DatabaseServiceStub(channel)
        stub.UpdateFullPrediction(
            database_pb2.UpdateFullPredictionRequest(id=index, results=concat_all_results(contours, ocr_result)))


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
