from fastapi import FastAPI
from pydantic import BaseModel
import time

app = FastAPI()

class SimpleMessage(BaseModel):
    msg: str
    timestamp: str

@app.post("/message")
async def receive_message(data: SimpleMessage):
    print(f"--- New Message ---")
    print(f"Content: {data.msg}")
    print(f"Timestamp: {data.timestamp}")
    
    return {
        "status": "success",
        "message": "Message received",
        "received_at": str(time.ctime())
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=80)