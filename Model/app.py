import streamlit as st
from model import analyze_skin, get_recommended_products
from PIL import Image

def main():
    # Display centered logo
    st.logo("/Users/vidit/Documents/College/Skinalyze/FrontEnd/Logow.png", size='large')
    st.subheader("Upload your image for analysis")
    # File uploader
    uploaded_file = st.file_uploader("", type=["jpg", "png", "jpeg"])
    if uploaded_file is not None:
        # Display image preview
        image = Image.open(uploaded_file)
        # Create two columns - one for image, one for analysis
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

        # Skin type selection
        st.subheader("Select your skin type:")
        # Define skin type mapping
        skin_type_mapping = {
            "Oily / Normal-Oily": 1,
            "Dry / Dry-Normal": 2,
            "Normal": 3,
            "Combination": 4,
            "Sensitive": 5,
            "All Types": 6
        }

        # Initialize selected skin type ID (default)
        skin_type_id = None
        skin_type = None
        cols = st.columns(3)
        for i, (skin_type_name, skin_type_value) in enumerate(skin_type_mapping.items()):
            if i < 3:  # First three items
                with cols[i]:
                    if st.button(skin_type_name, use_container_width=True):
                        skin_type_id = skin_type_value
                        skin_type = skin_type_name

            # Second row of buttons (3 columns)
        cols = st.columns(3)
        for i, (skin_type_name, skin_type_value) in enumerate(skin_type_mapping.items()):
            if i >= 3:  # Next three items
                with cols[i - 3]:
                    if st.button(skin_type_name, use_container_width=True):
                        skin_type_id = skin_type_value
                        skin_type = skin_type_name

        # Product recommendations section (full width)
        if detected_conditions:
            # Get recommended products based on detected conditions and skin type
            with st.spinner('Fetching product recommendations...'):
                recommended_products = get_recommended_products(detected_conditions, skin_type_id)

            st.header("Recommended Products:")

            if isinstance(recommended_products, str):  # Check if the result is a string message
                st.warning(recommended_products)
            else:
                for product in recommended_products:
                    # Display product name as title and ingredients as body
                    st.subheader(product["product_name"])
                    col1, col2, col3 = st.columns(3)
                    with col1:
                        st.write("<span style=color:#ff9f0a><b>Provided by:</b></span>", product["brand"],unsafe_allow_html=True)
                    with col2:
                        st.write("<span style=color:#ff9f0a><b>For:</b></span>", product["concern"], unsafe_allow_html=True)
                    with col3:
                        st.write("<span style=color:#ff9f0a><b>Key Ingredient:</b></span>" , product["key_ingredients"],unsafe_allow_html=True)
                    #with st.expander("All Ingredients"):
                    st.write("<span style=color:#ff9f0a><b>All Ingredients</b></span>",unsafe_allow_html=True)
                    st.write(product["all_ingredients"])
                    st.divider()  #Add a divider between products
        else:
            st.info("No conditions detected.")

if __name__ == "__main__":
    main()
