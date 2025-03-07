import streamlit as st
import requests
import random
import google.generativeai as genai
from dotenv import load_dotenv
import os

# Load environment variables
load_dotenv()

# Config
st.set_page_config(
    page_title="Practice",
    page_icon="✍️",
    layout="centered",
)

st.title("✍️ Start Practice Writing !")
st.divider()
