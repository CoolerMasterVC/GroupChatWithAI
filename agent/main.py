from fastapi import FastAPI, BackgroundTasks
from datetime import datetime, timezone
from models import SummarizationRequest, Frame
from llm import call_llm
from transport import send_response_chunks
from config import TRANSPORT_URL, SEGMENT_SIZE_BYTES, OLLAMA_URL, MODEL_NAME
from fastapi.middleware.cors import CORSMiddleware
import httpx
import asyncio

app = FastAPI(title="Agent Service (LLM Summarization)")

# Добавляем CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

async def process_request_and_send(sender: str, timestamp: str, prompt: str):
    try:
        summary = await call_llm(prompt)
        await send_response_chunks(summary)
    except Exception as e:
        print(f"Error processing request {timestamp}: {e}")
        # Отправка кадра с ошибкой
        error_timestamp = datetime.now(timezone.utc).isoformat()
        error_frame = Frame(
            send_time=error_timestamp,
            total_segments=1,
            segment_number=0,
            payload=""
        )
        async with httpx.AsyncClient() as client:
            try:
                await client.post(TRANSPORT_URL, json=error_frame.model_dump())
            except:
                pass

@app.post("/api/summary")
async def process_request(request: SummarizationRequest, background_tasks: BackgroundTasks):
    background_tasks.add_task(process_request_and_send, request.sender, request.timestamp, request.payload)
    return {"status": "ok"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)