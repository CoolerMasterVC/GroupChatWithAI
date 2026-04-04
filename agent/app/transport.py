import asyncio
import httpx
from datetime import datetime, timezone
from app.models import Frame
from app.config import TRANSPORT_URL, SEGMENT_SIZE_BYTES
from app.chunking import split_into_sentences, split_into_byte_chunks

async def send_response_chunks(response_text: str, send_time: str):
    """
    Отправляет текст ответа, разбитый на логические куски (по 2-3 предложения),
    каждый кусок дополнительно разбивается на байтовые сегменты.
    Все сегменты имеют одинаковый send_time и сквозную нумерацию.
    """
    # 1. Логическое разбиение на куски по 2-3 предложения
    logical_chunks = split_into_sentences(response_text, max_sentences=2)
    if not logical_chunks:
        logical_chunks = [""]

    # 2. Для каждого логического куска получаем байтовые сегменты и собираем общий список
    all_segments = []          # список кортежей (номер_логического_куска, payload)
    for chunk_idx, chunk in enumerate(logical_chunks):
        byte_segments = split_into_byte_chunks(chunk, SEGMENT_SIZE_BYTES)
        for seg in byte_segments:
            all_segments.append((chunk_idx, seg))

    total_segments = len(all_segments)

    print(f"\n=== ОТПРАВКА ОТВЕТА (send_time={send_time}) ===", flush=True)
    print(f"Логических кусков: {len(logical_chunks)}", flush=True)
    for i, chunk in enumerate(logical_chunks):
        preview = chunk[:100] + "..." if len(chunk) > 100 else chunk
        print(f"  Логический кусок {i+1}: {preview}", flush=True)
    print(f"Всего байтовых сегментов: {total_segments}", flush=True)

    async with httpx.AsyncClient() as client:
        for seg_num, (chunk_idx, payload) in enumerate(all_segments):
            frame = Frame(
                send_time=send_time,
                total_segments=total_segments,
                segment_number=seg_num,
                payload=payload
            )

            print(f"\n--- Сегмент {seg_num+1}/{total_segments} (из логического куска {chunk_idx+1}) ---", flush=True)
            print(f"  payload: {payload}", flush=True)
            print(f"  полный frame: {frame.model_dump_json(indent=2)}", flush=True)

            # Реальная отправка
            #for attempt in range(3):
            #    try:
            #        await client.post(TRANSPORT_URL, json=frame.model_dump())
            #        break
            #    except httpx.RequestError:
            #        if attempt == 2:
            #            print(f"Failed to send segment {seg_num} for message {current_timestamp} after 3 attempts")
            #        else:
            #            await asyncio.sleep(2 ** attempt)

    print("=== ОТПРАВКА ЗАВЕРШЕНА ===\n", flush=True)