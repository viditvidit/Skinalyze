import streamlit as st
from model import analyze_skin, get_recommended_products, get_recommended_products_by_type
from PIL import Image
import os

port = int(os.environ.get("PORT", 8080))
os.system(f"streamlit config set server.port {port}")
os.system("streamlit config set server.address 0.0.0.0")


def main():
    image = "Logo.png"
    st.logo(image, size="large")
    st.subheader("Select Image Source", anchor=False)

    #Session states
    if 'upload_mode' not in st.session_state:
        st.session_state.upload_mode = None
    if 'skin_type_selected' not in st.session_state:
        st.session_state.skin_type_selected = None
    if 'detected_conditions' not in st.session_state:
        st.session_state.detected_conditions = None

    col1, col2 = st.columns(2)
    with col1:
        upload_button = st.button("Upload Image", use_container_width=True)
    with col2:
        camera_button = st.button("Take Photo", use_container_width=True)

    if upload_button:
        st.session_state.upload_mode = 'upload'
        st.session_state.skin_type_selected = None
        st.session_state.detected_conditions = None
    if camera_button:
        st.toast("Make sure to keep the camera closer to the skin for better detection")
        st.session_state.upload_mode = 'camera'
        st.session_state.skin_type_selected = None
        st.session_state.detected_conditions = None

    uploaded_file = None
    if st.session_state.upload_mode == 'upload':
        uploaded_file = st.file_uploader("Upload your image", type=["jpg", "png", "jpeg"])
    elif st.session_state.upload_mode == 'camera':
        uploaded_file = st.camera_input("Take a photo")

    if uploaded_file is not None:
        # Image analysis logic
        image = Image.open(uploaded_file)
        col1, col2 = st.columns(2)
        with col1:
            st.image(image, width=200)
        with col2:
            with st.spinner('Analyzing your image...'):
                uploaded_file.seek(0)
                st.session_state.detected_conditions = analyze_skin(uploaded_file)

            if st.session_state.detected_conditions:
                st.subheader("Conditions detected:", anchor=False)
                for condition in st.session_state.detected_conditions:
                    st.write(f"- {condition}")

        #Product recommendations section
        if st.session_state.detected_conditions:
            # Skin type selection
            st.subheader("Select your skin type:", anchor=False)
            skin_type_mapping = {
                "Oily / Normal-Oily": 1,
                "Dry / Dry-Normal": 2,
                "Normal": 3,
                "Combination": 4,
                "Sensitive": 5,
                "All Types": 6
            }

            cols = st.columns(3)
            for i, (skin_type_name, skin_type_value) in enumerate(skin_type_mapping.items()):
                if i < 3:  # First three items
                    with cols[i]:
                        if st.button(skin_type_name, use_container_width=True, key=f"skin_type_{i}"):
                            st.session_state.skin_type_selected = skin_type_value

            cols = st.columns(3)
            for i, (skin_type_name, skin_type_value) in enumerate(skin_type_mapping.items()):
                if i >= 3:  # Next three items
                    with cols[i - 3]:
                        if st.button(skin_type_name, use_container_width=True, key=f"skin_type_{i}"):
                            st.session_state.skin_type_selected = skin_type_value

            if st.session_state.skin_type_selected is not None:
                st.header("Recommended Products:", anchor=False)
                product_type_mapping = {
                    "Cleansers": 1,
                    "Serums": 2,
                    "Toners": 3,
                    "Moisturisers": 4,
                    "Sunscreens": 5
                }
                tab1, tab2, tab3, tab4, tab5, tab6 = st.tabs([
                    "All Products",
                    "Cleansers",
                    "Serums",
                    "Toners",
                    "Moisturisers",
                    "Sunscreens"
                ])
                with tab1:
                    with st.spinner('Fetching product recommendations...'):
                        recommended_products = get_recommended_products(
                            st.session_state.detected_conditions,
                            st.session_state.skin_type_selected
                        )
                    display_products(recommended_products)

                for i, (tab, (product_type, type_id)) in enumerate(zip(
                        [tab2, tab3, tab4, tab5, tab6],
                        product_type_mapping.items()
                )):
                    with tab:
                        with st.spinner('Fetching product recommendations...'):
                            recommended_products = get_recommended_products_by_type(
                                st.session_state.detected_conditions,
                                st.session_state.skin_type_selected,
                                type_id
                            )
                        display_products(recommended_products)
        else:
            st.balloons()
            st.info("No conditions detected.")

def display_products(recommended_products):
    if isinstance(recommended_products, str):
        st.warning(recommended_products)
    else:
        for idx, product in enumerate(recommended_products):
            col1, col2 = st.columns([1, 2])
            with col1:
                st.image(product["image_url"], width=160)
            with col2:
                st.write("")
                st.write(
                    f"<div style='font-size: 18px; padding-bottom:8px; color:#f57f17'><b>{product['brand']}</b></div>",
                    unsafe_allow_html=True)
                st.subheader(product["product_name"], anchor=False)
                col1, col2 = st.columns(2)
                with col1:
                    st.write("<span style=color:#f57f17><b>For:</b></span>", product["concern"], unsafe_allow_html=True)
                with col2:
                    st.write("<span style=color:#f57f17><b>Key Ingredient:</b></span>", product["key_ingredients"],
                             unsafe_allow_html=True)
            with st.expander("Product Link & Ingredients"):
                st.write(product["all_ingredients"])
                st.link_button("Visit Product Page", product["product_url"])
            st.divider()


if __name__ == "__main__":
    main()