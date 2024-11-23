import streamlit as st
from model import analyze_skin, get_recommended_products, get_recommended_products_by_type
from PIL import Image

def main():
    image = "Logow.png"
    col1, col2, col3 = st.columns(3)
    with col2:
        st.image(image, width=400)
    st.subheader("Select Image Source")

    # Create columns for buttons
    col1, col2 = st.columns(2)

    with col1:
        upload_button = st.button("Upload Image", use_container_width=True)

    with col2:
        camera_button = st.button("Take Photo", use_container_width=True)

    # Session state to manage UI flow
    if 'upload_mode' not in st.session_state:
        st.session_state.upload_mode = None
    if upload_button:
        st.session_state.upload_mode = 'upload'
    if camera_button:
        st.session_state.upload_mode = 'camera'

    # Handle image input based on selected mode
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
                # We need to seek to the beginning of the file because we read it once for preview
                uploaded_file.seek(0)
                detected_conditions = analyze_skin(uploaded_file)

            if detected_conditions:
                st.subheader("Conditions detected:")
                for condition in detected_conditions:
                    st.write(f"- {condition}")

        # Product recommendations section (full width)
        if detected_conditions:
            # Skin type selection
            st.subheader("Select your skin type:")
            # Define skin type mapping
            skin_type_mapping = {"Oily / Normal-Oily": 1, "Dry / Dry-Normal": 2, "Normal": 3, "Combination": 4,
                                 "Sensitive": 5, "All Types": 6}
            skin_type_id = None
            cols = st.columns(3)
            for i, (skin_type_name, skin_type_value) in enumerate(skin_type_mapping.items()):
                if i < 3:  # First three items
                    with cols[i]:
                        if st.button(skin_type_name, use_container_width=True):
                            skin_type_id = skin_type_value
            cols = st.columns(3)
            for i, (skin_type_name, skin_type_value) in enumerate(skin_type_mapping.items()):
                if i >= 3:  # Next three items
                    with cols[i - 3]:
                        if st.button(skin_type_name, use_container_width=True):
                            skin_type_id = skin_type_value

            st.header("Recommended Products:")

            # Product type mapping
            product_type_mapping = {"Cleansers": 1, "Serums": 2, "Toners": 3, "Moisturisers": 4, "Sunscreens": 5}

            tab1, tab2, tab3, tab4, tab5, tab6 = st.tabs(["All Products", "Cleansers", "Serums", "Toners", "Moisturisers", "Sunscreens"])

            with tab1:
                with st.spinner('Fetching product recommendations...'):
                    recommended_products = get_recommended_products(detected_conditions, skin_type_id)
                display_products(recommended_products)
            with tab2:
                with st.spinner('Fetching product recommendations...'):
                    recommended_products = get_recommended_products_by_type(detected_conditions, skin_type_id,
                                                                            product_type_mapping["Cleansers"])
                display_products(recommended_products)
            with tab3:
                with st.spinner('Fetching product recommendations...'):
                    recommended_products = get_recommended_products_by_type(detected_conditions, skin_type_id,
                                                                            product_type_mapping["Serums"])
                display_products(recommended_products)
            with tab4:
                with st.spinner('Fetching product recommendations...'):
                    recommended_products = get_recommended_products_by_type(detected_conditions, skin_type_id,
                                                                            product_type_mapping["Toners"])
                display_products(recommended_products)
            with tab5:
                with st.spinner('Fetching product recommendations...'):
                    recommended_products = get_recommended_products_by_type(detected_conditions, skin_type_id,
                                                                            product_type_mapping["Moisturisers"])
                display_products(recommended_products)
            with tab6:
                with st.spinner('Fetching product recommendations...'):
                    recommended_products = get_recommended_products_by_type(detected_conditions, skin_type_id,
                                                                            product_type_mapping["Sunscreens"])
                display_products(recommended_products)

        else:
            st.balloons()
            st.info("No conditions detected.")


def display_products(recommended_products):
    if isinstance(recommended_products, str):  # Check if the result is a string message
        st.warning(recommended_products)
    else:
        for product in recommended_products:
            # Display product name as title and ingredients as body
            st.subheader(product["product_name"])
            col1, col2, col3 = st.columns(3)
            with col1:
                st.write("<span style=color:#ff9f0a><b>Provided by:</b></span>", product["brand"],
                         unsafe_allow_html=True)
            with col2:
                st.write("<span style=color:#ff9f0a><b>For:</b></span>", product["concern"],
                         unsafe_allow_html=True)
            with col3:
                st.write("<span style=color:#ff9f0a><b>Key Ingredient:</b></span>", product["key_ingredients"],
                         unsafe_allow_html=True)
            with st.expander("All Ingredients"):
                st.write(product["all_ingredients"])
            st.divider()


if __name__ == "__main__":
    main()