# Use Python 3.9 as the base image
FROM python:3.9

# Set environment variables for Hugging Face caching
ENV HF_HOME=/app/huggingface

# Set the working directory inside the container
WORKDIR /app

# Copy only requirements.txt first to leverage Docker's layer caching
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Preload Hugging Face model to cache it in the image
RUN python -c "from transformers import AutoModelForSequenceClassification, AutoTokenizer; \
               AutoTokenizer.from_pretrained('distilbert-base-uncased-finetuned-sst-2-english'); \
               AutoModelForSequenceClassification.from_pretrained('distilbert-base-uncased-finetuned-sst-2-english')"

# Copy the rest of the application code
COPY . .

# Expose the application port
EXPOSE 8000

# Command to run the application
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]

