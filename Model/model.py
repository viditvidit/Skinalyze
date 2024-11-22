import os
import sys
import numpy as np
from PIL import Image
from ultralytics import YOLO
import requests

# Assume models are in a 'models' directory in the same repository
model_paths = {
    "pigmentation": "./models/pigmentation.pt",
    "darkspot": "./models/darkspot.pt",
    "acne": "./models/acne.pt",
}

# Load YOLO models from local paths
try:
    pigmentation_model = YOLO(model_paths["pigmentation"])
    darkspot_model = YOLO(model_paths["darkspot"])
    acne_model = YOLO(model_paths["acne"])
except Exception as e:
    print(f"Failed to load models: {e}")
    sys.exit(1)

def preprocess_image(file):
    try:
        img = Image.open(file)
        # Convert to RGB if not already
        if img.mode != 'RGB':
            img = img.convert('RGB')
        # Convert to numpy array
        img_rgb = np.array(img)
        return img, img_rgb
    except Exception as e:
        print(f"Image preprocessing error: {e}")
        raise

def detect_objects(model, image):
    results = model(image)
    boxes = results[0].boxes.data.cpu().numpy()
    return boxes

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

def get_recommended_products(detected_conditions, skin_type_id):
    API_BASE_URL = "https://clear-vision-438804-u6.el.r.appspot.com/"
    PRODUCTS_BY_CONCERN_ENDPOINT = "/products/select"
    recommended_products = []
    for condition in detected_conditions:
        # Map condition to concern ID
        concern_id = None
        if condition == "Acne":
            concern_id = 1
        elif condition == "Pigmentation" or condition == "Dark Spots":
            concern_id = 2
        else:
            continue
        # Fetch products based on both concern ID and skin type
        url = f"{API_BASE_URL}{PRODUCTS_BY_CONCERN_ENDPOINT}/{concern_id}/{skin_type_id}"
        try:
            products_by_concern_data = fetch_data(url)
            if products_by_concern_data:  # Check if data is not empty
                recommended_products.extend(products_by_concern_data)
        except requests.exceptions.RequestException as e:
            print(f"Error fetching products: {e}")
            continue
    # Return recommended products if found, otherwise indicate no products
    if recommended_products:
        return recommended_products
    else:
        return "No products for given condition and skin type."

def get_recommended_products_by_type(detected_conditions, skin_type_id, product_type_id):
    API_BASE_URL = "https://clear-vision-438804-u6.el.r.appspot.com/"
    PRODUCTS_BY_CONCERN_ENDPOINT = "/products/selectspec"
    recommended_products = []

    for condition in detected_conditions:
        # Map condition to concern ID
        concern_id = None
        if condition == "Acne":
            concern_id = 1
        elif condition == "Pigmentation" or condition == "Dark Spots":
            concern_id = 2
        else:
            continue
        # Fetch products based on concern ID, skin type, and product type
        url = f"{API_BASE_URL}{PRODUCTS_BY_CONCERN_ENDPOINT}/{concern_id}/{skin_type_id}/{product_type_id}"
        try:
            products_by_concern_data = fetch_data(url)
            if products_by_concern_data:  # Check if data is not empty
                recommended_products.extend(products_by_concern_data)
        except requests.exceptions.RequestException as e:
            print(f"Error fetching products: {e}")
            continue
    # Return recommended products if found, otherwise indicate no products
    if recommended_products:
        return recommended_products
    else:
        return "No products for given condition, skin type, and product type."

def fetch_data(url):
    response = requests.get(url)
    response.raise_for_status()  # Raise exception for bad status codes
    return response.json()