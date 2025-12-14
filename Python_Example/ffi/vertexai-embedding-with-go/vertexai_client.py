"""
Vertex AI Multimodal Embedding Python Client

This module provides a Python interface to the Go-based Vertex AI embedding library.
"""

import ctypes
import json
import os
from typing import List, Optional, Dict, Any
import base64


class VertexAIEmbedding:
    """Client for Vertex AI Multimodal Embedding API using Go shared library."""

    def __init__(
        self,
        project_id: str,
        location: str = "us-central1",
        credentials_json: Optional[str] = None,
    ):
        """
        Initialize the Vertex AI Embedding client.

        Args:
            project_id: GCP project ID
            location: GCP region (default: us-central1)
            credentials_json: Path to credentials JSON file or JSON string.
                            If None, uses default credentials.
        """
        self.project_id = project_id
        self.location = location

        # Load credentials if provided
        self.credentials_json = ""
        if credentials_json:
            if os.path.isfile(credentials_json):
                with open(credentials_json, "r") as f:
                    self.credentials_json = f.read()
            else:
                self.credentials_json = credentials_json

        # Load the shared library
        lib_path = os.path.join(os.path.dirname(__file__), "go-ffi", "libvertexai.so")
        if not os.path.exists(lib_path):
            raise FileNotFoundError(f"Shared library not found at {lib_path}")

        self.lib = ctypes.CDLL(lib_path)

        # Define function signatures
        self.lib.GetTextEmbedding.argtypes = [
            ctypes.c_char_p,  # projectID
            ctypes.c_char_p,  # location
            ctypes.c_char_p,  # text
            ctypes.c_char_p,  # credentialsJSON
        ]
        self.lib.GetTextEmbedding.restype = ctypes.c_char_p

        self.lib.GetImageEmbedding.argtypes = [
            ctypes.c_char_p,  # projectID
            ctypes.c_char_p,  # location
            ctypes.c_char_p,  # imageBase64
            ctypes.c_char_p,  # credentialsJSON
        ]
        self.lib.GetImageEmbedding.restype = ctypes.c_char_p

    def _call_go_function(self, func, *args) -> Dict[str, Any]:
        """Call a Go function and parse the JSON result."""
        # Convert all arguments to bytes
        byte_args = [
            arg.encode("utf-8") if isinstance(arg, str) else arg for arg in args
        ]

        # Call the function
        result_ptr = func(*byte_args)

        # Convert result to Python string
        result_json = ctypes.string_at(result_ptr).decode("utf-8")

        # Parse JSON
        result = json.loads(result_json)

        return result

    def get_text_embedding(self, text: str) -> List[float]:
        """
        Get embedding vector for text.

        Args:
            text: Input text to embed

        Returns:
            List of floats representing the embedding vector

        Raises:
            RuntimeError: If the API call fails
        """
        result = self._call_go_function(
            self.lib.GetTextEmbedding,
            self.project_id,
            self.location,
            text,
            self.credentials_json,
        )

        if "error" in result and result["error"]:
            raise RuntimeError(f"Failed to get text embedding: {result['error']}")

        return result.get("text_embedding", [])

    def get_image_embedding(
        self, image_path: str = None, image_bytes: bytes = None
    ) -> List[float]:
        """
        Get embedding vector for an image.

        Args:
            image_path: Path to image file (mutually exclusive with image_bytes)
            image_bytes: Raw image bytes (mutually exclusive with image_path)

        Returns:
            List of floats representing the embedding vector

        Raises:
            ValueError: If neither or both arguments are provided
            RuntimeError: If the API call fails
        """
        if (image_path is None and image_bytes is None) or (
            image_path is not None and image_bytes is not None
        ):
            raise ValueError("Provide exactly one of image_path or image_bytes")

        # Read image if path is provided
        if image_path:
            with open(image_path, "rb") as f:
                image_bytes = f.read()

        # Encode to base64
        image_base64 = base64.b64encode(image_bytes).decode("utf-8")

        result = self._call_go_function(
            self.lib.GetImageEmbedding,
            self.project_id,
            self.location,
            image_base64,
            self.credentials_json,
        )

        if "error" in result and result["error"]:
            raise RuntimeError(f"Failed to get image embedding: {result['error']}")

        return result.get("image_embedding", [])


# Example usage
if __name__ == "__main__":
    # Initialize client
    # Option 1: Use default credentials (from GOOGLE_APPLICATION_CREDENTIALS env var)
    client = VertexAIEmbedding(project_id="your-project-id", location="us-central1")

    # Option 2: Specify credentials file
    # client = VertexAIEmbedding(
    #     project_id="your-project-id",
    #     location="us-central1",
    #     credentials_json="/path/to/credentials.json"
    # )

    # Get text embedding
    try:
        text = "Hello, this is a test sentence."
        embedding = client.get_text_embedding(text)
        print(f"Text embedding dimension: {len(embedding)}")
        print(f"First 5 values: {embedding[:5]}")
    except Exception as e:
        print(f"Error: {e}")

    # Get image embedding
    try:
        image_path = "path/to/your/image.jpg"
        embedding = client.get_image_embedding(image_path=image_path)
        print(f"Image embedding dimension: {len(embedding)}")
        print(f"First 5 values: {embedding[:5]}")
    except Exception as e:
        print(f"Error: {e}")
