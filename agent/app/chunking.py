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
    """
    Разбивает текст на сегменты, каждый из которых в UTF-8 занимает не более chunk_size байт.
    Гарантирует, что символы не разрезаются.
    """
    if not text:
        return [""]

    result = []
    current_chunk = ""
    current_chunk_bytes = 0

    for char in text:
        char_bytes = len(char.encode('utf-8'))
        # Если добавление текущего символа превысит лимит, сохраняем текущий кусок и начинаем новый
        if current_chunk_bytes + char_bytes > chunk_size:
            if current_chunk:
                result.append(current_chunk)
            # Начинаем новый кусок с текущего символа
            current_chunk = char
            current_chunk_bytes = char_bytes
        else:
            current_chunk += char
            current_chunk_bytes += char_bytes

    # Добавляем последний кусок
    if current_chunk:
        result.append(current_chunk)

    return result