import os
import sys
import json

import grpc
import requests
import pika
from PIL import Image

from ml.proto.database_pb2 import UpdateFullPredictionRequest
from ml.proto.database_pb2_grpc import DatabaseServiceStub
from ml.proto.image_pb2 import RecognitionResult, Scope
from text_detection import detect_objects_on_image, process_ocr_result, concat_all_results

USER = 'user'
PASSWORD = 'password'
BROKER_HOSTNAME = 'rabbitmq'
PORT = 5672
QUEUE_NAME = 'predictions'

DATABASE = ''

CONNECTION_URL = f'amqp://{USER}:{PASSWORD}@{BROKER_HOSTNAME}:{PORT}/'


def handle_message(ch, method, properties, body):
    body_str = body.decode('utf-8')
    request = json.loads(body_str)

    index = request.get('index')
    image_url = request.get('image')

    image_data = Image.open(requests.get(image_url, stream=True).raw)

    contours = detect_objects_on_image(image_data)
    ocr_result = process_ocr_result()

    pb2_results = []
    for result in concat_all_results(contours, ocr_result):
        pb2_result = RecognitionResult(recognition=result[0],
                                       metric=result[1],
                                       scope=Scope(
                                           x1=result[2][0], y1=result[2][1],
                                           x2=result[2][2], y2=result[2][3]))
        pb2_results.append(pb2_result)

    with grpc.insecure_channel(DATABASE) as channel:
        stub = DatabaseServiceStub(channel)
        stub.UpdateFullPrediction(UpdateFullPredictionRequest(id=index, results=pb2_results))


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
