from concurrent import futures

import grpc

from proto import image_pb2_grpc, image_pb2
from text_detection import detect_objects_on_image, process_ocr_result, concat_all_results


class ImageProcessingService(image_pb2_grpc.ImageProcessingServiceServicer):
    def ProcessImage(self, request_iterator, context):
        for request in request_iterator:
            image = request.image

            image_path = 'images/temp'
            with open(image_path, 'wb') as f:
                f.write(image)

            contours = detect_objects_on_image(image_path)
            ocr_result = process_ocr_result()

            concat_results = concat_all_results(contours, ocr_result)

            pb2_results = []
            for result in concat_results:
                pb2_result = image_pb2.RecognitionResult(recognition=result[0],
                                                         metric=result[1],
                                                         scope=image_pb2.Scope(
                                                             x1=result[2][0], y1=result[2][1],
                                                             x2=result[2][2], y2=result[2][3]))
                pb2_results.append(pb2_result)

            yield image_pb2.ImageResponse(index=request.index, results=pb2_results)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    image_pb2_grpc.add_ImageProcessingServiceServicer_to_server(ImageProcessingService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
