import streamlit as st
from PIL import Image
import io
from streamlit_drawable_canvas import st_canvas
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Config
st.set_page_config(
    page_title="Practice",
    page_icon="✍️",
    layout="centered",
)

st.title("✍️ Start Practice Writing!")
st.divider()

# Check if we have a sentence to practice
if "sentence" not in st.session_state or not st.session_state.sentence:
    st.warning("No sentence found. Please go back to setup and generate a sentence first.")
    if st.button("Go to Setup", key="setup_button"):
        st.switch_page("00_setup.py")
else:
    # Display the sentence to practice
    st.subheader("Write this sentence in Japanese:")
    st.write(st.session_state.sentence)

    # Provide the option to either upload an image or draw on a canvas
    option = st.radio("Choose an input method:", ("Upload Image", "Draw on Canvas"), key="input_method")

    if option == "Upload Image":
        # Upload area for handwritten practice
        st.subheader("Upload your handwritten practice:")
        uploaded_file = st.file_uploader("Choose an image file", type=["png", "jpg", "jpeg"], key="image_uploader")

        # Preview the uploaded image
        if uploaded_file is not None:
            image = Image.open(uploaded_file)
            st.image(image, caption="Your handwritten practice", use_container_width=True)

            # Save image to session state for review
            if st.button("Submit for Review", key="submit_image"):
                # Convert image to bytes for storage in session state
                img_bytes = io.BytesIO()
                image.save(img_bytes, format="PNG")
                st.session_state.practice_image = img_bytes.getvalue()

                # Add any additional info needed for review
                st.session_state.submission_time = st.session_state.get("submission_time", 0) + 1

                # Transition to review stage
                st.success("Submission successful! Redirecting to review...")
                st.session_state.stage = "review"

                # Navigate to the review page
                st.switch_page("02_review.py")

        else:
            # Show placeholder or instructions when no image is uploaded
            st.info("Please upload an image of your handwritten Japanese sentence.")

    elif option == "Draw on Canvas":
        # Canvas for drawing
        st.subheader("Draw your handwritten practice:")
        canvas_result = st_canvas(
            stroke_width=10,
            stroke_color="black",
            background_color="white",
            height=200,
            width=600,
            drawing_mode="freedraw",
            key="canvas",
        )

        # Save the drawing as an image when the user clicks submit
        if canvas_result.image_data is not None:
            image = Image.fromarray(canvas_result.image_data.astype("uint8"))

            # Display the canvas image
            st.image(image, caption="Your handwritten practice", use_container_width=True)

            # Save image to session state for review
            if st.button("Submit for Review", key="submit_canvas"):
                # Convert image to bytes for storage in session state
                img_bytes = io.BytesIO()
                image.save(img_bytes, format="PNG")
                st.session_state.practice_image = img_bytes.getvalue()

                # Add any additional info needed for review
                st.session_state.submission_time = st.session_state.get("submission_time", 0) + 1

                # Transition to review stage
                st.success("Submission successful! Redirecting to review...")
                st.session_state.stage = "review"

                # Navigate to the review page
                st.switch_page("02_review.py")


# Add a back button
if st.button("Next"):
    st.switch_page("02_review.py")