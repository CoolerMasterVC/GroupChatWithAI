import asyncio
import httpx
from datetime import datetime, timezone
from models import Frame
from config import TRANSPORT_URL, SEGMENT_SIZE_BYTES
from chunking import split_into_sentences, split_into_byte_chunks

async def send_response_chunks(response_text: str):
    """
    Отправляет текст ответа в виде нескольких сообщений (по 2-3 предложения),
    каждое сообщение разбивается на сегменты размером SEGMENT_SIZE_BYTES.
    Для каждого сообщения генерируется новый timestamp.
    """
    logical_chunks = split_into_sentences(response_text, max_sentences=2)
    if not logical_chunks:
        logical_chunks = [""]

    async with httpx.AsyncClient() as client:
        for idx, logical_chunk in enumerate(logical_chunks):
            current_timestamp = datetime.now(timezone.utc).isoformat()
            byte_segments = split_into_byte_chunks(logical_chunk, SEGMENT_SIZE_BYTES)
            total_byte_segments = len(byte_segments)
            
            print(f"\n--- Сообщение с send_time: {current_timestamp} ---", flush=True)
            print(f"Логический кусок: {logical_chunk}", flush=True)   # первые 100 символов

            for seg_num, payload in enumerate(byte_segments):
                print(f"  Сегмент {seg_num+1}/{total_byte_segments}: {payload}", flush=True)   
                frame = Frame(
                    send_time=current_timestamp,
                    total_segments=total_byte_segments,
                    segment_number=seg_num,
                    payload=payload
                )
                print(frame.model_dump_json(indent=2), flush=True)
                #for attempt in range(3):
                #    try:
                #        await client.post(TRANSPORT_URL, json=frame.model_dump())
                #        break
                #    except httpx.RequestError:
                #        if attempt == 2:
                #            print(f"Failed to send segment {seg_num} for message {current_timestamp} after 3 attempts")
                #        else:
                #            await asyncio.sleep(2 ** attempt)
            print("--- Конец логической части ---\n")
        print("---Конец Сообщения ---\n")