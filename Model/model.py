import cv2
import numpy as np
from ultralytics import YOLO
import requests

# Load the models for pigmentation, dark spot, and acne detection
pigmentation_model = YOLO('pigmentation.pt')
darkspot_model = YOLO('darkspot.pt')
acne_model = YOLO('acne.pt')

# Preprocess input image for YOLO models
def preprocess_image(file):
    img = cv2.imdecode(np.frombuffer(file.read(), np.uint8), cv2.IMREAD_COLOR)
    img_rgb = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
    return img, img_rgb

# Run object detection using YOLO
def detect_objects(model, image):
    results = model(image)
    boxes = results[0].boxes.data.cpu().numpy()
    return boxes

# Analyze the uploaded image
def analyze_skin(file):
    img, img_rgb = preprocess_image(file)
    pigmentation_boxes = detect_objects(pigmentation_model, img_rgb)
    darkspot_boxes = detect_objects(darkspot_model, img_rgb)
    acne_boxes = detect_objects(acne_model, img_rgb)

    detected_conditions = []
    if len(pigmentation_boxes) > 0:
        detected_conditions.append("Pigmentation")
    if len(darkspot_boxes) > 0:
        detected_conditions.append("Dark Spots")
    if len(acne_boxes) > 0:
        detected_conditions.append("Acne")

    return detected_conditions

# Fetch recommended products based on detected conditions
def get_recommended_products(detected_conditions):
    API_BASE_URL = "http://localhost:8080"
    PRODUCTS_BY_CONCERN_ENDPOINT = "/products/byconcern/"
    recommended_products = []
    for condition in detected_conditions:
        # Map condition to concern ID
        concern_id = None
        if condition == "Acne":
            concern_id = 1
        elif condition == "Pigmentation" or "Dark Spots":
            concern_id = 2
        else:
            continue

        # Fetch the recommended products based on the concern ID
        products_by_concern_data = fetch_data(f"{API_BASE_URL}{PRODUCTS_BY_CONCERN_ENDPOINT}{concern_id}")
        recommended_products.extend(products_by_concern_data)

    return recommended_products

# Function to fetch data from API
def fetch_data(url):
    response = requests.get(url)
    return response.json()