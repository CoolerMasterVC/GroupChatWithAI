from pydantic import BaseModel, Field

class SummarizationRequest(BaseModel):
    sender: str = Field(..., description="Отправитель запроса")
    timestamp: str = Field(..., description="Время отправки запроса (идентификатор)")
    payload: str = Field(..., description="Текст для суммаризации")

class Frame(BaseModel):
    send_time: str = Field(..., description="Уникальный идентификатор сообщения")
    total_segments: int = Field(..., description="Общее количество сегментов")
    segment_number: int = Field(..., description="Номер сегмента (0-based)")
    payload: str = Field(..., description="Фрагмент данных")