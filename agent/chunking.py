import re
from typing import List

def split_into_sentences(text: str, max_sentences: int = 2) -> List[str]:
    sentences = re.split(r'(?<=[.!?])\s+', text.strip())
    chunks = []
    for i in range(0, len(sentences), max_sentences):
        chunk = ' '.join(sentences[i:i+max_sentences])
        chunks.append(chunk)
    return chunks

def split_into_byte_chunks(text: str, chunk_size: int) -> List[str]:
    encoded = text.encode('utf-8')
    byte_chunks = []
    for i in range(0, len(encoded), chunk_size):
        byte_chunk = encoded[i:i+chunk_size]
        byte_chunks.append(byte_chunk.decode('utf-8', errors='ignore'))
    return byte_chunks