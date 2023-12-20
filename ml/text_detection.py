import glob
import os

import cv2
import numpy as np
import onnxruntime as ort
from PIL import Image
from paddleocr import PaddleOCR

from proto import database_pb2

detector_classes = ['water_meter_data']


def detect_objects_on_image(buf):
    input, img_width, img_height = prepare_input(buf)
    output = run_model(input)

    contours = process_output(output, img_width, img_height)
    crop_by_contours(buf, contours)

    return contours


def prepare_input(img):
    img_width, img_height = img.size

    img = img.resize((640, 640))
    img = img.convert("RGB")

    input = np.array(img)
    input = input.transpose(2, 0, 1)
    input = input.reshape(1, 3, 640, 640) / 255.0

    return input.astype(np.float32), img_width, img_height


def run_model(input):
    model = ort.InferenceSession("models/water_meters_detector.onnx", providers=['CPUExecutionProvider'])
    outputs = model.run(["output0"], {"images": input})

    return outputs[0]


def iou(box1, box2):
    return intersection(box1, box2) / union(box1, box2)


def union(box1, box2):
    box1_x1, box1_y1, box1_x2, box1_y2 = box1[:4]
    box2_x1, box2_y1, box2_x2, box2_y2 = box2[:4]
    box1_area = (box1_x2 - box1_x1) * (box1_y2 - box1_y1)
    box2_area = (box2_x2 - box2_x1) * (box2_y2 - box2_y1)

    return box1_area + box2_area - intersection(box1, box2)


def intersection(box1, box2):
    box1_x1, box1_y1, box1_x2, box1_y2 = box1[:4]
    box2_x1, box2_y1, box2_x2, box2_y2 = box2[:4]
    x1 = max(box1_x1, box2_x1)
    y1 = max(box1_y1, box2_y1)
    x2 = min(box1_x2, box2_x2)
    y2 = min(box1_y2, box2_y2)

    return (x2 - x1) * (y2 - y1)


def process_output(output, img_width, img_height):
    output = output[0].astype(float)
    output = output.transpose()

    boxes = []
    for row in output:
        prob = row[4:].max()
        if prob < 0.5: continue

        class_id = row[4:].argmax()
        label = detector_classes[class_id]

        xc, yc, w, h = row[:4]
        x1 = (xc - w / 2) / 640 * img_width
        y1 = (yc - h / 2) / 640 * img_height
        x2 = (xc + w / 2) / 640 * img_width
        y2 = (yc + h / 2) / 640 * img_height

        boxes.append([x1, y1, x2, y2, label, prob])

    boxes.sort(key=lambda x: x[5], reverse=True)

    result = []
    while len(boxes) > 0:
        result.append(boxes[0])
        boxes = [box for box in boxes if iou(box, boxes[0]) < 0.7]

    return result


def crop_by_contours(img, contours):
    img = np.array(img)

    for i, contour in enumerate(contours):
        x1, y1, x2, y2 = map(int, contour[:4])
        cv_contour = np.array([[x1, y1], [x2, y2]])

        x, y, w, h = cv2.boundingRect(cv_contour)
        cropped = img[y - 10:y + h + 10, x - 10:x + w + 10]

        Image.fromarray(cropped).save(f"./images/{i}.jpg")


def predict_text():
    os.makedirs(name="./images", exist_ok=True)

    ocr = PaddleOCR(use_angle_cls=True, lang='en', show_log=False)

    crops_list = glob.glob('./images/*.jpg')

    results = []
    for crop in crops_list:
        result = ocr.ocr(crop, cls=True)
        results.append(result)

    return results


def remove_old_images():
    folder_path = './images/'
    files_to_delete = os.listdir(folder_path)

    for image in files_to_delete:
        if image is not '.gitkeep':
            image_path = os.path.join(folder_path, image)
            os.remove(image_path)


def process_ocr_result():
    ocr_results = predict_text()
    remove_old_images()

    results = []

    for ocr_result in ocr_results:
        if len(ocr_result[0]) > 1:
            conf = ocr_result[0][0][-1][-1]
            result = ocr_result[0][0][-1][-2]
        else:
            result = ocr_result[0][-1][-1][-2]
            conf = ocr_result[0][-1][-1][-1]
        
        results.append([result, conf])

    return results


def concat_all_results(contours, ocr_result):
    processed_results = []

    for idx, result in enumerate(ocr_result):
        box = np.array(contours[idx][:4], dtype=np.int32).tolist()
        text = result[0]
        proba = result[1]

        result = database_pb2.RecognitionResult(recognition=text, metric=proba,
                                             scope=database_pb2.Scope(x1=box[0], y1=box[1], x2=box[2], y2=box[3]))

        processed_results.append(result)

    return processed_results