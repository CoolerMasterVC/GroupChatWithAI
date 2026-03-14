import os
from dotenv import load_dotenv

load_dotenv()

TRANSPORT_URL = os.getenv("TRANSPORT_URL", "http://transport-service:8000/api/transfer")
SEGMENT_SIZE_BYTES = int(os.getenv("SEGMENT_SIZE_BYTES", 350))
OLLAMA_URL = os.getenv("OLLAMA_URL", "http://ollama:11434")
MODEL_NAME = os.getenv("MODEL_NAME", "mistral")