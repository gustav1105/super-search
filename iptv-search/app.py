from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from sentence_transformers import SentenceTransformer
import faiss
import numpy as np
import logging
import json
from fastapi.responses import JSONResponse
import os
# Configure logging
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allows all origins; replace "*" with specific origins for better security
    allow_credentials=True,
    allow_methods=["*"],  # Allows all methods (GET, POST, etc.)
    allow_headers=["*"],  # Allows all headers
)

# Globals
model = None
index_map = {}  # Store FAISS indices by property
metadata_map = {}  # Store metadata by property

@app.on_event("startup")
async def load_model_and_initialize():
    global model, index_map, metadata_map
    try:
        logger.info("Loading Sentence-BERT model...")
        model = SentenceTransformer('all-MiniLM-L6-v2')
        
        # Define properties for embedding segregation
        properties = [
            'stream_id', 'title', 'plot', 'genre',
            'release_date', 'rating', 'director', 'cast'
        ]
        d = 384  # Embedding dimension
        
        # Check if backups exist
        metadata_path = "/app/backups/metadata.json"
        index_map = {}
        metadata_map = {}
        if os.path.exists(metadata_path):
            logger.info("Loading metadata from backup...")
            with open(metadata_path, "r") as f:
                metadata_map = json.load(f)
        else:
            metadata_map = {prop: [] for prop in properties}
        
        for prop in properties:
            index_path = f"/app/backups/faiss_{prop}.index"
            if os.path.exists(index_path):
                logger.info(f"Loading FAISS index for property '{prop}' from backup...")
                index_map[prop] = faiss.read_index(index_path)
            else:
                logger.info(f"Creating new FAISS index for property '{prop}'...")
                index_map[prop] = faiss.IndexFlatL2(d)
                if prop not in metadata_map:
                    metadata_map[prop] = []
        
        logger.info("Model, FAISS indices, and metadata initialized successfully.")
    except Exception as e:
        logger.error(f"Error during startup: {e}")
        raise e

@app.get("/health")
def health_check():
    if model:
        return {"status": "ok"}
    return {"status": "loading"}, 503

@app.post("/add")
async def add_embeddings(request: Request):
    try:
        # Log the received payload
        raw_body = await request.body()
        data = json.loads(raw_body)
        items = data.get("metadata", [])
        
        if not items:
            logger.error("No metadata provided in request.")
            return JSONResponse({"status": "error", "message": "No metadata provided"}, status_code=400)
        
        for item in items:
            metadata_entry = {}  # A dictionary to store all metadata properties
            
            for prop, value in item.items():
                metadata_entry[prop] = value  # Include all properties in metadata
                
                # Generate embeddings only for properties in index_map
                if prop in index_map and isinstance(value, str):
                    embedding = model.encode([value])[0]
                    index_map[prop].add(np.array([embedding], dtype=np.float32))
            
            # Store the entire metadata entry
            for prop in index_map.keys():  # Ensure metadata_map stores entries under existing properties
                if prop in item:
                    metadata_map[prop].append(metadata_entry)

        logger.info(f"Added {len(items)} items to FAISS indices.")
        return {"status": "success", "added": len(items)}
    except Exception as e:
        logger.error("Error in /add endpoint: %s", e, exc_info=True)
        return JSONResponse({"status": "error", "message": str(e)}, status_code=500)

@app.post("/query")
async def query_embeddings(request: Request):
    try:
        data = await request.json()
        property_name = data.get("property", "").lower()
        query_value = data.get("query", "")
        top_k = data.get("top_k", 5)
        
        if not property_name or property_name not in index_map:
            return JSONResponse({"status": "error", "message": "Invalid or missing property"}, status_code=400)
        if not query_value:
            return JSONResponse({"status": "error", "message": "Query value is empty"}, status_code=400)
        
        if property_name == 'release_date':
            # Special case for release_date: Exact or approximate match
            results = [
                {"metadata": item, "distance": 0.0}
                for item in metadata_map[property_name]
                if query_value in item.get("release_date", "")
            ]
        else:
            # Semantic search for free-text fields
            query_embedding = model.encode([query_value])[0]
            distances, indices = index_map[property_name].search(np.array([query_embedding], dtype=np.float32), top_k)
            results = [
                {"metadata": metadata_map[property_name][idx], "distance": float(distances[0][i])}
                for i, idx in enumerate(indices[0]) if idx < len(metadata_map[property_name])
            ]

        logger.info("Query successful. Returning results.")
        return {"results": results}
    except Exception as e:
        logger.error(f"Error processing query: {e}")
        return JSONResponse({"status": "error", "message": str(e)}, status_code=500)

@app.post("/save")
async def save_data():
    try:
        for prop, index in index_map.items():
            faiss.write_index(index, f"/app/faiss_{prop}.index")
        
        with open("/app/metadata.json", "w") as f:
            json.dump(metadata_map, f)

        logger.info("FAISS indices and metadata saved successfully.")
        return {"status": "success", "message": "Data saved successfully"}
    except Exception as e:
        logger.error(f"Error saving data: {e}")
        return {"status": "error", "message": str(e)}

