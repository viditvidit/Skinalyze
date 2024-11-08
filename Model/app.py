import streamlit as st
from model import analyze_skin, get_recommended_products

def main():
    # Display centered logo
    st.logo("/Users/vidit/Documents/College/Skinalyze/FrontEnd/Logow.png", size='large')

    st.header("Upload your image for analysis")
    # File uploader
    uploaded_file = st.file_uploader("", type=["jpg", "png", "jpeg"])

    if uploaded_file is not None:
        with st.spinner('Analyzing your image...'):
            detected_conditions = analyze_skin(uploaded_file)
        if detected_conditions:
            st.success("The following conditions were detected:")
            for condition in detected_conditions:
                st.write(f"- {condition}")

            # Get recommended products based on detected conditions
            with st.spinner('Fetching product recommendations...'):
                recommended_products = get_recommended_products(detected_conditions)
            st.header("Recommended Products:")
            for product in recommended_products:
                # Display product name as title and ingredients as body
                st.subheader(product["product_name"])
                st.write("**Ingredients:**")
                st.write(product["all_ingredients"])
                st.divider()  # Add a divider between products
        else:
            st.info("No conditions detected.")

if __name__ == "__main__":
    main()