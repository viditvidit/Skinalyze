import cv2
import numpy as np
from ultralytics import YOLO
import requests
import os

urls = {
    "pigmentation": "https://raw.githubusercontent.com/viditvidit/Skinalyze/master/Model/pigmentation.pt",
    "darkspot": "https://raw.githubusercontent.com/viditvidit/Skinalyze/master/Model/darkspot.pt",
    "acne": "https://raw.githubusercontent.com/viditvidit/Skinalyze/master/Model/acne.pt",
}

# Download models
for name, url in urls.items():
    local_path = f"./{name}.pt"
    if not os.path.exists(local_path):  # Download only if not already downloaded
        print(f"Downloading {name}.pt...")
        response = requests.get(url)
        with open(local_path, "wb") as f:
            f.write(response.content)
        print(f"{name}.pt downloaded!")

# Load YOLO models
pigmentation_model = YOLO('./pigmentation.pt')
darkspot_model = YOLO('./darkspot.pt')
acne_model = YOLO('./acne.pt')

def preprocess_image(file):
    img = cv2.imdecode(np.frombuffer(file.read(), np.uint8), cv2.IMREAD_COLOR)
    img_rgb = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
    return img, img_rgb

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
    API_BASE_URL = "https://star-actively-minnow.ngrok-free.app.ngrok.io/products"
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
    API_BASE_URL = "https://star-actively-minnow.ngrok-free.app.ngrok.io/products"
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