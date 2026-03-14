import httpx
from config import OLLAMA_URL, MODEL_NAME

async def call_llm(prompt: str) -> str:
    async with httpx.AsyncClient() as client:
        try:
            full_prompt = f"""Ты — помощник, который выделяет главное из сообщений пользователей.
Напиши краткую суммаризацию (не менее 2 предложений), отражающую основные события: чем занимался пользователь, его планы, погоду и т.д.

Пример:
Сообщение: Привет! Сегодня сходил в кино, фильм понравился. Потом поужинал в ресторане.
Суммаризация: Пользователь сходил в кино и поужинал в ресторане.

Теперь сделай аналогично:
Сообщение: {prompt}
Суммаризация:"""
            response = await client.post(
                f"{OLLAMA_URL}/api/generate",
                json={
                    "model": MODEL_NAME,
                    "prompt": full_prompt,
                    "stream": False
                },
                timeout=180.0
            )
            response.raise_for_status()
            data = response.json()
            summary = data["response"].strip()
            print(f"\n[LLM ответ] {summary}\n", flush=True)
            return summary
        except Exception as e:
            print(f"\n[Ошибка LLM] {e}\n")
            if hasattr(e, 'response') and e.response is not None:
                print(f"Статус: {e.response.status_code}, тело: {e.response.text}", flush=True)
            raise
