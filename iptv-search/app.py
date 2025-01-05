from fastapi import FastAPI, Request
from sentence_transformers import SentenceTransformer
import faiss
import numpy as np
import logging
from fastapi.responses import JSONResponse

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

app = FastAPI()
model_ready = False
index_ready = False

@app.on_event("startup")
async def load_model():
    global model, index, model_ready, index_ready
    try:
        logger.info("Loading Sentence-BERT model...")
        model = SentenceTransformer('all-MiniLM-L6-v2')
        d = 384
        index = faiss.IndexFlatL2(d)
        model_ready = True
        index_ready = True
        logger.info("Model and FAISS index loaded successfully.")
    except Exception as e:
        logger.error(f"Error during startup: {e}")
        raise e

@app.get("/health")
def health_check():
    if model_ready and index_ready:
        return {"status": "ok"}
    return {"status": "loading"}, 503

metadata_store = []  # Global list to store metadata

@app.post("/add")
async def add_embeddings(request: Request):
    global metadata_store
    try:
        data = await request.json()
        sentences = data.get("metadata", [])
        if not sentences:
            return JSONResponse({"status": "error", "message": "No sentences provided"}, status_code=400)
        
        # Generate embeddings
        embeddings = model.encode(sentences)
        index.add(np.array(embeddings, dtype=np.float32))
        
        # Store metadata for later lookup
        metadata_store.extend(sentences)
        
        logger.info("Added %d embeddings to FAISS index.", len(sentences))
        return {"status": "success", "added": len(sentences)}
    except Exception as e:
        logger.error(f"Error processing request: {e}")
        return JSONResponse({"status": "error", "message": str(e)}, status_code=500)

@app.post("/query")
async def query_embedding(request: Request):
    global metadata_store
    try:
        data = await request.json()
        query = data.get("query", "")
        top_k = data.get("top_k", 5)
        if not query:
            return JSONResponse({"status": "error", "message": "Query string is empty"}, status_code=400)
        
        # Generate query vector
        query_vector = model.encode([query])[0]
        distances, indices = index.search(np.array([query_vector], dtype=np.float32), top_k)
        
        # Map indices to metadata for readable results
        results = [
            {"metadata": metadata_store[idx], "distance": float(distances[0][i])}  # Convert numpy.float32 to float
            for i, idx in enumerate(indices[0]) if idx < len(metadata_store)
        ]

        logger.info("Query successful. Returning results.")
        return {"results": results}
    except Exception as e:
        logger.error(f"Error processing query: {e}")
        return JSONResponse({"status": "error", "message": str(e)}, status_code=500)

@app.post("/save")
async def save_data():
    try:
        # Save FAISS index
        faiss.write_index(index, "/app/faiss.index")

        # Save metadata
        with open("/app/metadata.json", "w") as f:
            json.dump(metadata_store, f)

        logger.info("FAISS index and metadata saved successfully.")
        return {"status": "success", "message": "Data saved successfully"}
    except Exception as e:
        logger.error(f"Error saving data: {e}")
        return {"status": "error", "message": str(e)}

